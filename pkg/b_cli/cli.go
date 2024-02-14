package cli

import (
	"github.com/spf13/cobra"
	"os"
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
	cockroachCmd.AddCommand(
		startSingleNodeCmd,
	)
}

// Main is the entry point for the cli, with a single line calling it intended
// to be the body of an action package main `main` func elsewhere. It is
// abstracted for reuse by duplicated `main` funcs in different distributions.
func Main() {
	if len(os.Args) == 1 {
		//os.Args = append(os.Args, "help")
		os.Args = append(os.Args, "start-single-node")
	}

	err := doMain()
	if err != nil {
		os.Exit(1)
	}

}

func doMain() error {
	return Run(os.Args[1:])
}

func Run(args []string) error {
	cockroachCmd.SetArgs(args)
	return cockroachCmd.Execute()
}
