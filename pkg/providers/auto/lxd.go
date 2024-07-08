package auto

// TODO: put behind build tag
// TODO: should be dynamic because creating the provider may make network calls
// TODO: abstract so you can provide url and credentials to this and any other provider

import (
	"fmt"
	"os"

	"github.com/ctr2cloud/ctr2cloud/pkg/providers/lxd"
)

func init() {
	provider, err := lxd.NewProvider("")
	if err != nil {
		fmt.Fprintf(os.Stderr, "lxd provider not available: %v\n", err)
		return
	}
	Providers["lxd"] = provider
}
