package httpclient

import (
	"encoding/json"
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

func (h *HttpClient) Get(path string, header http.Header) (*http.Response, error) {
	url := fmt.Sprintf(formatUrl, h.BaseUrl, path)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return h.client.Do(request)
}

func (h *HttpClient) Patch(path string, model any, header http.Header) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", h.BaseUrl, path)
	request, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header = h.Headers
	result, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	request.Body.Read(result)
	return h.client.Do(request)
}
