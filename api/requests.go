package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type ListResponse[T any] struct {
	Items []T `json:"items"`
}

type QueryParams struct {
	K, V string
}

func (client *FlopshotClient) QueryData(dataType string, data any, query []QueryParams) error {

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%s/%s", ClientUrl, dataType),
		nil,
	)

	if err != nil {
		return err
	}

	if query != nil {
		q := req.URL.Query()
		for _, v := range query {
			q.Add(v.K, v.V)
		}

		req.URL.RawQuery = q.Encode()
	}

	_, err = client.ExecR(req, &data)
	return err
}

func (client *FlopshotClient) WriteData(dataType string, data any) error {

	val, err := json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("%s/%s/write", ClientUrl, dataType),
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
