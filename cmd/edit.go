package cmd

import (
	"fmt"

	"flopshot.io/dev/cli/api"
	"flopshot.io/dev/cli/edit"
	"github.com/manifoldco/promptui"
	"github.com/charmbracelet/huh"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
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

		p := huh.NewSelect[string]().
			Title("Select Type").
			Options(huh.NewOptions(allTypesFormatted...)...)

		p.Run()

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

		labels := make([]string, len(resp.Items))
		queryObjects := make([]*edit.EditType, len(resp.Items))

		// At this point we know there are items to render
		for i, t := range resp.Items {

			editType, _ := edit.GetType[edit.EditType](selection)

			// Copy values from the type returned into the edit type
			mapstructure.Decode(t, editType)

			labels[i] = (*editType).Label()
			queryObjects[i] = editType
		}

		prompt = promptui.Select{
			Label: "Select Item",
			Items: labels,
		}
		promptPos, _, _ = prompt.Run()

		selectedObj := queryObjects[promptPos]

		renderObject(selectedObj)
	},
}

func renderObject(obj *edit.EditType) {
	objFields, _ := edit.TypeFields(obj)

	renderEditFields(obj, objFields)
}

func renderEditFields(obj *edit.EditType, fields []edit.Field) {

	// Format the fields
	formattedFields := make([]string, len(fields))
	for i, v := range fields {
		formattedFields[i] = fmt.Sprintf("%s (%s): %s", v.Name, v.Type, v.Value)
	}

	// Append
	formattedFields = append([]string{"[BACK]"}, formattedFields...)
	formattedFields = append([]string{"[BACK]"}, formattedFields...)

	// Select a field to modify
	promptSelect := promptui.Select{
		Label: fmt.Sprintf("Edit: %s", (*obj).Label()),
		Items: formattedFields,
	}

	promptPos, _, err := promptSelect.Run()

	// It will likely be an error if they ctrl-c
	if err != nil {
		return
	}

	selectedField := &fields[promptPos]

	prompt := promptui.Prompt{
		Label: selectedField.Name,
		Default: selectedField.Value.String(),
		AllowEdit: true,
	}

	val, err := prompt.Run()
	if err != nil {
		return
	}

	edit.UpdateField(obj, selectedField, val)

	// Recurse to update another field if desired
	renderEditFields(obj, fields)
}

func init() {
	rootCmd.AddCommand(editCmd)
}
