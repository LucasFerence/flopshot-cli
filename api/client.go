package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/99designs/keyring"
)

const ClientUrl = "http://localhost:5050"

// Constants for keyring identifiers
const keyringService = "flopshot"
const loginTokenKey = "loginToken"

func openRing() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: keyringService,

		// This will still prompt the user for allowance, but will remember the change
		KeychainTrustApplication: true,
	})
}

type FlopshotClient struct {

	// Underlying client for executing requests
	HttpClient *http.Client

	// calculated token, might be nil
	Token string
}

func NewFlopshotClient() FlopshotClient {

	client := FlopshotClient{
		HttpClient: &http.Client{},
	}

	return client
}

func (client *FlopshotClient) InitAuth(token string) {

	ring, err := openRing()

	if err != nil {
		fmt.Println(err)
		return
	}

	err = ring.Set(keyring.Item{
		Key:  loginTokenKey,
		Data: []byte(token),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	client.Token = token
}

func (client *FlopshotClient) RemoveAuth() {

	ring, err := openRing()

	if err != nil {
		fmt.Println(err)
		return
	}

	ring.Remove(loginTokenKey)
}

func (client *FlopshotClient) IsAuthenticated() (bool, string) {

	ring, err := openRing()
	if err != nil {
		fmt.Println(err)
		return false, ""
	}

	keys, _ := ring.Keys()
	tokenExists := slices.Contains(keys, loginTokenKey)

	if tokenExists {
		val, _ := ring.Get(loginTokenKey)
		return tokenExists, string(val.Data)
	}

	return false, ""
}

func (client *FlopshotClient) Exec(req *http.Request) (*BufferedResponse, error) {

	isAuth, token := client.IsAuthenticated()

	// If authenticated, add headers for auth
	if isAuth {
		req.Header.Add("content-type", "application/json")
		req.Header.Add("authorization", "Bearer "+token)
	}

	resp, err := client.HttpClient.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Do this at the end of the execution function
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	bufferedResp := &BufferedResponse{
		RawResponse: resp,
		Body:        body,
	}

	return bufferedResp, nil
}

func (client *FlopshotClient) ExecR(req *http.Request, respType any) (*BufferedResponse, error) {

	bufResp, err := client.Exec(req)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = json.Unmarshal(bufResp.Body, respType)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	// Return nil if no error
	return bufResp, nil
}
