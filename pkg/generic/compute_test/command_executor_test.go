package computetest

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
	"github.com/ctr2cloud/ctr2cloud/pkg/providers/lxd"
	"github.com/juju/zaputil/zapctx"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

type testExecCase struct {
	Command              string
	ExpectedOutput       string
	ExpectedOutputPrefix string
	ExpectError          bool
	Skip                 bool
}

func (t testExecCase) requireResult(r *require.Assertions, output string, err error) {
	if t.ExpectedOutput != "" {
		r.Equal(t.ExpectedOutput, output)
	}
	if t.ExpectedOutputPrefix != "" {
		r.True(strings.HasPrefix(output, t.ExpectedOutputPrefix))
	}
	if t.ExpectError {
		r.Error(err)
		return
	}
	r.NoError(err)
}

const testExecInstanceName = "test-exec"

func TestExec(t *testing.T) {
	r := require.New(t)
	p, err := lxd.NewProvider("")
	r.NoError(err)

	err = p.Create(compute.InstanceSpec{
		Name:  testExecInstanceName,
		Image: "ubuntu:20.04",
	})
	r.NoError(err)

	instances, err := p.List()
	r.NoError(err)
	var instanceId string
	for _, i := range instances {
		if i.Name == testExecInstanceName {
			instanceId = i.Id
			break
		}
	}
	r.NotEmpty(instanceId)
	defer func() {
		err := p.Delete(instanceId)
		r.NoError(err)
	}()

	tests := []testExecCase{
		{
			Command:        "echo hello",
			ExpectedOutput: "hello\n",
			ExpectError:    false,
		},
		{
			Command:        "echo hello;false",
			ExpectedOutput: "hello\n",
			ExpectError:    true,
		},
		{
			Command:        "nonexistentcommand",
			ExpectedOutput: "",
			ExpectError:    true,
		},
		{
			Command:              "dpkg-query -W apt",
			ExpectedOutputPrefix: "apt\t",
		},
		// special characters
		{
			Command:     "dpkg-query -W docker.io",
			ExpectError: true,
		},
		// nonexsistent
		{
			Command:     "dpkg-query -W asdfasdf",
			ExpectError: true,
		},
		// no output
		{
			Command:     "false",
			ExpectError: true,
		},
		// docker.io not installable until first apt updatd
		{
			Command:     "apt install -yq docker.io",
			ExpectError: true,
		},
	}

	// run each case on its own
	for _, test := range tests {
		t.Run(test.Command, func(t *testing.T) {
			if test.Skip {
				t.SkipNow()
			}
			logger := zaptest.NewLogger(t)
			r := require.New(t)
			ctx := zapctx.WithLogger(context.Background(), logger)
			ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()
			executor, err := p.GetCommandExecutor(instanceId)
			r.NoError(err)
			output, err := executor.ExecString(ctx, test.Command)
			test.requireResult(r, output, err)
		})
	}

	// now run each case on a dirty command executor
	t.Run("exec-dirty", func(t *testing.T) {
		logger := zaptest.NewLogger(t)
		r := require.New(t)
		ctx := zapctx.WithLogger(context.Background(), logger)
		ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()
		executor, err := p.GetCommandExecutor(instanceId)
		r.NoError(err)
		for _, test := range tests {
			if test.Skip {
				continue
			}
			output, err := executor.ExecString(ctx, test.Command)
			test.requireResult(r, output, err)
		}
	})

}
