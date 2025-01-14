package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/teintinu/khayyam/internal"
)

func init() {
	rootCmd.AddCommand(runCmd)
	internal.Logger.LogFlagDeclare(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [flags] <executable...>",
	Short: "Build and run",
	Long:  `Builds and runs one or more applications. If not specified will run all applications in repository`,
	//Args:  cobra.MinimumNArgs(1),
	// DisableFlagsInUseLine: true,
	// SilenceErrors:         true,
	RunE: func(cmd *cobra.Command, args []string) error {
		internal.Logger.LogFlagInit()
		repo := mustLoadRepository(true)
		var packagesToRun []string = []string{}

		if err := internal.CheckEngines(repo); err != nil {
			return err
		}

		if len(args) == 0 {
			for pkgName, pkg := range repo.Packages {
				if pkg.Executable {
					packagesToRun = append(packagesToRun, pkgName)
				}
			}
		} else {
			for _, pkgName := range args {
				pkg := repo.Packages[pkgName]
				if pkg == nil {
					return errors.New("no such package: " + pkgName)
				}
				packagesToRun = append(packagesToRun, pkgName)
			}
		}
		return internal.Run(repo, packagesToRun, &internal.RunOpts{
			Watch: false,
		})
	},
}
