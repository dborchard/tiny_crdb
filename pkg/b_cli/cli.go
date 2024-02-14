package cli

import (
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var cockroachCmd = &cobra.Command{
	Use:           "cockroach [command] (flags)",
	Short:         "CockroachDB command-line interface and server",
	SilenceUsage:  true,
	SilenceErrors: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func init() {
	cobra.EnableCommandSorting = false

	// Set an error function for flag parsing which prints the usage message.
	cockroachCmd.SetFlagErrorFunc(func(c *cobra.Command, err error) error {
		return c.Usage()
	})

	cockroachCmd.AddCommand(
		startSingleNodeCmd,
	)
}

// Main is the entry point for the cli, with a single line calling it intended
// to be the body of an action package main `main` func elsewhere. It is
// abstracted for reuse by duplicated `main` funcs in different distributions.
func Main() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "help")
	}

	// We ignore the error in this lookup, because
	// we want cobra to handle lookup errors with a verbose
	// help message in Run() below.
	cmd, _, _ := cockroachCmd.Find(os.Args[1:])

	cmdName := commandName(cmd)
	err := doMain(cmd, cmdName)
	if err != nil {
		os.Exit(1)
	}

}

func doMain(cmd *cobra.Command, cmdName string) error {
	return Run(os.Args[1:])
}

// Run ...
func Run(args []string) error {
	cockroachCmd.SetArgs(args)
	return cockroachCmd.Execute()
}

// commandName computes the name of the command that args would invoke. For
// example, the full name of "cockroach debug zip" is "debug zip". If args
// specify a nonexistent command, commandName returns "cockroach".
func commandName(cmd *cobra.Command) string {
	rootName := cockroachCmd.CommandPath()
	if cmd != nil {
		return strings.TrimPrefix(cmd.CommandPath(), rootName+" ")
	}
	return rootName
}
