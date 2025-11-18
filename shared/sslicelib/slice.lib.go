package sslicelib

func FilterNotInSlice[T comparable](original []T, reference []T) []T {
	referenceMap := make(map[T]struct{})
	for _, item := range reference {
		referenceMap[item] = struct{}{}
	}

	var result []T
	for _, item := range original {
		if _, found := referenceMap[item]; !found {
			result = append(result, item)
		}
	}
	return result
}

func Filter[T any](slice []T, fn func(index int, item T) bool) []T {
	var result []T
	for idx, item := range slice {
		if fn(idx, item) {
			result = append(result, item)
		}
	}
	return result
}

func Map[T any, R any](slice []T, fn func(index int, item T) R) []R {
	var result []R
	for idx, item := range slice {
		result = append(result, fn(idx, item))
	}
	return result
}

func Find[T any](slice []T, fn func(index int, item T) bool) (result T, ok bool) {
	for idx, item := range slice {
		if fn(idx, item) {
			return item, true
		}
	}
	return result, false
}
