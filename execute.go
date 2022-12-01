// Copyright © 2022 Patrik Neu
// SPDX-License-Identifier: Apache-2.0

package cobrax

import (
	"bytes"
	"context"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"io"
)

// ExecBackgroundCtx runs the cobra command in the background. This function can also be used outside of tests.
// Pass nil for stdIn, stdOur, or stdErr to use os.Std*.
func ExecBackgroundCtx(ctx context.Context, cmd *cobra.Command, stdIn io.Reader, stdOut, stdErr io.Writer, args ...string) *errgroup.Group {
	cmd.SetIn(stdIn)
	cmd.SetOut(stdOut)
	cmd.SetErr(stdErr)

	if args == nil {
		args = []string{}
	}
	cmd.SetArgs(args)

	eg := &errgroup.Group{}
	eg.Go(func() error {
		defer cmd.SetIn(nil)
		return cmd.ExecuteContext(ctx)
	})

	return eg
}

// Exec runs the provided cobra command with the given reader as STD_IN and the given args. This function can also be used outside of tests.
func Exec(cmd *cobra.Command, stdIn io.Reader, args ...string) (stdOut string, stdErr string, err error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return ExecCtx(ctx, cmd, stdIn, args...)
}

// ExecCtx is the same as [Exec] but with a user-supplied context. This function can also be used outside of tests.
func ExecCtx(ctx context.Context, cmd *cobra.Command, stdIn io.Reader, args ...string) (stdOut string, stdErr string, err error) {
	stdOutB, stdErrB := &bytes.Buffer{}, &bytes.Buffer{}
	err = ExecBackgroundCtx(ctx, cmd, stdIn, stdOutB, stdErrB, args...).Wait()

	return stdOutB.String(), stdErrB.String(), err
}

// CommandExecutor is a struct that can be used to execute a cobra command multiple times.
type CommandExecutor struct {
	New            func() *cobra.Command
	Ctx            context.Context
	PersistentArgs []string
}

// Exec runs the cobra command with the given reader as STD_IN and the given args appended to the persistent args.
func (c *CommandExecutor) Exec(stdin io.Reader, args ...string) (stdOut string, stdErr string, err error) {
	return ExecCtx(c.Ctx, c.New(), stdin, append(c.PersistentArgs, args...)...)
}

// ExecBackground runs the cobra command in the background with the given reader as STD_IN and the given args appended to the persistent args.
func (c *CommandExecutor) ExecBackground(stdin io.Reader, stdOut, stdErr io.Writer, args ...string) *errgroup.Group {
	return ExecBackgroundCtx(c.Ctx, c.New(), stdin, stdOut, stdErr, append(c.PersistentArgs, args...)...)
}
