package compute

import "context"

type InstanceStatus struct {
	Name string
	Id   string
}

type InstanceSpec struct {
	Name  string
	Image string
}

type Address struct {
	Address            string
	Netmask            string
	Type               string
	IsPubliclyRoutable bool
}

type Provider interface {
	List() ([]InstanceStatus, error)
	Create(*InstanceSpec) error
	Delete(string) error
	GetIpAddresses(context.Context, string) ([]Address, error)
	GetCommandExecutor(string) (*CommandExecutor, error)
}
