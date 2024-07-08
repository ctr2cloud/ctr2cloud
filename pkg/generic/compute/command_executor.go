package compute

import (
	"bytes"
	"context"
	"strconv"
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
	for res := range resChan {
		if res.Error != nil {
			return nil, res.Error
		}
		_, err := buf.Write(res.Data)
		if err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func (e *CommandExecutor) ExecString(ctx context.Context, cmd string) (string, error) {
	data, err := e.Exec(ctx, cmd)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
