package compute

import (
	"bytes"
	"context"
	"strconv"

	"github.com/juju/zaputil/zapctx"
	"go.uber.org/zap"
)

type ExecStreamDataType int

const (
	_ ExecStreamDataType = iota
	ExecStreamDataTypeStdout
	ExecStreamDataTypeStderr
)

type ExecStreamResult struct {
	Data     []byte
	DataType ExecStreamDataType
	Error    error
}

type MinimalCommandExecutor interface {
	// ExecStream executes a shell command and returns the results
	ExecStream(context.Context, string) chan ExecStreamResult
	Close() error
}

type CommandExecutorError struct {
	Code int
}

// IsNotFound returns true if the error is a CommandExecutorError with code 127.
// This is the code returned by sh when a command is not found.
func (e CommandExecutorError) IsNotFound() bool {
	return e.Code == 127
}

func (e CommandExecutorError) Error() string {
	if e.IsNotFound() {
		return "command not found"
	}
	return "command failed with return code " + strconv.Itoa(e.Code)
}

func (e CommandExecutorError) Unwrap() error {
	return e
}

// CommandExecutor wraps a MinimalCommandExecutor and provides some convenient helper functions
type CommandExecutor struct {
	MinimalCommandExecutor
}

func (e *CommandExecutor) Exec(ctx context.Context, cmd string) ([]byte, error) {
	buf := new(bytes.Buffer)
	resChan := e.MinimalCommandExecutor.ExecStream(ctx, cmd)
	logger := zapctx.Logger(ctx).With(zap.String("sub", "CommandExecutor.Exec"))
	var err error
	for res := range resChan {
		// do not return the error immediately to allow resChan to close
		if res.Error != nil {
			logger.Debug("command error", zap.Error(res.Error))
			err = res.Error
		}
		if res.Data != nil {
			_, err := buf.Write(res.Data)
			if err != nil {
				return buf.Bytes(), err
			}
		}
	}
	return buf.Bytes(), err
}

func (e *CommandExecutor) ExecString(ctx context.Context, cmd string) (string, error) {
	data, err := e.Exec(ctx, cmd)
	if err != nil {
		return string(data), err
	}
	return string(data), nil
}
