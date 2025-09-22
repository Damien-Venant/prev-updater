package queryslice

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuerySliceFilter(t *testing.T) {
	tests := []struct {
		name      string
		predicate func(predicate any) bool
		values    []any
		results   []any
	}{
		{
			name: "FilterByModulo2",
			predicate: func(pre any) bool {
				res, _ := pre.(int)
				return res%2 == 0
			},
			values:  []any{1, 2, 3, 4, 5},
			results: []any{2, 4},
		},
		{
			name: "FilterByNameEquality",
			predicate: func(pre any) bool {
				res, _ := pre.(string)
				return res == "test"
			},
			values:  []any{"super", "foo", "bar", "test"},
			results: []any{"test"},
		},
	}

	for _, test := range tests {
		testName := fmt.Sprintf("%s_%s", "TestQuerySliceFilter", test.name)
		t.Run(testName, func(t *testing.T) {
			result := Filter(test.values, test.predicate)
			assert.Equal(t, test.results, result)
		})
	}
}
