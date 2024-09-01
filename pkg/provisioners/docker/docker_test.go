package docker

import (
	"testing"
	"time"

	"github.com/ctr2cloud/ctr2cloud/internal/test"
)

const testEnsureDockerDaemon = "test-ensure-docker-daemon"

func TestDocker(t *testing.T) {
	executorFactory := test.GetLXDExecutorFactory(t, testEnsureDockerDaemon)
	ctx, r := test.DefaultPreamble(t, time.Minute*5)

	executor, err := executorFactory()
	r.NoError(err)

	provisioner := Provisioner{CommandExecutor: executor}

	test.RequireIdempotence(r, func() (bool, error) {
		return provisioner.EnsureDockerDaemon(ctx)
	})

	spec := ContainerSpec{
		Image: "nginx",
		Name:  "nginx",
		Mounts: map[string]string{
			"/var/www": "/var/www",
		},
		Restart: true,
	}
	test.RequireIdempotence(r, func() (bool, error) {
		return provisioner.EnsureContainer(ctx, spec)
	})
}
