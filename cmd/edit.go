package cmd

import (
	"fmt"

	"flopshot.io/dev/cli/api"
	"flopshot.io/dev/cli/edit"
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

		// Convert all types to options
		selectOptions := make([]huh.Option[string], len(allTypes))
		caser := cases.Title(language.English)
		for i, v := range allTypes {
			selectOptions[i] = huh.NewOption(caser.String(v), v)
		}

		selectVal := ""
		shouldSearch := true
		err := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Affirmative("Search").
					Negative("Create").
					Value(&shouldSearch),
			),
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select Type").
					Options(selectOptions...).
					Value(&selectVal),
			),
		).Run()

		editObj, _ := edit.GetType[edit.EditType](selectVal)

		if shouldSearch {

		}

		if shouldSearch {
			resp := api.ListResponse[any]{}
			err = flopshotClient.QueryData(selectVal, &resp, []api.QueryParams{})

			if err != nil {
				fmt.Println(err)
				return
			}

			// Check to make sure items were returned
			if resp.Items == nil || len(resp.Items) == 0 {
				fmt.Println("No items found!")
				return
			}

			selectObjOpts := make([]huh.Option[*edit.EditType], len(resp.Items))

			// At this point we know there are items to render
			for i, t := range resp.Items {

				// Copy values from the type returned into the edit type
				mapstructure.Decode(t, editObj)

				selectObjOpts[i] = huh.NewOption((*editObj).Label(), editObj)
			}

			err = huh.NewSelect[*edit.EditType]().
				Title("Select Item").
				Options(selectObjOpts...).
				Value(&editObj).
				Run()

			if err != nil {
				fmt.Println(err)
				return
			}
		}

		shouldWrite := renderObject(editObj)
		if shouldWrite {
			flopshotClient.WriteData(selectVal, editObj)
		}
	},
}

func renderObject(obj *edit.EditType) bool {
	objFields, _ := edit.TypeFields(obj)

	return renderEditFields(obj, objFields)
}

func renderEditFields(obj *edit.EditType, fields []edit.Field) bool {

	textForms := make([]huh.Field, len(fields))
	fieldValues := make([]string, len(fields))
	for i, v := range fields {

		if v.RefType != "" {
			fieldValues[i] = fmt.Sprint(v.Value.FieldByName("Id"))
		} else {
			fieldValues[i] = fmt.Sprint(v.Value)
		}

		textForms[i] = huh.NewText().
			Title(fmt.Sprintf("%s (%s)", v.Name, v.Type)).
			Value(&fieldValues[i]).
			Lines(2)
	}

	shouldWrite := false
	_ = huh.NewForm(
		huh.NewGroup(
			textForms...,
		),
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("Write: %s?", (*obj).Label())).
				Value(&shouldWrite),
		),
	).Run()

	if shouldWrite {

		for i := 0; i < len(fields); i++ {
			// todo: this can be optimized to update all in one batch
			err := edit.UpdateField(obj, &fields[i], fieldValues[i])
			if err != nil {
				fmt.Printf("Could not update! Invalid Field. Error: [%s]\n", err)
				return false
			}
		}
	}

	return shouldWrite
}

func init() {
	rootCmd.AddCommand(editCmd)
}
