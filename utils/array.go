package utils

func Filter[T any](data []T, test func(T) bool) (res []T) {
	for _, val := range data {
		if test(val) {
			res = append(res, val)
		}
	}
	return
}
