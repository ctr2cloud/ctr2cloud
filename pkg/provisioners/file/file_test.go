package file

import (
	"math/rand"
	"testing"
	"time"

	"github.com/ctr2cloud/ctr2cloud/internal/test"
	"github.com/stretchr/testify/require"
)

type testEnsureFileContentsCase struct {
	Path          string
	Contents      []byte
	ExpectError   bool
	ExpectedError error
}

const testEnsureFileContentsInstanceName = "test-ensure-file-contents"

func TestEnsureFileContents(t *testing.T) {
	executorFactory := test.GetLXDExecutorFactory(t, testEnsureFileContentsInstanceName)

	randomBytes := make([]byte, 10<<20)
	_, err := rand.New(rand.NewSource(1)).Read(randomBytes)
	require.NoError(t, err)

	tests := []testEnsureFileContentsCase{
		{
			Path:     "/tmp/hello-test",
			Contents: []byte("hello"),
		},
		{
			Path:          "/sys/nonexistent",
			Contents:      []byte("hello"),
			ExpectError:   true,
			ExpectedError: ErrPermissionDenied,
		},
		{
			Path:     "/tmp/test2",
			Contents: randomBytes,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Path, func(t *testing.T) {
			ctx, r := test.DefaultPreamble(t, time.Second*10)
			executor, err := executorFactory()
			r.NoError(err)

			provisioner := Provisioner{executor}

			updated, err := provisioner.EnsureFileContents(ctx, tc.Path, tc.Contents)
			if tc.ExpectError {
				r.Error(err)
				if tc.ExpectedError != nil {
					r.ErrorAs(err, &tc.ExpectedError)
				}
				return
			}
			r.NoError(err)
			r.True(updated)

			// now ensure that second call does nothing
			updated, err = provisioner.EnsureFileContents(ctx, tc.Path, tc.Contents)
			r.NoError(err)
			r.False(updated)
		})
	}

}

const testGetFileContentsInstanceName = "test-get-file-contents"

func TestGetFileContents(t *testing.T) {
	executor, err := test.GetLXDExecutorFactory(t, testGetFileContentsInstanceName)()
	require.NoError(t, err)
	ctx, _ := test.DefaultPreamble(t, time.Second*10)

	fProvisioner := Provisioner{executor}

	paths := []string{"/etc/passwd", "/etc/group", "/etc/hosts", "/usr/bin/zcmp", "/usr/bin/zgrep"}

	for _, path := range paths {
		t.Run(path, func(t *testing.T) {
			_, err = fProvisioner.GetFileContents(ctx, path)
			require.NoError(t, err)
		})
	}

}
