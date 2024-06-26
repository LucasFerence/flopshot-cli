package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"flopshot.io/dev/cli/api"
	"flopshot.io/dev/cli/edit"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/icza/gox/osx"
	"github.com/spf13/cobra"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "A brief description of your command",

	Run: func(cmd *cobra.Command, args []string) {

		// Short circuit and exit if they are already logged in
		loggedIn, _ := flopshotClient.IsAuthenticated()
		if loggedIn {
			fmt.Println("Already logged in!")
			return
		}

		// Get the device code response
		deviceCodeRespose := deviceCodeResponse{}
		err := execDeviceCodeRequest(&deviceCodeRespose)

		if err != nil {
			fmt.Println(err)
			return
		}

		shouldOpenBrowser := true
		err = huh.NewConfirm().
			Title("Open Browser?").
			Value(&shouldOpenBrowser).
			Run()

		if err != nil || !shouldOpenBrowser {
			return
		}

		// Open URL in default browser
		osx.OpenDefault(deviceCodeRespose.VerificationUriComplete)

		tokenResponse := tokenResponse{}
		awaitTokenAction := func() {

			// Sleep for 5 seconds initially, then check for a token response
			// every 5 seconds to see if they're authorized
			time.Sleep(5 * time.Second)
			for i := 0; i < 15; i++ {

				code, err := execTokenRequest(deviceCodeRespose.DeviceCode, &tokenResponse)
				if err != nil {
					fmt.Println(err)
					return
				}

				if code == 200 {
					break
				}

				// Sleep for one second before trying to get a token
				time.Sleep(5 * time.Second)
			}
		}

		err = spinner.New().
			Title(fmt.Sprintf("Confirm Code: %s", deviceCodeRespose.UserCode)).
			Action(awaitTokenAction).
			Run()

		if err != nil || tokenResponse.AccessToken == "" {
			fmt.Println("Login failed. Please try again")
			return
		}

		flopshotClient.InitAuth(tokenResponse.AccessToken)

		email := ""
		err = huh.NewInput().
			Value(&email).
			Title("Enter email").
			Description("Check for existing user").
			Validate(func(s string) error {
				if len(s) == 0 {
					return errors.New("Must input an email!")
				}

				return nil
			}).
			Run()

		if err != nil {
			return
		}

		// Query for existing users
		userResp := api.ListResponse[edit.User]{}
		err = flopshotClient.QueryData(
			edit.UserType,
			&userResp,
			[]api.QueryParams{{K: "email", V: email}},
		)

		if len(userResp.Items) > 0 {
			fmt.Println("User already exists!")
		} else {

			fmt.Println("No user. Creating one automatically...")

			newUser := edit.User{
				Email: email,
			}
			flopshotClient.WriteData(edit.UserType, newUser)
		}
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loginCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loginCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// --- Command API Utility ---

// Response model for initial device code request
type deviceCodeResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationUri         string `json:"verification_uri"`
	ExpiresIn               int    `json:"expires_in"`
	Interval                int    `json:"interval"`
	VerificationUriComplete string `json:"verification_uri_complete"`
}

// Execute a request to retrieve the device code init data
// Pass a pointer of the response data to be populated
func execDeviceCodeRequest(data *deviceCodeResponse) error {
	form := url.Values{}
	form.Add("client_id", "Ap9LIyJxGcc0vVvisrLiaLsCVLWDahqv")
	form.Add("audience", "http://localhost:5050")

	req, err := http.NewRequest(
		"POST",
		"https://dev-fkh-ll2p.us.auth0.com/oauth/device/code",
		strings.NewReader(form.Encode()),
	)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	// Execute and ignore raw response
	_, err = flopshotClient.ExecR(req, &data)

	// Return the err (it will be nil if successful)
	return err
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

func execTokenRequest(deviceCode string, data *tokenResponse) (int, error) {

	form := url.Values{}
	form.Add("grant_type", "urn:ietf:params:oauth:grant-type:device_code")
	form.Add("device_code", deviceCode)
	form.Add("client_id", "Ap9LIyJxGcc0vVvisrLiaLsCVLWDahqv")

	req, err := http.NewRequest(
		"POST",
		"https://dev-fkh-ll2p.us.auth0.com/oauth/token",
		strings.NewReader(form.Encode()),
	)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	resp, err := flopshotClient.ExecR(req, &data)
	if err != nil {
		fmt.Println(err)
		return -1, err
	}

	return resp.RawResponse.StatusCode, nil
}
