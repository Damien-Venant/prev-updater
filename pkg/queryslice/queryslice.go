package queryslice

import (
	"slices"
)

type Predicate[T any] func(predicate T) bool
type Trans[T any, K any] func(trans T, index int) K

func Filter[T any](source []T, predicateFunc Predicate[T]) []T {
	result := make([]T, 0, len(source))
	for _, val := range source {
		if predicateFunc(val) {
			result = append(result, val)
		}
	}
	return result
}

// Transform
func Transform[T any, K any](source []T, transFunc Trans[T, K]) []K {
	result := make([]K, 0, len(source))
	for index, val := range source {
		result = append(result, transFunc(val, index))
	}
	return result
}

// TransformParallel is used to make transform in parallel with goroutines
// If you transform some data with an call like API Call, Query or another I/O operation we recommend to used this function else use Transform
func TransformParallel[T any, K any](source []T, transFunc Trans[T, K]) []K {
	type valIndex struct {
		index int
		value K
	}
	srcLength := len(source)
	resultChan := make(chan valIndex, srcLength)
	for index, val := range source {
		go func() {
			res := transFunc(val, index)
			resultChan <- valIndex{index: index, value: res}
		}()
	}
	resultValIndex := make([]valIndex, 0, srcLength)
	for i := 0; i < srcLength; i++ {
		resultValIndex = append(resultValIndex, <-resultChan)
	}

	close(resultChan)
	slices.SortFunc(resultValIndex, func(a, b valIndex) int {
		return a.index - b.index
	})

	result := Transform(resultValIndex, func(val valIndex, _ int) K {
		return val.value
	})

	return result
}

func FindIndex[T any](source []T, filterFunc Predicate[T]) int {
	for index, val := range source {
		if filterFunc(val) {
			return index
		}
	}
	return -1
}
