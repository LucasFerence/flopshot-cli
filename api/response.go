package api

import (
	"net/http"
)

// Wrapper around http.Response, with data into bytes
type BufferedResponse struct {
	RawResponse *http.Response

	Body []byte
}
