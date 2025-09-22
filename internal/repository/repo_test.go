package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigureRouteWithVersion(t *testing.T) {
	tests := []struct {
		name       string
		route      string
		parameters []any
		version    AzureDevOpsRepository
		wants      string
	}{
		{
			name:       "MultiParam",
			route:      "/api/%s/%d",
			parameters: []any{"test", 10},
			version:    AzureDevOpsRepository{version: "7.1"},
			wants:      "/api/test/10?api-version=7.1",
		},
		{
			name:       "MonoParam",
			route:      "/api/%d",
			parameters: []any{10},
			version:    AzureDevOpsRepository{version: "7.5"},
			wants:      "/api/10?api-version=7.5",
		},
		{
			name:       "MultiParamWithQueryParameter",
			route:      "/api/%s/%d?date=%s",
			parameters: []any{"test", 10, "test"},
			version:    AzureDevOpsRepository{version: "7.1"},
			wants:      "/api/test/10?api-version=7.1",
		},
		{
			name:       "MonoParamWithQueryParameter",
			route:      "/api/%d?date=%s",
			parameters: []any{10, "test"},
			version:    AzureDevOpsRepository{version: "7.5"},
			wants:      "/api/10?api-version=7.5",
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%s_%s", "TestConfigureRouteWithVersion", test.name)
		t.Run(testName, func(t *testing.T) {
			result := test.version.configureRouteWithVersion(test.route, test.parameters...)
			if result != test.wants {
				assert.NotEqual(t, test.wants, result)
			}
		})
	}
}

func TestErrorCodeMapping(t *testing.T) {
	tests := []struct {
		Name           string
		ErrorCode      int
		ExpectedResult error
	}{
		{
			Name:           "InternalServerError",
			ErrorCode:      http.StatusInternalServerError,
			ExpectedResult: InternalServerError,
		},
		{
			Name:           "BadRequestError",
			ErrorCode:      http.StatusBadRequest,
			ExpectedResult: BadRequestError,
		},
		{
			Name:           "NotFoundError",
			ErrorCode:      http.StatusNotFound,
			ExpectedResult: NotFoundError,
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("TestErrorCodeMapping%s", test.Name)
		t.Run(testName, func(t *testing.T) {
			err := errorCodeMapping(test.ErrorCode)
			assert.ErrorIs(t, err, test.ExpectedResult)
		})
	}
}

func TestErrorCodeMappingErrorUnknow(t *testing.T) {
	test := struct {
		ErrorCode      int
		ExpectedResult error
	}{
		ErrorCode:      http.StatusBadGateway,
		ExpectedResult: IdkError,
	}

	err := errorCodeMapping(test.ErrorCode)
	assert.ErrorIs(t, err, test.ExpectedResult)
}

func TestReadAndUnMarshall(t *testing.T) {
	type Person struct {
		FirstName string `json:"first-name"`
		LastName  string `json:"last-name"`
	}
	var person Person

	resultModel, _ := json.Marshal(Person{
		FirstName: "damien",
		LastName:  "venant",
	})

	reader := bytes.NewReader(resultModel)

	err := readAndUnmarshal[Person](reader, &person)

	assert.Nil(t, err)
	assert.Equal(t, "damien", person.FirstName)
	assert.Equal(t, "venant", person.LastName)
}
