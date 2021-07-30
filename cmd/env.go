package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Dump environment information.",
	Long: `Validates and prints information about the current environment.

Returns a non-zero status code if there environment fails any checks.
The printed output of this command is intended for human consumption and
not (yet?) intended to be parsed.`,
	Run: func(cmd *cobra.Command, args []string) {
		println("TODO")
		// repo := mustLoadRepository()
		// env, err := internal.AnalyzeEnvironment(repo)
		// if err != nil {
		// 	internal.Log.Error(err.Error())
		// 	os.Exit(1)
		// }
		// internal.DumpEnvironment(env)
		// if !env.OK {
		// 	os.Exit(1)
		// }
	},
}
