package main

import (
	"context"
	"os"

	"github.com/charmbracelet/fang"
	"github.com/orangekame3/qasmtools/cmd/qasm/commands"
	"github.com/orangekame3/vercheck"
	"github.com/spf13/cobra"
)

func main() {
	// Check for newer versions
	vercheck.Check(vercheck.Options{
		CurrentVersion: Version,
		RepoOwner:      "orangekame3",
		RepoName:       "qasmtools",
	})

	if err := fang.Execute(context.Background(), rootCmd, fang.WithVersion(GetVersion())); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "qasm",
	Short: "QASM CLI tools",
	Long:  `A collection of tools for working with QASM files.`,
}

func init() {
	// Add all subcommands
	rootCmd.AddCommand(
		commands.NewFormatCommand(),
		commands.NewHighlightCommand(),
		commands.NewLintCommand(),
		commands.NewParseCommand(),
		commands.NewBenchmarkCommand(),
	)
}
