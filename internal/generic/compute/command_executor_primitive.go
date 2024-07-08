package compute

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"time"

	public "github.com/ctr2cloud/ctr2cloud/pkg/generic/compute"
)

var PrimitiveCommandExecutorShell = []string{"env", "PS1=___${?}___> ", "/bin/sh", "-i"}

var PrimitiveCommandExecutorShellRegex = regexp.MustCompile(`___([\d]+)___> `)

// PrimitiveCommandExecutor allows running multiple commands in a single shell session.
// This is useful for providers which do not use SSH but still allow some kind of shell access.
type PrimitiveCommandExecutor struct {
	stdinReader io.ReadCloser
	stdinWriter io.WriteCloser

	stdout chan []byte
	stderr chan []byte

	firstCommand bool
}

func (e *PrimitiveCommandExecutor) GetShellIO() (io.ReadCloser, io.WriteCloser, io.WriteCloser) {
	return e.stdinReader, &chanByteWriterCloser{e.stdout}, &chanByteWriterCloser{e.stderr}
}

func NewPrimitiveCommandExecutor() *PrimitiveCommandExecutor {
	inRead, outRead := io.Pipe()
	stdout := make(chan []byte)
	stderr := make(chan []byte)
	return &PrimitiveCommandExecutor{
		stdinReader:  inRead,
		stdinWriter:  outRead,
		stdout:       stdout,
		stderr:       stderr,
		firstCommand: true,
	}
}

func (e *PrimitiveCommandExecutor) ExecStream(ctx context.Context, cmd string) chan public.ExecStreamResult {
	resChan := make(chan public.ExecStreamResult)

	// if this is the first command, we need to wait for a clean shell
	if e.firstCommand {
		firstCtx, cancel := context.WithTimeout(ctx, time.Second*2)
		defer cancel()
		loop := true
		for loop {
			select {
			case <-firstCtx.Done():
				resChan <- public.ExecStreamResult{Error: ctx.Err()}
				close(resChan)
				return resChan
			case data, ok := <-e.stderr:
				if !ok {
					resChan <- public.ExecStreamResult{Error: io.EOF}
					close(resChan)
					return resChan
				}
				shellMatch := PrimitiveCommandExecutorShellRegex.FindAllIndex(data, 10)
				if shellMatch == nil {
					continue
				}
				// if the shell prompt is at the end of the data, we are ready for a command
				if shellMatch[len(shellMatch)-1][1] == len(data) {
					loop = false
				}
			}
		}
		e.firstCommand = false
	}

	_, err := e.stdinWriter.Write([]byte(cmd + "\n"))
	if err != nil {
		resChan <- public.ExecStreamResult{Error: err}
		close(resChan)
		return resChan
	}

	go func() {
		defer func() {
			close(resChan)
		}()
		for {
			select {
			case <-ctx.Done():
				resChan <- public.ExecStreamResult{Error: ctx.Err()}
				return
			case data, ok := <-e.stderr:
				if !ok {
					resChan <- public.ExecStreamResult{Error: io.EOF}
					return
				}
				shellMatch := PrimitiveCommandExecutorShellRegex.FindSubmatchIndex(data)
				if shellMatch != nil {
					if len(shellMatch) != 4 {
						resChan <- public.ExecStreamResult{Error: fmt.Errorf("unable to parse return code: %s", data)}
						return
					}
					cleanData := data[shellMatch[1]:]
					if len(cleanData) > 0 {
						resChan <- public.ExecStreamResult{Data: cleanData, DataType: public.ExecStreamDataTypeStderr}
					}
					returnCode, err := strconv.ParseInt(string(data[shellMatch[2]:shellMatch[3]]), 10, 32)
					if err != nil {
						resChan <- public.ExecStreamResult{Error: fmt.Errorf("unable to parse return code: %w", err)}
						return
					}
					if returnCode != 0 {
						resChan <- public.ExecStreamResult{Error: public.CommandExecutorError{Code: int(returnCode)}}
					}
					return
				}
				resChan <- public.ExecStreamResult{Data: data, DataType: public.ExecStreamDataTypeStderr}
			case data, ok := <-e.stdout:
				if !ok {
					resChan <- public.ExecStreamResult{Error: io.EOF}
					return
				}
				resChan <- public.ExecStreamResult{Data: data, DataType: public.ExecStreamDataTypeStdout}
			}
		}
	}()

	return resChan
}

func (e *PrimitiveCommandExecutor) Close() error {
	return e.stdinWriter.Close()
}
