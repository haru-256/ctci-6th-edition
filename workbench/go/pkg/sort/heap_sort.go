package sort

import (
	"cmp"

	"github.com/haru-256/ctci-6th-edition/pkg/heap"
)

// toPointerSlice converts a slice of values to a slice of pointers.
// This is a helper function used internally to interface with the heap package
// which works with pointer slices for memory efficiency.
func toPointerSlice[T any](arr []T) []*T {
	ptrs := make([]*T, len(arr))
	for i := range arr {
		ptrs[i] = &arr[i]
	}
	return ptrs
}

// fromPointerSlice converts a slice of pointers back to a slice of values.
// This is a helper function used internally to convert the result from the heap package
// back to the expected return type.
func fromPointerSlice[T any](ptrs []*T) []T {
	arr := make([]T, len(ptrs))
	for i := range ptrs {
		arr[i] = *ptrs[i]
	}
	return arr
}

// HeapSort sorts a slice using the heap sort algorithm.
//
// HeapSort provides guaranteed O(n log n) time complexity with O(1) extra space complexity.
// It is not a stable sort, meaning equal elements may not preserve their relative order.
//
// The algorithm works by:
// 1. Building a max heap from the input array
// 2. Repeatedly extracting the maximum element and placing it at the end
// 3. Rebuilding the heap property for the remaining elements
//
// Time Complexity: O(n log n) - guaranteed for all cases
// Space Complexity: O(n) - creates a new sorted slice, not in-place
// Stability: Not stable
//
// Parameters:
//   - arr: slice of any ordered type to be sorted
//
// Returns:
//   - A new slice containing the elements sorted in ascending order
//
// Example:
//
//	numbers := []int{64, 34, 25, 12, 22, 11, 90}
//	sorted := sort.HeapSort(numbers)
//	// sorted: [11, 12, 22, 25, 34, 64, 90]
//
//	words := []string{"banana", "apple", "cherry"}
//	sortedWords := sort.HeapSort(words)
//	// sortedWords: ["apple", "banana", "cherry"]
func HeapSort[T cmp.Ordered](arr []T) ([]T, error) {
	ptrs := toPointerSlice(arr)
	ptrs, err := heap.HeapSort(ptrs)
	if err != nil {
		return nil, err
	}
	return fromPointerSlice(ptrs), nil
}
