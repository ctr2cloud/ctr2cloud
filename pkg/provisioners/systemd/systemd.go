package systemd

import (
	"context"
	"fmt"
	"strings"

	"github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

type Provisioner struct {
	*compute.CommandExecutor
}

// EnsureServiceEnabledNow ensures that that a service is in the desired state
//
// forceRestart should only be used if dependencies have changed
func (p *Provisioner) EnsureServiceEnabledNow(ctx context.Context, service string, forceRestart bool) (bool, error) {
	logger := zapctx.Logger(ctx)
	checkCmd := fmt.Sprintf("systemctl is-active is-enabled %s", service)
	isEnabledOutput, err := p.CommandExecutor.Exec(ctx, checkCmd)
	isEnabledNow := err == nil
	if isEnabledNow {
		logger.Debug("service already enabled and running", zap.String("service", service))
		if !forceRestart {
			return false, nil
		}
	} else {
		logger.Debug("service not enabled or running", zap.String("service", service), zap.ByteString("output", isEnabledOutput))
	}

	enableCmdBase := "systemctl daemon-reload"
	enableCmdPart := ""
	if isEnabledNow && forceRestart {
		enableCmdPart = fmt.Sprintf("systemctl restart %s", service)
	} else {
		enableCmdPart = fmt.Sprintf("systemctl reset-failed %s; systemctl enable --now %s", service, service)
	}
	enableCmd := fmt.Sprintf("%s; %s", enableCmdBase, enableCmdPart)
	logger.Debug("enabling service", zap.String("service", service), zap.String("cmd", enableCmd))
	enableOutput, err := p.CommandExecutor.Exec(ctx, enableCmd)
	if err != nil {
		logger.Info("failed to enable service", zap.String("service", service), zap.ByteString("output", enableOutput), zap.Error(err))
		return true, fmt.Errorf("enabling service: %w", err)
	}
	return true, nil
}

func (p *Provisioner) GetOSRelease(ctx context.Context) (map[string]string, error) {
	logger := zapctx.Logger(ctx)
	cmd := "cat /etc/os-release"
	rawOSRelease, err := p.CommandExecutor.ExecString(ctx, cmd)
	if err != nil {
		return nil, fmt.Errorf("get os-release: %w", err)
	}

	rawOSRelease = strings.Trim(rawOSRelease, "\n")
	osRelease := make(map[string]string)
	for _, line := range strings.Split(rawOSRelease, "\n") {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			logger.Warn("invalid os-release line", zap.String("line", line))
			continue
		}
		osRelease[parts[0]] = strings.Trim(parts[1], "\"")
	}
	return osRelease, nil
}
