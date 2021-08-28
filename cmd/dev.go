package cmd

import (
	"github.com/spf13/cobra"
	"github.com/teintinu/khayyam/internal"
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

		return internal.Dev(repo)
	},
}
