package docker

import (
	"context"
	"fmt"

	"github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
	"github.com/ctr2cloud/ctr2cloud/pkg/provisioners/apt"
	"github.com/ctr2cloud/ctr2cloud/pkg/provisioners/systemd"
)

type Provisioner struct {
	*compute.CommandExecutor
}

func (p *Provisioner) ensureDockerSocket(ctx context.Context) (bool, error) {
	sProvisioner := systemd.Provisioner{CommandExecutor: p.CommandExecutor}

	updated, err := sProvisioner.EnsureServiceEnabledNow(ctx, "docker.socket", false)
	if err != nil {
		return updated, fmt.Errorf("ensure docker socket enabled: %w", err)
	}

	_, err = p.CommandExecutor.Exec(ctx, "docker ps")
	if err != nil {
		return updated, fmt.Errorf("unable to run docker ps: %w", err)
	}
	return updated, nil
}

func (p *Provisioner) EnsureDockerDaemon(ctx context.Context) (bool, error) {
	aProvisioner := apt.Provisioner{CommandExecutor: p.CommandExecutor}

	installed, err := aProvisioner.EnsurePackageInstalled(ctx, "docker.io")
	if err != nil {
		return false, fmt.Errorf("ensure docker.io installed: %w", err)
	}

	enabled, err := p.ensureDockerSocket(ctx)
	if err != nil {
		return installed, fmt.Errorf("ensure docker socket enabled: %w", err)
	}
	return installed || enabled, nil
}
