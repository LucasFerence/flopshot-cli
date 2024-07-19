package cmd

import (
	"github.com/spf13/cobra"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Create or edit data",
	Run: func(cmd *cobra.Command, args []string) {


	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
