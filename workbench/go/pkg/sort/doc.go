/*
Package sort provides efficient sorting algorithms with generic type support.

This package implements various sorting algorithms optimized for different use cases.
All sorting functions work with any type that implements the cmp.Ordered interface,
providing type safety and performance.

# Available Algorithms

- HeapSort: O(n log n) time complexity, O(1) extra space, not stable
- QuickSort: O(n log n) average case, O(n²) worst case, O(log n) extra space, not stable

# Performance Characteristics

HeapSort:
- Time: O(n log n) guaranteed
- Space: O(1) extra space
- Stability: Not stable
- Best for: Guaranteed performance, memory-constrained environments

QuickSort:
- Time: O(n log n) average, O(n²) worst case
- Space: O(log n) for recursion stack
- Stability: Not stable
- Best for: Average case performance, cache-friendly access patterns

# Basic Usage

	import "github.com/haru-256/ctci-6th-edition/pkg/sort"

	// Sort integers using heap sort
	numbers := []int{64, 34, 25, 12, 22, 11, 90}
	sorted := sort.HeapSort(numbers)
	// sorted: [11, 12, 22, 25, 34, 64, 90]

	// Sort strings using quick sort
	words := []string{"banana", "apple", "cherry", "date"}
	sorted = sort.QuickSort(words)
	// sorted: ["apple", "banana", "cherry", "date"]

	// Sort float64 values
	prices := []float64{19.99, 9.99, 29.99, 4.99}
	sorted = sort.HeapSort(prices)
	// sorted: [4.99, 9.99, 19.99, 29.99]

# Algorithm Selection Guide

Use HeapSort when:
- You need guaranteed O(n log n) performance
- Memory usage is critical (O(1) extra space)
- You're working with large datasets where worst-case performance matters

Use QuickSort when:
- Average case performance is more important than worst case
- You have good cache locality requirements
- The dataset is expected to be somewhat randomized

# Thread Safety

The sorting functions in this package are not thread-safe. If you need to sort
data concurrently, ensure proper synchronization or use separate goroutines
with separate data copies.

# Examples

	// Example: Sorting different types
	// Example: Sorting different types
	// (requires import "fmt")
	func ExampleUsage() {
		// Sort integers
		ints := []int{5, 2, 8, 1, 9}
		sortedInts := sort.HeapSort(ints)
		fmt.Println(sortedInts) // [1, 2, 5, 8, 9]

		// Sort strings
		strs := []string{"zebra", "apple", "banana"}
		sortedStrs := sort.QuickSort(strs)
		fmt.Println(sortedStrs) // [apple, banana, zebra]

		// Sort custom ordered types
		type Price float64
		prices := []Price{19.99, 9.99, 29.99}
		sortedPrices := sort.HeapSort(prices)
		fmt.Println(sortedPrices) // [9.99, 19.99, 29.99]
	}

	// Example: Performance comparison
	func ExamplePerformanceTest() {
		// Large dataset
		data := make([]int, 10000)
		for i := range data {
			data[i] = rand.Intn(10000)
		}

		// Copy for fair comparison
		heapData := make([]int, len(data))
		quickData := make([]int, len(data))
		copy(heapData, data)
		copy(quickData, data)

		// Time heap sort
		start := time.Now()
		sort.HeapSort(heapData)
		heapTime := time.Since(start)

		// Time quick sort
		start = time.Now()
		sort.QuickSort(quickData)
		quickTime := time.Since(start)

		fmt.Printf("HeapSort: %v, QuickSort: %v\n", heapTime, quickTime)
	}
*/
package sort
