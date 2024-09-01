package computetest

import (
	"strings"
	"testing"
	"time"

	"github.com/ctr2cloud/ctr2cloud/internal/test"
	"github.com/stretchr/testify/require"
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
	executorFactory := test.GetLXDExecutorFactory(t, testExecInstanceName)

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
		// docker.io not installable until first apt updated
		{
			Command:     "apt install -yq docker.io",
			ExpectError: false,
		},
	}

	// run each case on its own
	for _, tc := range tests {
		t.Run(tc.Command, func(t *testing.T) {
			if tc.Skip {
				t.SkipNow()
			}
			ctx, r := test.DefaultPreamble(t, 10*time.Second)
			executor, err := executorFactory()
			r.NoError(err)
			output, err := executor.ExecString(ctx, tc.Command)
			tc.requireResult(r, output, err)
		})
	}

	// now run each case on a dirty command executor
	t.Run("exec-dirty", func(t *testing.T) {
		ctx, r := test.DefaultPreamble(t, 10*time.Second)
		executor, err := executorFactory()
		r.NoError(err)
		for _, tc := range tests {
			if tc.Skip {
				continue
			}
			output, err := executor.ExecString(ctx, tc.Command)
			tc.requireResult(r, output, err)
		}
	})

}
