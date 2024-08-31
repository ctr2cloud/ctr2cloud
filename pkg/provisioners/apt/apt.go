package apt

import (
	"context"
	"errors"
	"fmt"
	"path"
	"strings"

	"github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
	"github.com/ctr2cloud/ctr2cloud/pkg/provisioners/file"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

type Provisioner struct {
	*compute.CommandExecutor
}

func (p *Provisioner) Update(ctx context.Context) error {
	_, err := p.CommandExecutor.Exec(ctx, "apt update")
	return err
}

var ErrNotFound = errors.New("package not found")

func (p *Provisioner) GetPackageVersion(ctx context.Context, packageName string) (string, error) {
	logger := zapctx.Logger(ctx)
	currentState, err := p.CommandExecutor.ExecString(ctx, "dpkg-query -W "+packageName)
	logger.Debug("dpkg-query", zap.Error(err), zap.String("output", currentState))
	if err != nil {
		return "", fmt.Errorf("dpkg query: %w", err)
	}
	stateLines := strings.Split(currentState, "\n")
	if len(stateLines) != 1 {
		return "", fmt.Errorf("unexpected dpkg-query output: %s", currentState)
	}

	// example line if match:
	// openssh-server	1:8.2p1-4ubuntu0.11
	//
	// example line if no match:
	// ssh
	stateLine := strings.Split(stateLines[0], "\t")
	if len(stateLine) != 2 {
		return "", ErrNotFound
	}
	return stateLine[1], nil
}

func (p *Provisioner) EnsurePackageInstalled(ctx context.Context, packageName string) (bool, error) {
	logger := zapctx.Logger(ctx)
	_, err := p.GetPackageVersion(ctx, packageName)
	if err == nil {
		logger.Debug("package already installed", zap.String("package", packageName))
		return false, nil
	}

	aptUpdateRes, err := p.CommandExecutor.Exec(ctx, "apt update")
	logger.Debug("apt update", zap.Error(err), zap.ByteString("output", aptUpdateRes))
	if err != nil {
		return false, fmt.Errorf("apt update: %w", err)
	}

	aptInstallRes, err := p.CommandExecutor.Exec(ctx, "apt install -qy "+packageName)
	logger.Debug("apt install", zap.Error(err), zap.ByteString("output", aptInstallRes))
	if err != nil {
		return false, fmt.Errorf("apt install: %w", err)
	}
	return true, nil
}

// ensureAptKey ensures that the given ASCII armorred key is installed in the apt keyring
func (p *Provisioner) ensureAptKey(ctx context.Context, keyName, key string) (bool, error) {
	fProvisioner := file.Provisioner{CommandExecutor: p.CommandExecutor}
	keyPath := path.Join("/etc/apt/trusted.gpg.d", keyName+".asc")
	return fProvisioner.EnsureFileContentsString(ctx, keyPath, key)
}

// EnsureRepository ensures that the given repository is added to the apt sources
func (p *Provisioner) ensureRepository(ctx context.Context, name, specification string) (bool, error) {
	fProvisioner := file.Provisioner{CommandExecutor: p.CommandExecutor}
	repoPath := path.Join("/etc/apt/sources.list.d", name+".list")
	return fProvisioner.EnsureFileContentsString(ctx, repoPath, specification)
}

type EnsureRepositoryArgs struct {
	Name string
	// ASCII armored key
	Key string
	// sources.list line(s) for the repository
	Specification string
	// whether to update the apt cache after adding the repository
	Update bool
}

// EnsureRepository ensures that the given repository is added to the apt sources
// TODO: detect transport https
// TODO: support arch/codename interpolation
func (p *Provisioner) EnsureRepository(ctx context.Context, args EnsureRepositoryArgs) (bool, error) {
	logger := zapctx.Logger(ctx)
	keyUpdated, err := p.ensureAptKey(ctx, args.Name, args.Key)
	if err != nil {
		return false, fmt.Errorf("ensuring apt key: %w", err)
	}
	repositoryUpdated, err := p.ensureRepository(ctx, args.Name, args.Specification)
	if err != nil {
		return false, fmt.Errorf("ensuring repository: %w", err)
	}
	if args.Update && (keyUpdated || repositoryUpdated) {
		aptUpdateRes, err := p.CommandExecutor.Exec(ctx, "apt update")
		logger.Debug("apt update", zap.Error(err), zap.ByteString("output", aptUpdateRes))
		if err != nil {
			return false, fmt.Errorf("apt update after repo add: %w", err)
		}
	}
	return keyUpdated || repositoryUpdated, nil
}
