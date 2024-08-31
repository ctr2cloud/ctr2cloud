package test

import (
	"context"
	"testing"
	"time"

	"github.com/juju/zaputil/zapctx"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

// DefaultPreamble sets up a test with a logger and a context with a timeout
func DefaultPreamble(t *testing.T, timeout time.Duration) (context.Context, *require.Assertions) {
	r := require.New(t)
	logger := zaptest.NewLogger(t)
	ctx := zapctx.WithLogger(context.Background(), logger)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	t.Cleanup(cancel)
	return ctx, r
}
