package infra

import (
	"fmt"

	"github.com/prev-updater/pkg/http-client"
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
	httpClient = &httpclient.HttpClient{
		BaseUrl: config.BaseUrl,
	}
	httpClient.Headers.Add("Authorization", bearerToken)
}

func GetHttpClient() *httpclient.HttpClient {
	return httpClient
}
