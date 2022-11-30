package cobrax

import (
	"context"
	"errors"
	"github.com/spf13/cobra"
)

// TestingT is an interface excerpt of testing.TB
type TestingT interface {
	Errorf(format string, args ...interface{})
	FailNow()
	Helper()
}

// ExecNoErr is a test helper that assumes a successful run from Exec.
func ExecNoErr(t TestingT, cmd *cobra.Command, args ...string) (stdOut string) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return ExecNoErrCtx(ctx, t, cmd, args...)
}

// ExecNoErrCtx is the same as ExecNoErr but with a user-supplied context.
func ExecNoErrCtx(ctx context.Context, t TestingT, cmd *cobra.Command, args ...string) string {
	t.Helper()
	stdOut, stdErr, err := ExecCtx(ctx, cmd, nil, args...)
	if err != nil {
		t.Errorf("Expected no error, got %#v", err)
		t.Errorf("std_out: %q, std_err: %q", stdOut, stdErr)
		t.FailNow()
	}
	if len(stdErr) != 0 {
		t.Errorf("Expected empty stdErr, got %q", stdErr)
		t.Errorf("std_out: %q", stdOut)
		t.FailNow()
	}
	return stdOut
}

// ExecExpectedErr is a test helper that assumes a failing run from Exec returning ErrNoPrintButFail
func ExecExpectedErr(t TestingT, cmd *cobra.Command, args ...string) (stdErr string) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return ExecExpectedErrCtx(ctx, t, cmd, args...)
}

// ExecExpectedErrCtx is the same as ExecExpectedErr but with a user-supplied context.
func ExecExpectedErrCtx(ctx context.Context, t TestingT, cmd *cobra.Command, args ...string) (stdErr string) {
	t.Helper()

	stdOut, stdErr, err := ExecCtx(ctx, cmd, nil, args...)
	if !errors.Is(err, ErrNoPrintButFail) {
		t.Errorf("Expected error %#v, got %#v", ErrNoPrintButFail, err)
		t.Errorf("std_out: %q, std_err: %q", stdOut, stdErr)
		t.FailNow()
	}
	if len(stdOut) != 0 {
		t.Errorf("Expected empty stdOut, got %q", stdOut)
		t.Errorf("std_err: %q", stdErr)
		t.FailNow()
	}
	return stdErr
}

// ExecNoErr is a test helper that assumes a successful run. The args are appended to the persistent args.
func (c *CommandExecutor) ExecNoErr(t TestingT, args ...string) string {
	t.Helper()
	return ExecNoErrCtx(c.Ctx, t, c.New(), append(c.PersistentArgs, args...)...)
}

// ExecExpectedErr is a test helper that assumes a failing run returning ErrNoPrintButFail. The args are appended to the persistent args.
func (c *CommandExecutor) ExecExpectedErr(t TestingT, args ...string) string {
	t.Helper()
	return ExecExpectedErrCtx(c.Ctx, t, c.New(), append(c.PersistentArgs, args...)...)
}
