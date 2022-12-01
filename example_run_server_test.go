package cobrax_test

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/zepatrik/cobrax"
)

// NewKetoCmd returns the whole tree of Keto commands, including the server and client commands.
// This would be imported from [github.com/ory/keto/cmd.NewRootCmd].
func NewKetoCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "keto",
		Short: "Global and consistent permission and authorization server by Ory",
	}
}

// ExampleExecDependency shows how a cobra command can be executed like [os/exec.Command], but without the OS overhead.
func Example_execDependency() {
	serverCtx, serverCancel := context.WithCancel(context.Background())
	defer serverCancel()

	eg := cobrax.ExecBackgroundCtx(serverCtx, NewKetoCmd(), nil, nil, nil, "serve", "--config", "keto.yml")
	defer func() {
		fmt.Printf("Keto server exited with error: %v", eg.Wait())
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmd := cobrax.CommandExecutor{
		New:            NewKetoCmd,
		Ctx:            ctx,
		PersistentArgs: []string{"--read-remote", "localhost:4466"},
	}

	// Blocks until the command exits.
	statusCtx, statusCancel := context.WithTimeout(ctx, 10*time.Second)
	defer statusCancel()
	_, _, _ = cmd.ExecCtx(statusCtx, nil, "status", "--block")

	stdOut, _, err := cmd.Exec(nil, "check", "relation-tuple", "--namespace", "default", "--object", "article:1", "--relation", "view", "--subject", "user:1")
	if err != nil {
		fmt.Printf("Keto client gave unexpected error: %v", err)
		return
	}
	fmt.Printf("Keto client gave output: %s", stdOut)

	stdOut, _, err = cmd.Exec(nil, "check", "relation-tuple", "--namespace", "default", "--object", "article:1", "--relation", "view", "--subject", "user:2")
	// ...and so on
}
