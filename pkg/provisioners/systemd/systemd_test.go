package systemd

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/ctr2cloud/ctr2cloud/internal/test"
	"github.com/ctr2cloud/ctr2cloud/pkg/provisioners/file"
)

const testEnsureServiceEnabledNowInstanceName = "test-ensure-service-enabled-now"

func TestEnsureServiceEnabledNow(t *testing.T) {
	executorFactory := test.GetLXDExecutorFactory(t, testEnsureServiceEnabledNowInstanceName)
	ctx, r := test.DefaultPreamble(t, time.Second*20)

	executor, err := executorFactory()
	r.NoError(err)

	// systemd dbus service needs some time to come up
	// TODO: move this check inside provisioner
	time.Sleep(time.Second * 1)

	sProvisioner := Provisioner{CommandExecutor: executor}
	fProvisioner := file.Provisioner{CommandExecutor: executor}

	serviceName := fmt.Sprintf("%s.service", testEnsureServiceEnabledNowInstanceName)
	servicePath := filepath.Join("/etc/systemd/system", serviceName)

	serviceContents := `[Unit]
Description=Test Service

[Service]
ExecStart=/bin/sleep 1000
Restart=always

[Install]
WantedBy=multi-user.target
	`

	updated, err := fProvisioner.EnsureFileContentsString(ctx, servicePath, serviceContents)
	r.NoError(err)
	r.True(updated)

	test.RequireIdempotence(r, func() (bool, error) {
		return sProvisioner.EnsureServiceEnabledNow(ctx, serviceName, false)
	})

	test.RequireNonIdempotence(r, func() (bool, error) {
		return sProvisioner.EnsureServiceEnabledNow(ctx, serviceName, true)
	})
}
