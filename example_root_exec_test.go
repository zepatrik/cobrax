package cobrax_test

import (
	"bufio"
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/zepatrik/cobrax"
)

// NewRootCmd creates a new root greetme command. This function takes care of all the flag initialization and registering all subcommands.
func NewRootCmd() *cobra.Command {
	var errorExitCode int

	cmd := &cobra.Command{
		Use:   "greetme",
		Short: "This is a friendly program to greet you.",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.ErrOrStderr(), "Hello, what is your name?")
			name, err := bufio.NewReader(cmd.InOrStdin()).ReadString('\n')
			if err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "Error reading name: %v", err)
				return cobrax.WithExitCode(cobrax.FailSilently(cmd), errorExitCode)
			}
			fmt.Fprintf(cmd.ErrOrStderr(), "Hello %s\n", name)
			return nil
		},
	}
	cmd.Flags().IntVar(&errorExitCode, "exit-code", 1, "Exit code to return on error")
	cmd.AddCommand(NewVersionCmd(""))

	return cmd
}

// NewVersionCmd returns a new version subcommand. This could be in a totally different package, and reused by multiple projects.
// The version can be passed as a string. If the version is empty, the version embedded by [runtime/debug.ReadBuildInfo] is used.
func NewVersionCmd(version string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Prints the version of this program.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if version == "" {
				bi, ok := debug.ReadBuildInfo()
				if !ok {
					fmt.Fprintln(cmd.ErrOrStderr(), "No version information available.")
					return cobrax.FailSilently(cmd)
				}
				version = bi.Main.Version
			}
			fmt.Fprintf(cmd.OutOrStdout(), "Version: %s", version)
			return nil
		},
	}
	cmd.AddCommand(newVersionSubCmd())

	return cmd
}

// NewVersionSubCmd returns a new version subcommand. This would be in the same package as the version command.
// It could depend on the version command as a parent, e.g. to access persistent flags, and would therefore not be exported.
func newVersionSubCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "deeply-nested",
		Short: "This commmand is an example for a deeply nested subcommand.",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Do stuff
			return nil
		},
	}
}

// ExampleExecuteRoot shows how to use the helpers for the [spf13/cobra.Command.RunE] and main function.
func Example_executeRoot() {
	cobrax.ExecuteRootCommand(NewRootCmd())
}
