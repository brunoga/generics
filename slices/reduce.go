package slices

func Reduce[T1, T2 any](slice []T1, f func(T1, T2) T2, initial T2) T2 {
	result := initial
	for _, v := range slice {
		result = f(v, result)
	}
	return result
}
