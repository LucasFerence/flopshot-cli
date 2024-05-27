package api

import (
	"bytes"
	"net/http"
)

type HeaderPair struct {
	Key, Value string
}

// Tool to model our needed parameters for creating a new http.Request
type Request struct {
	Method string
	Body   []byte
	Url    string
	// need to add query parameters: https://stackoverflow.com/questions/30652577/go-doing-a-get-request-and-building-the-querystring
	// List of key/value pairs to assign as headers
	Headers []HeaderPair
}

func (req *Request) build(additionalHeaders ...HeaderPair) (*http.Request, error) {

	// Create http.Request from Request
	rawReq, err := http.NewRequest(req.Method, req.Url, bytes.NewReader(req.Body))

	if err != nil {
		return nil, err
	}

	for _, header := range req.Headers {
		rawReq.Header.Add(header.Key, header.Value)
	}

	if additionalHeaders != nil {
		for _, header := range additionalHeaders {
			rawReq.Header.Add(header.Key, header.Value)
		}
	}

	return rawReq, nil
}

func PostRequest(url string, body []byte, headers ...HeaderPair) Request {

	return Request{
		Method:  "POST",
		Url:     url,
		Body:    body,
		Headers: headers,
	}
}

func QueryItem(itemType string)
