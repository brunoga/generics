package slices

func Filter[T1 any](slice []T1, f func(T1) bool) []T1 {
	result := make([]T1, 0, len(slice))
	for _, v := range slice {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}
