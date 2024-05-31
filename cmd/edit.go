package cmd

import (
	"fmt"

	"flopshot.io/dev/cli/api"
	"flopshot.io/dev/cli/edit"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"github.com/mitchellh/mapstructure"
)

// editCmd represents the edit command
var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {

		loggedIn, _ := flopshotClient.IsAuthenticated()
		if !loggedIn {
			fmt.Println("Must be logged in!")
			return
		}

		// Get all the supported types
		allTypes := edit.AllTypes()

		// Format all types
		allTypesFormatted := make([]string, len(allTypes))
		caser := cases.Title(language.English)
		for i, v := range allTypes {
			allTypesFormatted[i] = caser.String(v)
		}

		// Generate a prompt for displayng types
		prompt := promptui.Select{
			Label: "Select Type",
			Items: allTypesFormatted,
		}

		promptPos, _, err := prompt.Run()
		if err != nil {
			fmt.Println(err)
			return
		}

		// Get the raw selected type based on positioning of arguments
		selection := allTypes[promptPos]
		editType, _ := edit.FindType[edit.EditType](selection)

		resp := api.ListResponse[any]{}
		err = flopshotClient.QueryData(selection, &resp, []api.QueryParams{})

		if err != nil {
			fmt.Println(err)
			return
		}

		// Check to make sure items were returned
		if resp.Items == nil || len(resp.Items) == 0 {
			fmt.Println("No items found!")
			return
		}

		// At this point we know there are items to render
		for _, t := range resp.Items {

			// Decode any into the type that was returned
			mapstructure.Decode(t, editType)
			fmt.Println((*editType).Label())
		}
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
