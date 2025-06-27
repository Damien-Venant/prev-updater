package repository

import (
	"fmt"
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
			name:       "multi-param",
			route:      "/api/%s/%d",
			parameters: []any{"test", 10},
			version:    AzureDevOpsRepository{version: "7.1"},
			wants:      "/api/test/10?api-version=7.1",
		},
		{
			name:       "mono-param",
			route:      "/api/%d",
			parameters: []any{10},
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
