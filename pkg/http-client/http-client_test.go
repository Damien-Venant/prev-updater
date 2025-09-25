package httpclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockRoundTripper struct {
	req  *http.Request
	resp *http.Response
	err  error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.req = req
	return m.resp, m.err
}

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

func TestSetHeaderWithTabHeaderValue(t *testing.T) {
	headers := http.Header{}
	headers.Add("Content-Type", "application/json")
	headers.Add("Content-Type", "application/bson")
	headers.Add("Authorization", "Bearer dsqdsdksnqlkdnlqs")
	headers.Add("Authorization", "Bearer tests")
	request, _ := http.NewRequest("GET", "", nil)

	setHeader(request, headers)

	for key, value := range headers {
		headerValue, ok := request.Header[key]
		assert.True(t, ok)
		assert.Equal(t, value, headerValue)
	}
}

// mockRoundTripper intercepte les requêtes HTTP pour les tester sans réseau.

// helper pour créer un client avec le mock
func newTestHttpClient(t *testing.T, roundTripper http.RoundTripper) *HttpClient {
	logger := zerolog.Nop()
	client := &http.Client{Transport: roundTripper}
	return &HttpClient{
		BaseUrl: "http://localhost",
		Headers: http.Header{
			"X-Global": []string{"global"},
		},
		client: client,
		logger: &logger,
	}
}

func TestHttpClient_Get(t *testing.T) {
	expectedBody := `{"status":"ok"}`
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(expectedBody)),
		Header:     make(http.Header),
	}

	mock := &mockRoundTripper{resp: mockResp}

	client := newTestHttpClient(t, mock)

	customHeaders := http.Header{}
	customHeaders.Set("Authorization", "Bearer token123")

	resp, err := client.Get("test-path", customHeaders)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, expectedBody, string(body))

	assert.NotNil(t, mock.req)
	assert.Equal(t, "GET", mock.req.Method)
	assert.Equal(t, "http://localhost/test-path", mock.req.URL.String())
	assert.Equal(t, "global", mock.req.Header.Get("X-Global"))
	assert.Equal(t, "Bearer token123", mock.req.Header.Get("Authorization"))
}

func TestHttpClient_Patch(t *testing.T) {
	expectedBody := `{"result":"patched"}`
	mockResp := &http.Response{
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBufferString(expectedBody)),
		Header:     make(http.Header),
	}

	mock := &mockRoundTripper{resp: mockResp}

	client := newTestHttpClient(t, mock)

	body := []byte(`[{"op":"replace","path":"/name","value":"test"}]`)
	customHeaders := http.Header{}
	customHeaders.Set("Authorization", "Bearer patchtoken")

	resp, err := client.Patch("update-path", body, customHeaders)
	require.NoError(t, err)
	require.NotNil(t, resp)

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, expectedBody, string(respBody))

	assert.NotNil(t, mock.req)
	assert.Equal(t, "PATCH", mock.req.Method)
	assert.Equal(t, "http://localhost/update-path", mock.req.URL.String())
	assert.Equal(t, "application/json-patch+json", mock.req.Header.Get("Content-Type"))
	assert.Equal(t, "Bearer patchtoken", mock.req.Header.Get("Authorization"))

	reqBody, _ := ioutil.ReadAll(mock.req.Body)
	assert.Equal(t, string(body), string(reqBody))
}
