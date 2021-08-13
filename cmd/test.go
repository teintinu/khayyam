package cmd

import (
	"errors"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/teintinu/monoclean/internal"
)

var testOpts = internal.TestOptions{
	Colors: true,
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().BoolVar(&testOpts.Watch, "watch", false, "re-runs tests when source files change")
	testCmd.Flags().BoolVar(&testOpts.Coverage, "coverage", false, "Run testes with coverage")
}

var testCmd = &cobra.Command{
	Use:                   "test [flags] <script> [args...]",
	Short:                 "Run tests.",
	Long:                  `Run tests using jest.`,
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository(true)
		if err := internal.CheckEngines(repo); err != nil {
			return err
		}

		jestCmd := internal.CreateJestCommand(repo, testOpts)
		jestCmd.Stdin = os.Stdin
		jestCmd.Stdout = os.Stdout
		jestCmd.Stderr = os.Stderr
		err := jestCmd.Run()
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.ExitCode())
		}
		return err
	},
}
