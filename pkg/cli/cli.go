package cli

import (
	"github.com/spf13/cobra"
	"os"
)

var cockroachCmd = &cobra.Command{
	Use:   "cockroach [command] (flags)",
	Short: "CockroachDB command-line interface and server",
	// TODO(cdo): Add a pointer to the docs in Long.
	Long: `CockroachDB command-line interface and server.`,
	// Disable automatic printing of usage information whenever an error
	// occurs. Many errors are not the result of a bad command invocation,
	// e.g. attempting to start a node on an in-use port, and printing the
	// usage information in these cases obscures the cause of the error.
	// Commands should manually print usage information when the error is,
	// in fact, a result of a bad invocation, e.g. too many arguments.
	SilenceUsage: true,
	// Disable automatic printing of the error. We want to also print
	// details and hints, which cobra does not do for us. Instead
	// we do the printing in Main().
	SilenceErrors: true,
	// Version causes cobra to automatically support a --version flag
	// that reports this string.
	Version: "details:\n" + "v1" +
		"\n(use '" + os.Args[0] + " version --build-tag' to display only the build tag)",
	// Prevent cobra from auto-generating a completions command,
	// since we provide our own.
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
