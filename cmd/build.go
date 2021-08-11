package cmd

import (
	"github.com/evanw/esbuild/pkg/api"
	"github.com/spf13/cobra"
	"github.com/teintinu/monoclean/internal"
)

func init() {
	rootCmd.AddCommand(buildCmd)
	internal.Logger.FlagDeclare(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Builds workspace",
	Long:  "Builds all packages on workspace",
	Args:  cobra.RangeArgs(0, 0),
	RunE: func(cmd *cobra.Command, args []string) error {
		internal.Logger.FlagInit()
		repo := mustLoadRepository(true)
		if err := internal.CheckEngines(repo); err != nil {
			return err
		}
		return internal.BuildWorkspace(repo, &internal.BuildOpts{
			Target: api.ES5,
			Minify: true,
			Mode:   internal.BuildOnly,
		})
	},
}
