package collections

import "slices"

func RemoveDuplicates[T comparable](elements []T) []T {
	duplicates := make(map[T]struct{}, len(elements))
	for i := range elements {
		if _, ok := duplicates[elements[i]]; ok {
			elements = slices.Delete(elements, i, i+1)
		}
		duplicates[elements[i]] = struct{}{}
	}
	return elements
}
