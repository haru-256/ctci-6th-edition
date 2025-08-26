package sort

import (
	"cmp"
)

// QuickSort sorts a slice using the quicksort algorithm.
//
// QuickSort provides O(n log n) average time complexity with O(log n) space complexity
// for the recursion stack. In the worst case (already sorted data with poor pivot selection),
// it can degrade to O(n²) time complexity.
//
// The algorithm works by:
// 1. Choosing a pivot element (last element in this implementation)
// 2. Partitioning the array so elements ≤ pivot are on the left, > pivot on the right
// 3. Recursively sorting the left and right subarrays
//
// Time Complexity: O(n log n) average case, O(n²) worst case
// Space Complexity: O(log n) for recursion stack
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
//	sorted := sort.QuickSort(numbers)
//	// sorted: [11, 12, 22, 25, 34, 64, 90]
//
//	words := []string{"banana", "apple", "cherry"}
//	sortedWords := sort.QuickSort(words)
//	// sortedWords: ["apple", "banana", "cherry"]
func QuickSort[T cmp.Ordered](arr []T) []T {
	// Create a copy to avoid modifying the original slice
	result := make([]T, len(arr))
	copy(result, arr)

	quickSortInPlace(result, 0, len(result)-1)
	return result
}

// quickSortInPlace performs the actual quicksort algorithm in-place on a subarray.
// This is the internal recursive function that does the heavy lifting.
//
// Parameters:
//   - arr: the array to sort
//   - low: the starting index of the subarray to sort
//   - high: the ending index of the subarray to sort
func quickSortInPlace[T cmp.Ordered](arr []T, low, high int) {
	if low < high {
		// Partition the array and get the pivot index
		pivotIndex := partition(arr, low, high)

		// Recursively sort elements before and after partition
		quickSortInPlace(arr, low, pivotIndex-1)
		quickSortInPlace(arr, pivotIndex+1, high)
	}
}

// Partition rearranges the subarray arr[low..high] so that elements ≤ pivot
// are on the left and elements > pivot are on the right.
//
// This implementation uses the Lomuto partition scheme, which chooses the last
// element as the pivot. The function returns the final position of the pivot.
//
// Parameters:
//   - arr: the array to partition
//   - low: the starting index of the subarray
//   - high: the ending index of the subarray (pivot element)
//
// Returns:
//   - The final index position of the pivot element
//
// Example partitioning of [64, 34, 25, 12, 22, 11, 90] with pivot 90:
// Result: [64, 34, 25, 12, 22, 11, 90] with pivot at index 6
func Partition[T cmp.Ordered](arr []T, low, high int) int {
	return partition(arr, low, high)
}

// partition is the internal implementation of the partition function
func partition[T cmp.Ordered](arr []T, low, high int) int {
	// Choose the last element as pivot
	pivot := arr[high]

	// Index of smaller element, indicates the right position of pivot found so far
	i := low - 1

	for j := low; j < high; j++ {
		// If current element is smaller than or equal to pivot
		if arr[j] <= pivot {
			i++ // increment index of smaller element
			arr[i], arr[j] = arr[j], arr[i]
		}
	}

	// Place pivot in correct position
	arr[i+1], arr[high] = arr[high], arr[i+1]
	return i + 1
}
