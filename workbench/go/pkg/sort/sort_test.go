package sort

import (
	"math/rand"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// Test that both sorting algorithms produce the same results
func TestSortingAlgorithmsConsistency(t *testing.T) {
	testCases := []struct {
		name string
		data []int
	}{
		{"Empty", []int{}},
		{"Single", []int{42}},
		{"Two elements sorted", []int{1, 2}},
		{"Two elements reverse", []int{2, 1}},
		{"Already sorted", []int{1, 2, 3, 4, 5}},
		{"Reverse sorted", []int{5, 4, 3, 2, 1}},
		{"With duplicates", []int{3, 1, 4, 1, 5, 9, 2, 6, 5}},
		{"Random order", []int{64, 34, 25, 12, 22, 11, 90}},
		{"All same", []int{5, 5, 5, 5, 5}},
		{"Negative numbers", []int{-5, -1, -10, 0, 3, -3}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test both algorithms produce the same result
			heapResult := HeapSort(tc.data)
			quickResult := QuickSort(tc.data)

			assert.Equal(t, heapResult, quickResult,
				"HeapSort and QuickSort should produce the same result for %s", tc.name)

			// Verify they match Go's standard library
			if len(tc.data) > 0 {
				expected := make([]int, len(tc.data))
				copy(expected, tc.data)
				sort.Ints(expected)

				assert.Equal(t, expected, heapResult,
					"HeapSort should match standard library for %s", tc.name)
				assert.Equal(t, expected, quickResult,
					"QuickSort should match standard library for %s", tc.name)
			}
		})
	}
}

// Test with string data to ensure generics work correctly
func TestSortingAlgorithmsWithStrings(t *testing.T) {
	words := []string{"zebra", "apple", "banana", "cherry", "date", "elderberry"}

	heapResult := HeapSort(words)
	quickResult := QuickSort(words)

	expected := []string{"apple", "banana", "cherry", "date", "elderberry", "zebra"}

	assert.Equal(t, expected, heapResult, "HeapSort should work with strings")
	assert.Equal(t, expected, quickResult, "QuickSort should work with strings")
	assert.Equal(t, heapResult, quickResult, "Both algorithms should produce same result for strings")
}

// Test with float64 data
func TestSortingAlgorithmsWithFloats(t *testing.T) {
	floats := []float64{3.14159, 2.71828, 1.41421, 0.57721, 1.61803}

	heapResult := HeapSort(floats)
	quickResult := QuickSort(floats)

	expected := []float64{0.57721, 1.41421, 1.61803, 2.71828, 3.14159}

	assert.Equal(t, expected, heapResult, "HeapSort should work with floats")
	assert.Equal(t, expected, quickResult, "QuickSort should work with floats")
	assert.Equal(t, heapResult, quickResult, "Both algorithms should produce same result for floats")
}

// Test that original slices are not modified
func TestOriginalSliceNotModified(t *testing.T) {
	original := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
	originalCopy := make([]int, len(original))
	copy(originalCopy, original)

	heapResult := HeapSort(original)
	assert.Equal(t, originalCopy, original, "HeapSort should not modify original slice")

	quickResult := QuickSort(original)
	assert.Equal(t, originalCopy, original, "QuickSort should not modify original slice")

	// Results should be sorted
	assert.True(t, sort.IntsAreSorted(heapResult), "HeapSort result should be sorted")
	assert.True(t, sort.IntsAreSorted(quickResult), "QuickSort result should be sorted")
}

// Performance comparison test (not a benchmark, just a functional test)
func TestPerformanceComparison(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance comparison in short mode")
	}

	// Create test data
	size := 10000
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Intn(size)
	}

	// Test HeapSort
	heapData := make([]int, len(data))
	copy(heapData, data)
	start := time.Now()
	heapResult := HeapSort(heapData)
	heapTime := time.Since(start)

	// Test QuickSort
	quickData := make([]int, len(data))
	copy(quickData, data)
	start = time.Now()
	quickResult := QuickSort(quickData)
	quickTime := time.Since(start)

	// Both should produce the same sorted result
	assert.Equal(t, heapResult, quickResult, "Both algorithms should produce the same result")
	assert.True(t, sort.IntsAreSorted(heapResult), "Result should be sorted")

	t.Logf("Performance for %d elements:", size)
	t.Logf("HeapSort: %v", heapTime)
	t.Logf("QuickSort: %v", quickTime)

	// Both should complete in reasonable time (less than 1 second for 10k elements)
	assert.Less(t, heapTime, time.Second, "HeapSort should complete in reasonable time")
	assert.Less(t, quickTime, time.Second, "QuickSort should complete in reasonable time")
}

// Test worst-case scenario for QuickSort (already sorted data)
func TestQuickSortWorstCase(t *testing.T) {
	// Create already sorted data (worst case for basic QuickSort)
	size := 1000
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = i
	}

	result := QuickSort(data)

	// Should still work correctly despite worst-case input
	assert.True(t, sort.IntsAreSorted(result), "QuickSort should handle worst-case input")
	assert.Equal(t, data, result, "Result should be the same as input (already sorted)")
}

// Test heap sort with worst-case scenario
func TestHeapSortWorstCase(t *testing.T) {
	// HeapSort has consistent O(n log n) performance regardless of input
	size := 1000
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = i
	}

	result := HeapSort(data)

	assert.True(t, sort.IntsAreSorted(result), "HeapSort should handle any input efficiently")
	assert.Equal(t, data, result, "Result should be the same as input (already sorted)")
}

// Edge case: very large numbers
func TestSortingWithLargeNumbers(t *testing.T) {
	largeNumbers := []int{
		1000000000, 999999999, 1000000001, 500000000, 1500000000,
		-1000000000, -999999999, 0, 1, -1,
	}

	heapResult := HeapSort(largeNumbers)
	quickResult := QuickSort(largeNumbers)

	expected := []int{
		-1000000000, -999999999, -1, 0, 1,
		500000000, 999999999, 1000000000, 1000000001, 1500000000,
	}

	assert.Equal(t, expected, heapResult, "HeapSort should handle large numbers")
	assert.Equal(t, expected, quickResult, "QuickSort should handle large numbers")
}

// Test with custom ordered type
func TestCustomOrderedType(t *testing.T) {
	type Temperature float64

	temps := []Temperature{98.6, 32.0, 212.0, 100.0, 0.0, -40.0}

	heapResult := HeapSort(temps)
	quickResult := QuickSort(temps)

	expected := []Temperature{-40.0, 0.0, 32.0, 98.6, 100.0, 212.0}

	assert.Equal(t, expected, heapResult, "HeapSort should work with custom ordered types")
	assert.Equal(t, expected, quickResult, "QuickSort should work with custom ordered types")
}

// Stress test with multiple random datasets
func TestStressTestMultipleDatasets(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	// Test 100 different random datasets
	for i := 0; i < 100; i++ {
		size := rand.Intn(500) + 1 // 1 to 500 elements
		data := make([]int, size)
		for j := 0; j < size; j++ {
			data[j] = rand.Intn(1000) - 500 // -500 to 499
		}

		heapResult := HeapSort(data)
		quickResult := QuickSort(data)

		// Both should produce the same sorted result
		assert.Equal(t, heapResult, quickResult,
			"Both algorithms should produce same result for dataset %d", i)
		assert.True(t, sort.IntsAreSorted(heapResult),
			"HeapSort result should be sorted for dataset %d", i)
		assert.True(t, sort.IntsAreSorted(quickResult),
			"QuickSort result should be sorted for dataset %d", i)
	}
}
