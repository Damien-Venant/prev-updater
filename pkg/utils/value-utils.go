package utils

func Coalesce[K any](val interface{}, defaultValue K) K {
	result, ok := val.(K)
	if !ok {
		return defaultValue
	}
	return result
}
