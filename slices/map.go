package slices

// Maps converts a []T1 into a []T2 by applying the given function to each
// elemant of []T1.
func Map[T1, T2 any](slice []T1, f func(T1) T2) []T2 {
	result := make([]T2, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}
