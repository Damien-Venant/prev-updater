package httpclient

import (
	"bytes"
	"fmt"
	"net/http"
)

type (
	HttpClient struct {
		BaseUrl string
		Headers http.Header
		client  *http.Client
	}
)

const (
	formatUrl string = "%s/%s"
)

func New(baseUrl string, headers http.Header) *HttpClient {
	return &HttpClient{
		BaseUrl: baseUrl,
		Headers: headers,
		client:  &http.Client{},
	}
}

func (h *HttpClient) Get(path string, headers http.Header) (*http.Response, error) {
	url := fmt.Sprintf(formatUrl, h.BaseUrl, path)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header = h.Headers
	setHeader(request, headers)

	return h.client.Do(request)
}

func (h *HttpClient) Patch(path string, body []byte, headers http.Header) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", h.BaseUrl, path)
	request, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header = h.Headers
	request.Header.Set("Content-Type", "application/json-patch+json")
	setHeader(request, headers)
	if err != nil {
		return nil, err
	}

	return h.client.Do(request)
}

// TODO : write a test for this function
func setHeader(request *http.Request, headers http.Header) {
	for key, values := range headers {
		for _, val := range values {
			request.Header.Add(key, val)
		}
	}
}
