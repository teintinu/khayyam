package cmd

import (
	"runtime/debug"

	"github.com/spf13/cobra"
	"github.com/teintinu/khayyam/internal"
)

var detailed bool

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().BoolVar(&detailed, "detailed", false, "Prints version information for all dependencies too.")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of khayyam.",
	Long:  "Print the version of khayyam.",
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo, ok := debug.ReadBuildInfo()
		if !ok {
			panic("debug.ReadBuildInfo() failed")
		}
		printInfo := func(mod debug.Module) {
			if detailed {
				internal.Logger.Info(mod.Path, mod.Version)
			} else {
				internal.Logger.Info(mod.Version)
			}
		}
		printInfo(buildInfo.Main)
		if detailed {
			for _, dep := range buildInfo.Deps {
				internal.Logger.Info(dep.Path, dep.Version)
			}
		}
	},
}
