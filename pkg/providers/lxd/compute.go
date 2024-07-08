package lxd

import (
	"context"
	"fmt"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
	"github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
	"github.com/samber/lo"

	compute_internal "github.com/ctr2cloud/ctr2cloud/internal/generic/compute"
)

const ctr2cloudKey = "user.ctr2cloud"
const ctr2cloudNameKey = "user.ctr2cloud-name"

var _ compute.Provider = &Provider{}

type Provider struct {
	client lxd.InstanceServer
}

func NewProvider(url string) (*Provider, error) {
	c := &Provider{}
	var err error
	c.client, err = lxd.ConnectLXDUnix("", &lxd.ConnectionArgs{})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (p *Provider) List() ([]compute.InstanceStatus, error) {
	instances, err := p.client.GetInstancesWithFilter(api.InstanceTypeContainer, []string{fmt.Sprintf("config.%s=true", ctr2cloudKey)})
	if err != nil {
		return []compute.InstanceStatus{}, fmt.Errorf("getting instances: %w", err)
	}
	return lo.Map(instances, func(i api.Instance, _ int) compute.InstanceStatus {
		return compute.InstanceStatus{
			Name: i.Config[ctr2cloudNameKey],
			Id:   i.Name,
		}
	}), nil
}

func (p *Provider) Create(spec compute.InstanceSpec) error {
	id := fmt.Sprintf("ctr2cloud-%s-%s", spec.Name, lo.RandomString(5, lo.LettersCharset))
	createOp, err := p.client.CreateInstance(api.InstancesPost{
		Name: id,
		Source: api.InstanceSource{
			Type:     "image",
			Protocol: "simplestreams",
			Server:   "https://cloud-images.ubuntu.com/releases",
			Alias:    "22.04",
		},
		InstancePut: api.InstancePut{
			Config: map[string]string{
				"user.ctr2cloud":      "true",
				"user.ctr2cloud-name": spec.Name,
				"security.nesting":    "true",
			},
			Profiles: []string{"default"},
		},
		Type: api.InstanceTypeContainer,
	})
	if err != nil {
		return fmt.Errorf("creating container: %w", err)
	}
	err = createOp.Wait()
	if err != nil {
		return fmt.Errorf("waiting for container creation: %w", err)
	}
	startOp, err := p.client.UpdateInstanceState(id, api.InstanceStatePut{Action: "start"}, "")
	if err != nil {
		return fmt.Errorf("starting container: %w", err)
	}
	err = startOp.Wait()
	if err != nil {
		return fmt.Errorf("waiting for container start: %w", err)
	}
	return nil
}

func (p *Provider) Delete(id string) error {
	stopOp, err := p.client.UpdateInstanceState(id, api.InstanceStatePut{Action: "stop"}, "")
	if err != nil {
		return fmt.Errorf("stopping container: %w", err)
	}
	err = stopOp.Wait()
	if err != nil {
		return fmt.Errorf("waiting for container stop: %w", err)
	}
	deleteOp, err := p.client.DeleteContainer(id)
	if err != nil {
		return fmt.Errorf("deleting container: %w", err)
	}
	err = deleteOp.Wait()
	if err != nil {
		return fmt.Errorf("waiting for container deletion: %w", err)
	}
	return nil
}

func (p *Provider) GetCommandExecutor(id string) (*compute.CommandExecutor, error) {
	executor := compute_internal.NewPrimitiveCommandExecutor()
	stdin, stdout, stderr := executor.GetShellIO()
	op, err := p.client.ExecContainer(id, api.ContainerExecPost{
		Command:   compute_internal.PrimitiveCommandExecutorShell,
		WaitForWS: true,
	}, &lxd.ContainerExecArgs{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		return nil, fmt.Errorf("getting command exeuctor: %w", err)
	}
	go func() {
		err := op.Wait()
		if err != nil {
			fmt.Printf("command error: %v\n", err)
		}
	}()
	return &compute.CommandExecutor{MinimalCommandExecutor: executor}, nil
}

func (p *Provider) GetIpAddresses(ctx context.Context, id string) ([]compute.Address, error) {
	executor, err := p.GetCommandExecutor(id)
	if err != nil {
		return []compute.Address{}, fmt.Errorf("getting command executor: %w", err)
	}
	res, err := executor.ExecString(context.Background(), "ip addr")
	if err != nil {
		return []compute.Address{}, fmt.Errorf("executing ip a: %w", err)
	}

	return compute_internal.ParseIPAddrOutput(res), nil

}
