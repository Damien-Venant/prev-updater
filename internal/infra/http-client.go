package infra

import (
	"fmt"
	"net/http"

	httpclient "github.com/Damien-Venant/prev-updater/pkg/http-client"
	"github.com/rs/zerolog"
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

func ConfigureHttpClient(config *HttpClientConfiguration, logger *zerolog.Logger) {
	bearerToken := fmt.Sprintf("Bearer %s", config.Token)
	headers := http.Header{}
	headers.Add("Authorization", bearerToken)
	httpClient = httpclient.New(config.BaseUrl, headers, logger)
}

func GetHttpClient() *httpclient.HttpClient {
	return httpClient
}
