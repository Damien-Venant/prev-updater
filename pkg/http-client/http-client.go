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

func (h *HttpClient) Get(path string, headers http.Header) (*http.Response, error) {
	url := fmt.Sprintf(formatUrl, h.BaseUrl, path)
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header = h.Headers
	setHeader(request, headers)

	fmt.Println(url)
	return h.client.Do(request)
}

func (h *HttpClient) Patch(path string, model any, headers http.Header) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", h.BaseUrl, path)
	request, err := http.NewRequest("PATCH", url, nil)
	if err != nil {
		return nil, err
	}
	request.Header = h.Headers
	setHeader(request, headers)
	result, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}
	request.Body.Read(result)
	return h.client.Do(request)
}

func setHeader(request *http.Request, headers http.Header) {
	for key, values := range headers {
		for _, val := range values {
			request.Header.Add(key, val)
		}
	}
}
