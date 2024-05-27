package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/99designs/keyring"
)

// Constants for keyring identifiers
const keyringService = "flopshot"
const loginTokenKey = "loginToken"

func openRing() (keyring.Keyring, error) {
	return keyring.Open(keyring.Config{
		ServiceName: keyringService,
	})
}

type flopshotClient struct {

	// Underlying client for executing requests
	HttpClient *http.Client

	// calculated token, might be nil
	Token string
}

func NewFlopshotClient() flopshotClient {

	client := flopshotClient{
		HttpClient: &http.Client{},
	}

	return client
}

func (client *flopshotClient) InitializeAuth(token string) {

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

func (client *flopshotClient) IsAuthenticated() (bool, string) {

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

func (client *flopshotClient) ExecuteRaw(req Request) (*BufferedResponse, error) {

	var rawReq *http.Request
	var err error

	isAuth, token := client.IsAuthenticated()

	if isAuth {
		// If authenticated, add headers for auth
		rawReq, err = req.build(
			HeaderPair{
				Key: "authorization",
				Value: "Bearer" + token,
			},
			HeaderPair {
				Key: "content-type",
				Value: "application/json",
			},
		)
	} else {
		// Otherwise, just build it with no additional headers
		rawReq, err = req.build()
	}

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resp, err := client.HttpClient.Do(rawReq)

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

func (client *flopshotClient) Execute(req Request, respType any) (*BufferedResponse, error) {

	bufResp, err := client.ExecuteRaw(req)

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
