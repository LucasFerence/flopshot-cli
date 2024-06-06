package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"flopshot.io/dev/cli/api"
	"flopshot.io/dev/cli/edit"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

// bookCmd represents the book command
var bookCmd = &cobra.Command{
	Use:   "book",
	Short: "Book a tee time round",
	Run: func(cmd *cobra.Command, args []string) {

		// First get all the courses

		allCourseResp := api.ListResponse[edit.Course]{}
		flopshotClient.QueryData(edit.CourseType, &allCourseResp, []api.QueryParams{})

		// Create all course options
		courses := allCourseResp.Items
		courseOpts := make([]huh.Option[string], len(courses))
		for i, c := range courses {
			courseOpts[i] = huh.NewOption(c.Label(), c.Id)
		}

		var details struct {
			Course     string `json:"courseId"`
			Email      string `json:"userEmail"`
			Date       string `json:"roundDate"`
			NumPlayers int    `json:"playerCount"`
			Time       string `json:"preferredTime"`
		}

		err := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[string]().
					Title("Select course").
					Options(courseOpts...).
					Value(&details.Course),
			),
			huh.NewGroup(
				huh.NewInput().
					Title("Enter email").
					Value(&details.Email).
					Validate(func(s string) error {
						if len(s) == 0 {
							return errors.New("Must input an email!")
						}

						return nil
					}),
				huh.NewSelect[int]().
					Title("Number of players").
					Options(huh.NewOptions(1, 2, 3, 4)...).
					Value(&details.NumPlayers),
				huh.NewInput().
					Title("Date").
					Description("2012-04-23T18:25:43.511Z").
					Value(&details.Date),
				huh.NewInput().
					Title("Preferred Time").
					Description("2012-04-23T18:25:43.511Z").
					Value(&details.Time),
			),
		).Run()

		if err != nil {
			return
		}

		execBookRequest(&flopshotClient, &details)
	},
}

func execBookRequest(client *api.FlopshotClient, data any) error {

	val, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/book", api.ClientUrl),
		bytes.NewReader(val),
	)
	if err != nil {
		return err
	}

	_, err = client.Exec(req)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(bookCmd)
}
