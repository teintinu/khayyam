package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/teintinu/monoclean/internal"
)

func init() {
	rootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:                   "run [flags] <executable...>",
	Short:                 "Build and run",
	Long:                  `Builds and runs one or more executables. If not specified will run all executables in repository`,
	Args:                  cobra.MinimumNArgs(1),
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository()
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
		return internal.Run(repo, packagesToRun)
	},
}
