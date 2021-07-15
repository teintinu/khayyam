package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/deref/uni/internal"
	"github.com/spf13/cobra"
)

var testOpts = internal.TestOptions{}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().BoolVar(&testOpts.Watch, "watch", false, "re-runs command when source files change")
}

var testCmd = &cobra.Command{
	Use:   "test [flags]",
	Short: "Build and runs tests with jest.",
	Long: `Builds and runs tests with jest.

TODO	
Entrypoint files are modules that export a "main" function, which will be called
with the given 'args' as positional parameters.

Main functions may return an integer status code, and the return value will
will be awaited.  If no status code is returned, the default is 0.

After awaiting a return value, the process will be terminated immediately.  Any
pending events will not be executed; main is responsible for graceful shutdown.

Unhandled exceptions and promise rejections will be logged to stderr and the
process will immediately exit with status code 1.

Example:
TODO
export const main = async (...args: string[]) => {
  console.log("see uni run");
  return 0; // Return an exit code (optional).
}
`,
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
		if err := internal.CheckEngines(repo); err != nil {
			return err
		}

		var err = internal.Test(repo, testOpts)
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		return err
	},
}
