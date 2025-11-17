package stypelib

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
