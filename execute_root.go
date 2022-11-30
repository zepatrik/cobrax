package cobrax

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// ErrNoPrintButFail is returned to detect a failure state that was already reported to the user in some way
var ErrNoPrintButFail = errors.New("this error should never be printed")

// FailSilently is supposed to be used within a commands RunE function.
// It silences cobras default error handling and returns the ErrNoPrintButFail error.
func FailSilently(cmd *cobra.Command) error {
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	return ErrNoPrintButFail
}

type ExitCodeCarrier interface {
	error
	ExitCode() int
}

type exitCodeError struct {
	error
	code int
}

func WithExitCode(err error, exitCode int) error {
	return &exitCodeError{err, exitCode}
}

func (e *exitCodeError) ExitCode() int {
	return e.code
}

func (e *exitCodeError) Unwrap() error {
	return e.error
}

// ExecuteRootCommand executes the given command (usually the root command) and exits the process with a non-zero exit code if an error occurs.
// The error is only printed if it is not ErrNoPrintButFail, because in that case we assume the error was already printed.
// To customize the status code you can add an error to the error chain that implements the ExitCodeCarrier interface.
func ExecuteRootCommand(cmd *cobra.Command) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	ExecuteRootCommandContext(ctx, cmd)
}

// ExecuteRootCommandContext is the same as ExecuteRootCommand but with a user-supplied context.
func ExecuteRootCommandContext(ctx context.Context, cmd *cobra.Command) {
	if err := cmd.ExecuteContext(ctx); err != nil {
		if !errors.Is(err, ErrNoPrintButFail) {
			fmt.Println(err)
		}
		if se := ExitCodeCarrier(nil); errors.As(err, &se) {
			osExit(se.ExitCode())
			return
		}
		osExit(1)
		return
	}
	osExit(0)
}

// osExit is a variable to allow mocking os.Exit in tests
var osExit = os.Exit
