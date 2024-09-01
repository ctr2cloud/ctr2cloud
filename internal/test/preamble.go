package test

import (
	"context"
	"os"
	"os/signal"
	"testing"
	"time"

	"github.com/juju/zaputil/zapctx"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

// DefaultPreamble sets up a test with a logger and a context with a timeout
func DefaultPreamble(t *testing.T, timeout time.Duration) (context.Context, *require.Assertions) {
	r := require.New(t)
	logger := zaptest.NewLogger(t)
	logger = logger.With(zap.String("Test", t.Name()))
	ctx := zapctx.WithLogger(context.Background(), logger)
	ctx, cancel := context.WithTimeout(ctx, timeout)
	t.Cleanup(cancel)
	t.Cleanup(func() {
		inspect := false
		_, ok := os.LookupEnv("TEST_INSPECT_ALWAYS")
		if ok {
			inspect = true
		}
		_, ok = os.LookupEnv("TEST_INSPECT")
		if t.Failed() && ok {
			inspect = true
		}
		if !inspect {
			return
		}
		t.Log("Test failed, halting for inspection. Send Ctrl+C to continue.")
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt)
		<-sigChan
		t.Log("Continuing test cleanup")
	})
	return ctx, r
}
