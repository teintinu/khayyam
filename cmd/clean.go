package cmd

import (
	"github.com/spf13/cobra"
	"github.com/teintinu/khayyam/internal"
)

func init() {
	rootCmd.AddCommand(cleanCmd)
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "Removes build output.",
	Long:  `Removes build output.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := mustLoadRepository(false)
		return internal.CleanRepository(repo)
	},
}
