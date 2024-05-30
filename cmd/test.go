/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"flopshot.io/dev/cli/api"
	"flopshot.io/dev/cli/edit"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		// This will be later passed through the command as an option/sub-command or something
		const cmdTypeName = "user"

		listResp := api.ListResponse[any]{}
		flopshotClient.QueryData(cmdTypeName, &listResp, []api.QueryParams{{K: "p", V: "0"}})
		fmt.Println(listResp)

		editType := edit.FindType[any](cmdTypeName)
		fields := edit.TypeFields(editType)

		for _, f := range fields {
			fmt.Println(f)
		}

		// This will actually write the data correctly
		// flopshotClient.WriteData("user", &editType)
	},
}

func init() {
	rootCmd.AddCommand(testCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
