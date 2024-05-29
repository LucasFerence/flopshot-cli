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

		resp := flopshotClient.RegisterIdReq(edit.UserType)

		user := edit.User{
			Id:    resp.Id,
			Email: "test@gmail.com",
			Name:  "Lucas Test",
		}

		flopshotClient.WriteData(edit.UserType, user)

		listResp := api.ListResponse[edit.User]{}
		flopshotClient.QueryData(
			edit.UserType,
			&listResp,
			[]api.QueryParams{{K: "email", V: "ference.lucas@gmail.com"}},
		)

		fmt.Println(listResp)
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
