package cmd

import (
	"github.com/spf13/cobra"
	"github.com/teintinu/khayyam/internal"
)

var DevOpts = internal.TestOptions{}

func init() {
	rootCmd.AddCommand(devCmd)
	internal.Logger.LogFlagDeclare(devCmd)
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Test libraries and run applications in watch mode",
	Long:  `Test libraries and run applications in watch mode`,
	Args:  cobra.RangeArgs(0, 0),
	RunE: func(cmd *cobra.Command, args []string) error {
		internal.Logger.LogFlagInit()
		repo := mustLoadRepository(true)
		if err := internal.CheckEngines(repo); err != nil {
			return err
		}

		return internal.Dev(repo)
	},
}
