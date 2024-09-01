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

	fProvisioner := file.Provisioner{CommandExecutor: executor}
	sProvisioner := Provisioner{CommandExecutor: executor}

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

func TestOSRelease(t *testing.T) {
	executorFactory := test.GetLXDExecutorFactory(t, testEnsureServiceEnabledNowInstanceName)
	ctx, r := test.DefaultPreamble(t, time.Second*20)

	executor, err := executorFactory()
	r.NoError(err)

	sProvisioner := Provisioner{CommandExecutor: executor}

	osReleaseInfo, err := sProvisioner.GetOSRelease(ctx)
	r.NoError(err)
	r.Equal("debian", osReleaseInfo["ID"])
	r.Equal("Debian GNU/Linux", osReleaseInfo["NAME"])
}
