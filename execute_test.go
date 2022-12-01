package cobrax

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"testing"
)

func TestExec(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		cmd := &cobra.Command{
			Run: func(cmd *cobra.Command, args []string) {
				in, _ := io.ReadAll(cmd.InOrStdin())
				_, _ = fmt.Fprint(cmd.OutOrStdout(), string(in))
				_, _ = fmt.Fprint(cmd.ErrOrStderr(), args[0])
			},
		}
		stdOut, stdErr, err := Exec(cmd, bytes.NewBufferString("foo"), "bar")
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if stdOut != "foo" {
			t.Fatalf("unexpected stdOut: %q", stdOut)
		}
		if stdErr != "bar" {
			t.Fatalf("unexpected stdErr: %q", stdErr)
		}
	})

	t.Run("error", func(t *testing.T) {
		cmd := &cobra.Command{
			RunE: func(cmd *cobra.Command, args []string) error {
				_, _ = fmt.Fprint(cmd.OutOrStdout(), args[0])
				_, _ = fmt.Fprint(cmd.ErrOrStderr(), args[1])
				return FailSilently(cmd)
			},
		}
		stdOut, stdErr, err := Exec(cmd, nil, "foo", "bar")
		if err == nil {
			t.Fatalf("expected error")
		}
		if !errors.Is(err, ErrNoPrintButFail) {
			t.Fatalf("unexpected error: %v", err)
		}
		if stdOut != "foo" {
			t.Fatalf("unexpected stdOut: %q", stdOut)
		}
		if stdErr != "bar" {
			t.Fatalf("unexpected stdErr: %q", stdErr)
		}
	})
}
