package cmd

import (
	"github.com/spf13/cobra"
	"github.com/teintinu/khayyam/internal"
)

var initOpts = internal.InitOptions{}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&initOpts.CleanArchitecture, "clean-rchitecture", false, "Use Clean Architecture suggestion")
}

var initCmd = &cobra.Command{
	Use:                   "init [flags] <script> [args...]",
	Short:                 "Init.",
	Long:                  `Initialized the workspace.`,
	DisableFlagsInUseLine: true,
	SilenceErrors:         true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := internal.Init(initOpts); err != nil {
			return err
		}
		repo := mustLoadRepository(false)
		if err := internal.CheckEngines(repo); err != nil {
			return err
		}
		return internal.InstallDependencies(repo, internal.InstallDependenciesOptions{})
	},
}
