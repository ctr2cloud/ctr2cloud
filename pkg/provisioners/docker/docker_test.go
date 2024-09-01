package docker

import (
	"testing"
	"time"

	"github.com/ctr2cloud/ctr2cloud/internal/test"
)

const testEnsureDockerDaemon = "test-ensure-docker-daemon"

func TestEnsureDockerDaemon(t *testing.T) {
	executorFactory := test.GetLXDExecutorFactory(t, testEnsureDockerDaemon)
	ctx, r := test.DefaultPreamble(t, time.Second*60)

	executor, err := executorFactory()
	r.NoError(err)

	provisioner := Provisioner{CommandExecutor: executor}

	test.RequireIdempotence(r, func() (bool, error) {
		return provisioner.EnsureDockerDaemon(ctx)
	})
}
