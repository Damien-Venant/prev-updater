package queryslice

type Predicate[T any] func(predicate T) bool

func Filter[T any](source []T, predicate Predicate[T]) []T {
	result := make([]T, 0, len(source))
	for _, val := range source {
		if predicate(val) {
			result = append(result, val)
		}
	}
	return result
}
