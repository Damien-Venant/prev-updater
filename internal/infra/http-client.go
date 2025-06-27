package infra

import (
	"fmt"
	"net/http"

	httpclient "github.com/prev-updater/pkg/http-client"
)

type (
	HttpClientConfiguration struct {
		BaseUrl string
		Token   string
	}
)

var (
	httpClient *httpclient.HttpClient
)

func ConfigureHttpClient(config *HttpClientConfiguration) {
	bearerToken := fmt.Sprintf("Bearer %s", config.Token)
	headers := http.Header{}
	headers.Add("Authorization", bearerToken)
	httpClient = httpclient.New(config.BaseUrl, headers)
}

func GetHttpClient() *httpclient.HttpClient {
	return httpClient
}
