package apt

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
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

func (p *Provisioner) EnsurePackageInstalled(ctx context.Context, packageName string) error {
	logger := zapctx.Logger(ctx)
	_, err := p.GetPackageVersion(ctx, packageName)
	if err == nil {
		logger.Debug("package already installed", zap.String("package", packageName))
		return nil
	}

	aptUpdateRes, err := p.CommandExecutor.Exec(ctx, "apt update")
	logger.Debug("apt update", zap.Error(err), zap.ByteString("output", aptUpdateRes))
	if err != nil {
		return fmt.Errorf("apt update: %w", err)
	}

	aptInstallRes, err := p.CommandExecutor.Exec(ctx, "apt install -qy "+packageName)
	logger.Debug("apt install", zap.Error(err), zap.ByteString("output", aptInstallRes))
	if err != nil {
		return fmt.Errorf("apt install: %w", err)
	}
	return nil
}
