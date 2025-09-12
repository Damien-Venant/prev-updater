package httpclient

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetHeader(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add("Authorization", "Bearer dsqdsdksnqlkdnlqs")
	request, _ := http.NewRequest("GET", "", nil)

	setHeader(request, headers)

	for key, value := range headers {
		headerValue, ok := request.Header[key]
		assert.True(t, ok)
		assert.Equal(t, value, headerValue)
	}
}
