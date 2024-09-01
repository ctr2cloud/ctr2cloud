package file

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
	"github.com/ctr2cloud/ctr2cloud/pkg/pipeline"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

type Provisioner struct {
	*compute.CommandExecutor
}

var ErrFileNotFound = errors.New("file not found")
var ErrPermissionDenied = errors.New("permission denied")

// GetMD5Sum returns the hex encoded md5sum of a file
func (p *Provisioner) GetMD5Sum(ctx context.Context, path string) (string, error) {
	res, err := p.CommandExecutor.ExecString(ctx, "md5sum "+path)
	if err != nil {
		var cErr compute.CommandExecutorError
		if !errors.As(err, &cErr) {
			return "", fmt.Errorf("md5sum: %w", err)
		}
		if strings.Contains(res, "No such file or directory") {
			return "", ErrFileNotFound
		}
		if strings.Contains(res, "Permission denied") {
			return "", ErrPermissionDenied
		}
		return "", fmt.Errorf("md5sum: %w", err)
	}
	res = strings.Trim(res, "\n")
	lines := strings.Split(res, "\n")
	if len(lines) != 1 {
		return "", fmt.Errorf("unexpected md5sum output: %s", res)
	}
	return strings.Split(lines[0], " ")[0], nil
}

func (p *Provisioner) EnsureFileContents(ctx context.Context, path string, contents []byte) (bool, error) {
	logger := zapctx.Logger(ctx)
	targetMD5Raw := md5.Sum(contents)
	targetMD5 := hex.EncodeToString(targetMD5Raw[:])

	currentMD5, err := p.GetMD5Sum(ctx, path)
	if err == nil && currentMD5 == targetMD5 {
		logger.Debug("file already up to date", zap.String("path", path), zap.String("md5", currentMD5))
		return false, nil
	}

	b64Contents := base64.StdEncoding.EncodeToString(contents)

	writeCmd := fmt.Sprintf("echo %s | base64 -d > %s", b64Contents, path)
	logger.Debug("writing file", zap.String("path", path), zap.String("md5", targetMD5))
	_, err = p.CommandExecutor.Exec(ctx, writeCmd)
	if err != nil {
		return true, fmt.Errorf("writing file: %w", err)
	}

	finalMD5, err := p.GetMD5Sum(ctx, path)
	if err != nil {
		return true, fmt.Errorf("getting final md5sum: %w", err)
	}
	if finalMD5 != targetMD5 {
		return true, fmt.Errorf("final md5sum mismatch: expected %s, got %s", targetMD5, finalMD5)
	}
	return true, nil
}

// EnsureFileContentsP is the pipeline version of EnsureFileContents
func (p *Provisioner) EnsureFileContentsP(path string, contents []byte) pipeline.FuncT {
	return func(ctx *pipeline.Context) error {
		changed, err := p.EnsureFileContents(ctx, path, contents)
		if err != nil {
			return err
		}
		ctx.SetResult(changed)
		return nil
	}
}

func (p *Provisioner) EnsureFileContentsString(ctx context.Context, path, contents string) (bool, error) {
	return p.EnsureFileContents(ctx, path, []byte(contents))
}

// EnsureFileContentsStringP is the pipeline version of EnsureFileContentsString
func (p *Provisioner) EnsureFileContentsStringP(path, contents string) pipeline.FuncT {
	return func(ctx *pipeline.Context) error {
		changed, err := p.EnsureFileContentsString(ctx, path, contents)
		if err != nil {
			return err
		}
		ctx.SetResult(changed)
		return nil
	}
}

func (p *Provisioner) GetFileContents(ctx context.Context, path string) ([]byte, error) {
	remoteMD5Sum, err := p.GetMD5Sum(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("getting md5sum: %w", err)
	}
	contentsCmd := fmt.Sprintf("cat %s | base64 -w 0 ", path)
	encodedContents, err := p.CommandExecutor.Exec(ctx, contentsCmd)
	if err != nil {
		return nil, fmt.Errorf("cat: %w", err)
	}

	contents := make([]byte, base64.StdEncoding.DecodedLen(len(encodedContents)))
	n, err := base64.StdEncoding.Decode(contents, encodedContents)
	if err != nil {
		return nil, fmt.Errorf("decoding base64: %w", err)
	}
	contents = contents[:n]
	localMD5Bytes := md5.Sum(contents)
	localMD5 := hex.EncodeToString(localMD5Bytes[:])
	if localMD5 != remoteMD5Sum {
		return nil, fmt.Errorf("md5sum mismatch: expected %s, got %s", remoteMD5Sum, localMD5)
	}

	return contents, nil
}
