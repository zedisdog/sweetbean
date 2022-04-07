package tools

func InSlice[T comparable](x T, s []T) (result int) {
	result = -1
	for index, i := range s {
		if i == x {
			result = index
			break
		}
	}

	return
}
