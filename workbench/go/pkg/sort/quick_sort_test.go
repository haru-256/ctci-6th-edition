package sort

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuickSort_EmptySlice(t *testing.T) {
	var empty []int
	result := QuickSort(empty)
	assert.Empty(t, result, "QuickSort should handle empty slices")
}

func TestQuickSort_SingleElement(t *testing.T) {
	single := []int{42}
	result := QuickSort(single)
	assert.Equal(t, []int{42}, result, "QuickSort should handle single element")
}

func TestQuickSort_AlreadySorted(t *testing.T) {
	sorted := []int{1, 2, 3, 4, 5}
	result := QuickSort(sorted)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result, "QuickSort should handle already sorted arrays")
}

func TestQuickSort_ReverseSorted(t *testing.T) {
	reverse := []int{5, 4, 3, 2, 1}
	result := QuickSort(reverse)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result, "QuickSort should handle reverse sorted arrays")
}

func TestQuickSort_WithDuplicates(t *testing.T) {
	duplicates := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
	result := QuickSort(duplicates)
	expected := []int{1, 1, 2, 3, 4, 5, 5, 6, 9}
	assert.Equal(t, expected, result, "QuickSort should handle duplicates correctly")
}

func TestQuickSort_RandomOrder(t *testing.T) {
	random := []int{64, 34, 25, 12, 22, 11, 90}
	result := QuickSort(random)
	expected := []int{11, 12, 22, 25, 34, 64, 90}
	assert.Equal(t, expected, result, "QuickSort should sort random arrays correctly")
}

func TestQuickSort_Strings(t *testing.T) {
	words := []string{"banana", "apple", "cherry", "date"}
	result := QuickSort(words)
	expected := []string{"apple", "banana", "cherry", "date"}
	assert.Equal(t, expected, result, "QuickSort should work with strings")
}

func TestQuickSort_Float64(t *testing.T) {
	floats := []float64{3.14, 2.71, 1.41, 0.57}
	result := QuickSort(floats)
	expected := []float64{0.57, 1.41, 2.71, 3.14}
	assert.Equal(t, expected, result, "QuickSort should work with float64")
}

func TestQuickSort_NegativeNumbers(t *testing.T) {
	negatives := []int{-5, -1, -10, 0, 3, -3}
	result := QuickSort(negatives)
	expected := []int{-10, -5, -3, -1, 0, 3}
	assert.Equal(t, expected, result, "QuickSort should handle negative numbers")
}

func TestQuickSort_LargeDataset(t *testing.T) {
	// Create a large dataset
	size := 1000
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Intn(1000)
	}

	result := QuickSort(data)

	// Verify it's sorted
	assert.True(t, sort.IntsAreSorted(result), "Large dataset should be sorted")
	assert.Len(t, result, size, "Result should have same length as input")
}

func TestQuickSort_DoesNotModifyOriginal(t *testing.T) {
	original := []int{3, 1, 4, 1, 5}
	originalCopy := make([]int, len(original))
	copy(originalCopy, original)

	result := QuickSort(original)

	// Original should be unchanged
	assert.Equal(t, originalCopy, original, "QuickSort should not modify original slice")
	// Result should be different (sorted)
	assert.NotEqual(t, original, result, "Result should be different from original")
}

func TestQuickSort_CustomOrderedType(t *testing.T) {
	type Price float64

	prices := []Price{19.99, 9.99, 29.99, 4.99}
	result := QuickSort(prices)
	expected := []Price{4.99, 9.99, 19.99, 29.99}
	assert.Equal(t, expected, result, "QuickSort should work with custom ordered types")
}

func TestQuickSort_CorrectnessAgainstStandardLibrary(t *testing.T) {
	// Test multiple random datasets against Go's standard library sort
	for i := 0; i < 100; i++ {
		size := rand.Intn(100) + 1 // 1 to 100 elements
		data := make([]int, size)
		for j := 0; j < size; j++ {
			data[j] = rand.Intn(1000)
		}

		// Get expected result using standard library
		expected := make([]int, len(data))
		copy(expected, data)
		sort.Ints(expected)

		// Get our result
		result := QuickSort(data)

		assert.Equal(t, expected, result, "QuickSort result should match standard library sort for dataset %d", i)
	}
}

func TestPartition_BasicFunctionality(t *testing.T) {
	arr := []int{64, 34, 25, 12, 22, 11, 90}
	pivotIndex := Partition(arr, 0, len(arr)-1)

	// After partitioning, all elements to the left of pivot should be <= pivot
	// and all elements to the right should be > pivot
	pivot := arr[pivotIndex]

	for i := 0; i < pivotIndex; i++ {
		assert.LessOrEqual(t, arr[i], pivot, "Elements to the left of pivot should be <= pivot")
	}

	for i := pivotIndex + 1; i < len(arr); i++ {
		assert.Greater(t, arr[i], pivot, "Elements to the right of pivot should be > pivot")
	}
}

func TestPartition_SingleElement(t *testing.T) {
	arr := []int{42}
	pivotIndex := Partition(arr, 0, 0)
	assert.Equal(t, 0, pivotIndex, "Partition of single element should return index 0")
	assert.Equal(t, 42, arr[0], "Single element should remain unchanged")
}

func TestPartition_TwoElements(t *testing.T) {
	// Test with already sorted pair
	arr1 := []int{1, 2}
	pivotIndex1 := Partition(arr1, 0, 1)
	assert.Equal(t, 1, pivotIndex1, "Pivot should be at correct position")
	assert.Equal(t, []int{1, 2}, arr1, "Already sorted pair should remain sorted")

	// Test with reverse sorted pair
	arr2 := []int{2, 1}
	pivotIndex2 := Partition(arr2, 0, 1)
	assert.Equal(t, 0, pivotIndex2, "Pivot should be at correct position")
	assert.Equal(t, []int{1, 2}, arr2, "Reverse sorted pair should be corrected")
}

func TestPartition_AllEqual(t *testing.T) {
	arr := []int{5, 5, 5, 5, 5}
	pivotIndex := Partition(arr, 0, len(arr)-1)

	// With all equal elements, the partition should still work correctly
	assert.True(t, pivotIndex >= 0 && pivotIndex < len(arr), "Pivot index should be valid")

	// All elements should still be equal
	for i := 0; i < len(arr); i++ {
		assert.Equal(t, 5, arr[i], "All elements should remain equal")
	}
}

// Benchmark tests
func BenchmarkQuickSort_Random100(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = rand.Intn(1000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuickSort(data)
	}
}

func BenchmarkQuickSort_Random1000(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = rand.Intn(10000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuickSort(data)
	}
}

func BenchmarkQuickSort_Random10000(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = rand.Intn(100000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuickSort(data)
	}
}

func BenchmarkQuickSort_Sorted1000(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuickSort(data)
	}
}

func BenchmarkQuickSort_Reverse1000(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = len(data) - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		QuickSort(data)
	}
}

func BenchmarkPartition_Random1000(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = rand.Intn(10000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		testData := make([]int, len(data))
		copy(testData, data)
		b.StartTimer()

		Partition(testData, 0, len(testData)-1)
	}
}
