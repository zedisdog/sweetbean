package tools

import "math"

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

// GroupSlice group slice to a new slice with each group has same groupCount
func GroupSlice[T any](s []T, groupCount int) (result [][]T) {
	result = make([][]T, int(math.Ceil(float64(len(s))/float64(groupCount))))
	index := 0

	for i, d := range s {
		if result[index] == nil {
			result[index] = make([]T, 0, groupCount)
		}
		result[index] = append(result[index], d)
		if (i+1)%groupCount == 0 {
			index++
		}
	}

	return
}
