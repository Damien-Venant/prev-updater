package queryslice

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
	srcLength := len(source)
	resultChan := make(chan K, srcLength)
	for index, val := range source {
		go func() {
			resultChan <- transFunc(val, index)
		}()
	}
	result := make([]K, 0, srcLength)
	for i := 0; i < srcLength; i++ {
		result = append(result, <-resultChan)
	}

	close(resultChan)
	return result
}
