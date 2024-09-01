package test

import (
	"testing"

	"github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
	"github.com/ctr2cloud/ctr2cloud/pkg/providers/lxd"
	"github.com/stretchr/testify/require"
)

// GetLXDExecutorFactory returns a factory function that creates a CommandExecutor for a test LXD instance
func GetLXDExecutorFactory(t *testing.T, instanceName string) func() (*compute.CommandExecutor, error) {
	r := require.New(t)
	p, err := lxd.NewProvider("")
	r.NoError(err)

	err = p.Create(compute.InstanceSpec{
		Name:  instanceName,
		Image: "ubuntu:20.04",
	})
	r.NoError(err)

	instances, err := p.List()
	r.NoError(err)
	var instanceId string
	for _, i := range instances {
		if i.Name == instanceName {
			instanceId = i.Id
			break
		}
	}
	r.NotEmpty(instanceId)
	t.Cleanup(func() {
		err := p.Delete(instanceId)
		r.NoError(err)
	})

	return func() (*compute.CommandExecutor, error) {
		return p.GetCommandExecutor(instanceId)
	}

}
