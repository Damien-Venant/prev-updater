package queryslice

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

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

func TestQuerySliceTransform(t *testing.T) {
	tests := []struct {
		name          string
		transformFunc func(val any, _ int) any
		values        []any
		results       []any
	}{
		{
			name: "Multiplication",
			transformFunc: func(val any, _ int) any {
				intVal, _ := val.(int)
				return intVal * 2
			},
			values:  []any{1, 2, 3, 4, 5},
			results: []any{2, 4, 6, 8, 10},
		},
		{
			name: "Concatenation",
			transformFunc: func(val any, _ int) any {
				strVal, _ := val.(string)
				return fmt.Sprintf("test_%s", strVal)
			},
			values:  []any{"1", "2", "3", "4", "5"},
			results: []any{"test_1", "test_2", "test_3", "test_4", "test_5"},
		},
		{
			name: "Conversion",
			transformFunc: func(val any, _ int) any {
				intVal, _ := val.(int)
				return fmt.Sprintf("%d", intVal)
			},
			values:  []any{1, 2, 3, 4, 5},
			results: []any{"1", "2", "3", "4", "5"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestQuerySliceTransform_%s", test.name), func(te *testing.T) {
			result := Transform(test.values, test.transformFunc)
			assert.Equal(te, test.results, result)
		})
	}
}

func TestQuerySliceTransformParallel(t *testing.T) {
	tests := []struct {
		name          string
		transformFunc func(val any, _ int) any
		values        []any
		results       []any
	}{
		{
			name: "Multiplication",
			transformFunc: func(val any, _ int) any {
				intVal, _ := val.(int)
				r := rand.Intn(400)
				time.Sleep(time.Millisecond * time.Duration(r))
				return intVal * 2
			},
			values:  []any{1, 2, 3, 4, 5},
			results: []any{2, 4, 6, 8, 10},
		},
		{
			name: "Concatenation",
			transformFunc: func(val any, _ int) any {
				strVal, _ := val.(string)
				r := rand.Intn(400)
				time.Sleep(time.Millisecond * time.Duration(r))
				return fmt.Sprintf("test_%s", strVal)
			},
			values:  []any{"1", "2", "3", "4", "5"},
			results: []any{"test_1", "test_2", "test_3", "test_4", "test_5"},
		},
		{
			name: "Conversion",
			transformFunc: func(val any, _ int) any {
				intVal, _ := val.(int)
				r := rand.Intn(400)
				time.Sleep(time.Millisecond * time.Duration(r))
				return fmt.Sprintf("%d", intVal)
			},
			values:  []any{1, 2, 3, 4, 5},
			results: []any{"1", "2", "3", "4", "5"},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("TestQuerySliceTransformParallel_%s", test.name), func(te *testing.T) {
			result := Transform(test.values, test.transformFunc)
			assert.Equal(te, test.results, result)
		})
	}
}
