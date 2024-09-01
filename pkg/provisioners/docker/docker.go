package docker

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
	"github.com/ctr2cloud/ctr2cloud/pkg/provisioners/apt"
	"github.com/ctr2cloud/ctr2cloud/pkg/provisioners/systemd"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
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

type ContainerSpec struct {
	Image   string
	Name    string
	Mounts  map[string]string
	Restart bool
	Command string
}

func (s ContainerSpec) GetCommand() string {
	var cmd strings.Builder

	cmd.WriteString(fmt.Sprintf("docker run -d --name %s", s.Name))

	for hostPath, containerPath := range s.Mounts {
		cmd.WriteString(fmt.Sprintf(" -v \"%s:%s\"", hostPath, containerPath))
	}

	if s.Restart {
		cmd.WriteString(" --restart always")
	}

	cmd.WriteString(fmt.Sprintf(" %s", s.Image))

	if s.Command != "" {
		cmd.WriteString(fmt.Sprintf(" %s", s.Command))
	}

	return cmd.String()
}

func (s *ContainerSpec) matchesInspect(inspectRes dockerInspect) bool {
	if inspectRes.Config.Image != s.Image {
		return false
	}

	if inspectRes.Name != fmt.Sprintf("/%s", s.Name) {
		return false
	}

	mountSpecSlice := make([]string, 0, len(s.Mounts))
	for hostPath, containerPath := range s.Mounts {
		mountSpecSlice = append(mountSpecSlice, fmt.Sprintf("%s:%s", hostPath, containerPath))
	}
	slices.Sort(mountSpecSlice)

	mountInspectSlice := make([]string, 0, len(inspectRes.Mounts))
	for _, mount := range inspectRes.Mounts {
		mountInspectSlice = append(mountInspectSlice, fmt.Sprintf("%s:%s", mount.Source, mount.Destination))
	}
	slices.Sort(mountInspectSlice)
	if slices.Compare(mountSpecSlice, mountInspectSlice) != 0 {
		return false
	}

	if s.Restart && inspectRes.HostConfig.RestartPolicy.Name != "always" {
		return false
	}

	if s.Command != "" && s.Command != strings.Join(inspectRes.Config.Cmd, " ") {
		return false
	}

	return true
}

func (p *Provisioner) inspectContainer(ctx context.Context, name string) (dockerInspect, error) {
	inspectResBytes, err := p.CommandExecutor.Exec(ctx, fmt.Sprintf("docker inspect -f \"{{ json . }}\" %s", name))
	if err != nil {
		return dockerInspect{}, fmt.Errorf("docker inspect: %w", err)
	}

	var inspectRes dockerInspect
	err = json.Unmarshal(inspectResBytes, &inspectRes)
	if err != nil {
		return dockerInspect{}, fmt.Errorf("unmarshal inspect: %w", err)
	}

	return inspectRes, nil
}

// EnsureContainer ensures that a container is running with the desired spec.
// If the container does not match the specification, it will be deleted and recreated.
func (p *Provisioner) EnsureContainer(ctx context.Context, spec ContainerSpec) (bool, error) {
	logger := zapctx.Logger(ctx)
	inspectRes, err := p.inspectContainer(ctx, spec.Name)
	if err == nil {
		if spec.matchesInspect(inspectRes) {
			logger.Debug("container already exists with correct config", zap.String("name", spec.Name))
			return false, nil
		}
		logger.Debug("container exists but does not match spec, deleting", zap.String("name", spec.Name))
		_, err = p.CommandExecutor.Exec(ctx, fmt.Sprintf("docker rm -f %s", spec.Name))
		if err != nil {

			return false, fmt.Errorf("delete container: %w", err)
		}

	} else {
		logger.Debug("inspect container failed", zap.String("name", spec.Name), zap.Error(err))
	}

	cmd := spec.GetCommand()
	logger.Debug("creating container", zap.String("cmd", cmd))
	_, err = p.CommandExecutor.Exec(ctx, cmd)
	if err != nil {
		return true, fmt.Errorf("create container: %w", err)
	}
	return true, nil
}
