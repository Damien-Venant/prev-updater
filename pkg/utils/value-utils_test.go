package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoalesceString(t *testing.T) {
	type TestsStr struct {
		name          string
		inputValue    interface{}
		defaultValue  string
		expectedValue interface{}
	}

	tests := []TestsStr{
		{
			name:          "StringTest",
			inputValue:    "Test",
			defaultValue:  "",
			expectedValue: "Test",
		},
		{
			name:          "StringTestWithError",
			inputValue:    0.6,
			defaultValue:  "",
			expectedValue: "",
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("TestCoalesce_%s", test.name)
		t.Run(name, func(t *testing.T) {
			result := Coalesce(test.inputValue, test.defaultValue)
			assert.Equal(t, test.expectedValue, result)
		})
	}
}

func TestCoalesceInt(t *testing.T) {
	type TestsStr struct {
		name          string
		inputValue    interface{}
		defaultValue  int
		expectedValue interface{}
	}

	tests := []TestsStr{
		{
			name:          "IntTest",
			inputValue:    "Test",
			defaultValue:  0,
			expectedValue: 0,
		},
		{
			name:          "IntTestWithError",
			inputValue:    0.6,
			defaultValue:  0,
			expectedValue: 0,
		},
	}

	for _, test := range tests {
		name := fmt.Sprintf("TestCoalesce_%s", test.name)
		t.Run(name, func(t *testing.T) {
			result := Coalesce(test.inputValue, test.defaultValue)
			assert.Equal(t, test.expectedValue, result)
		})
	}
}
