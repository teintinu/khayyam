package cmd

import (
	"errors"
	"os"
	"os/exec"
	"sync"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/spf13/cobra"
	"github.com/teintinu/monoclean/internal"
)

var DevOpts = internal.TestOptions{}

func init() {
	rootCmd.AddCommand(devCmd)
	internal.Logger.FlagDeclare(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Test libraries and run executables in watch mode",
	Long:  `Test libraries and run executables in watch mode`,
	Args:  cobra.RangeArgs(0, 0),
	RunE: func(cmd *cobra.Command, args []string) error {
		internal.Logger.FlagInit()
		repo := mustLoadRepository(true)
		if err := internal.CheckEngines(repo); err != nil {
			return err
		}

		wg := &sync.WaitGroup{}

		wg.Add(2)
		go devTestApps(repo, wg)
		go devRunExecutables(repo, wg)

		wg.Wait()
		return nil
	},
}

func devTestApps(repo *internal.Repository, wg *sync.WaitGroup) {
	err := internal.Test(repo, internal.TestOptions{
		Watch: true,
	})

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		os.Exit(exitErr.ExitCode())
	}
	wg.Done()
}

func devRunExecutables(repo *internal.Repository, wg *sync.WaitGroup) {
	for _, pkg := range repo.Packages {
		internal.BundleWithEsbuild(repo, pkg, &internal.BuildOpts{
			Target: api.ESNext,
			Minify: false,
			Mode:   internal.WatchAndRun,
		})
	}
	wg.Done()
}
