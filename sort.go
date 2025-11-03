package poindexter

import "sort"

// SortInts sorts a slice of integers in ascending order in place.
func SortInts(data []int) {
	sort.Ints(data)
}

// SortIntsDescending sorts a slice of integers in descending order in place.
func SortIntsDescending(data []int) {
	sort.Sort(sort.Reverse(sort.IntSlice(data)))
}

// SortStrings sorts a slice of strings in ascending order in place.
func SortStrings(data []string) {
	sort.Strings(data)
}

// SortStringsDescending sorts a slice of strings in descending order in place.
func SortStringsDescending(data []string) {
	sort.Sort(sort.Reverse(sort.StringSlice(data)))
}

// SortFloat64s sorts a slice of float64 values in ascending order in place.
func SortFloat64s(data []float64) {
	sort.Float64s(data)
}

// SortFloat64sDescending sorts a slice of float64 values in descending order in place.
func SortFloat64sDescending(data []float64) {
	sort.Sort(sort.Reverse(sort.Float64Slice(data)))
}

// SortBy sorts a slice using a custom less function.
// The less function should return true if data[i] should come before data[j].
func SortBy[T any](data []T, less func(i, j int) bool) {
	sort.Slice(data, less)
}

// SortByKey sorts a slice by extracting a comparable key from each element.
// The key function should return a value that implements constraints.Ordered.
func SortByKey[T any, K int | float64 | string](data []T, key func(T) K) {
	sort.Slice(data, func(i, j int) bool {
		return key(data[i]) < key(data[j])
	})
}

// SortByKeyDescending sorts a slice by extracting a comparable key from each element in descending order.
func SortByKeyDescending[T any, K int | float64 | string](data []T, key func(T) K) {
	sort.Slice(data, func(i, j int) bool {
		return key(data[i]) > key(data[j])
	})
}

// IsSorted checks if a slice of integers is sorted in ascending order.
func IsSorted(data []int) bool {
	return sort.IntsAreSorted(data)
}

// IsSortedStrings checks if a slice of strings is sorted in ascending order.
func IsSortedStrings(data []string) bool {
	return sort.StringsAreSorted(data)
}

// IsSortedFloat64s checks if a slice of float64 values is sorted in ascending order.
func IsSortedFloat64s(data []float64) bool {
	return sort.Float64sAreSorted(data)
}

// BinarySearch performs a binary search on a sorted slice of integers.
// Returns the index where target is found, or -1 if not found.
func BinarySearch(data []int, target int) int {
	idx := sort.SearchInts(data, target)
	if idx < len(data) && data[idx] == target {
		return idx
	}
	return -1
}

// BinarySearchStrings performs a binary search on a sorted slice of strings.
// Returns the index where target is found, or -1 if not found.
func BinarySearchStrings(data []string, target string) int {
	idx := sort.SearchStrings(data, target)
	if idx < len(data) && data[idx] == target {
		return idx
	}
	return -1
}
