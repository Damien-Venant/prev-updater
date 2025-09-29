package httpclient

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/rs/zerolog"
)

type (
	HttpClient struct {
		BaseUrl string
		Headers http.Header
		client  *http.Client
		logger  *zerolog.Logger
	}
	HttpClientInterface interface {
		Get(path string, headers http.Header) (*http.Response, error)
		Patch(path string, body []byte, headers http.Header) (*http.Response, error)
	}
)

const (
	formatUrl string = "%s/%s"
)

func New(baseUrl string, headers http.Header, logger *zerolog.Logger) *HttpClient {
	return &HttpClient{
		BaseUrl: baseUrl,
		Headers: headers,
		client:  &http.Client{},
		logger:  logger,
	}
}

func (h *HttpClient) Get(path string, headers http.Header) (*http.Response, error) {
	url := fmt.Sprintf(formatUrl, h.BaseUrl, path)
	h.logger.
		Info().
		Dict("request-data", zerolog.Dict().Str("url", url).Str("method", "GET")).
		Msg("Send request")
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		h.logger.
			Error().
			Stack().
			Err(err).
			Send()
		return nil, err
	}
	request.Header = h.Headers
	setHeader(request, headers)

	return h.client.Do(request)
}

func (h *HttpClient) Patch(path string, body []byte, headers http.Header) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", h.BaseUrl, path)
	h.logger.
		Info().
		Dict("request-data", zerolog.Dict().Str("url", url).Str("method", "PATCH").Str("body", string(body))).
		Msg("Send request")
	request, err := http.NewRequest("PATCH", url, bytes.NewReader(body))
	if err != nil {
		h.logger.
			Error().
			Stack().
			Err(err).
			Send()
		return nil, err
	}
	request.Header = h.Headers
	request.Header.Set("Content-Type", "application/json-patch+json")
	setHeader(request, headers)

	return h.client.Do(request)
}

func setHeader(request *http.Request, headers http.Header) {
	for key, values := range headers {
		for _, val := range values {
			request.Header.Add(key, val)
		}
	}
}
