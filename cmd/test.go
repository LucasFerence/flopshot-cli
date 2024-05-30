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

		// Query for all users (could provide pagination in query params if we wanted to)
		listResp := api.ListResponse[any]{}
		flopshotClient.QueryData(cmdTypeName, &listResp, []api.QueryParams{})
		fmt.Println(listResp)

		// Find the type based on the name identifier provided
		// By getting the fields we can display the type in a form
		editType, err := edit.FindType[any](cmdTypeName)
		if err != nil {
			fmt.Println(err)
			return
		}

		fields, err := edit.TypeFields(editType)
		if err != nil {
			fmt.Println(err)
			return
		}

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
