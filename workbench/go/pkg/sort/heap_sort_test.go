package sort

import (
	"math/rand"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeapSort_EmptySlice(t *testing.T) {
	var empty []int
	result, err := HeapSort(empty)
	require.NoError(t, err)
	assert.Empty(t, result, "HeapSort should handle empty slices")
}

func TestHeapSort_SingleElement(t *testing.T) {
	single := []int{42}
	result, err := HeapSort(single)
	require.NoError(t, err)
	assert.Equal(t, []int{42}, result, "HeapSort should handle single element")
}

func TestHeapSort_AlreadySorted(t *testing.T) {
	sorted := []int{1, 2, 3, 4, 5}
	result, err := HeapSort(sorted)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result, "HeapSort should handle already sorted arrays")
}

func TestHeapSort_ReverseSorted(t *testing.T) {
	reverse := []int{5, 4, 3, 2, 1}
	result, err := HeapSort(reverse)
	require.NoError(t, err)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, result, "HeapSort should handle reverse sorted arrays")
}

func TestHeapSort_WithDuplicates(t *testing.T) {
	duplicates := []int{3, 1, 4, 1, 5, 9, 2, 6, 5}
	result, err := HeapSort(duplicates)
	require.NoError(t, err)
	expected := []int{1, 1, 2, 3, 4, 5, 5, 6, 9}
	assert.Equal(t, expected, result, "HeapSort should handle duplicates correctly")
}

func TestHeapSort_RandomOrder(t *testing.T) {
	random := []int{64, 34, 25, 12, 22, 11, 90}
	result, err := HeapSort(random)
	require.NoError(t, err)
	expected := []int{11, 12, 22, 25, 34, 64, 90}
	assert.Equal(t, expected, result, "HeapSort should sort random arrays correctly")
}

func TestHeapSort_Strings(t *testing.T) {
	words := []string{"banana", "apple", "cherry", "date"}
	result, err := HeapSort(words)
	require.NoError(t, err)
	expected := []string{"apple", "banana", "cherry", "date"}
	assert.Equal(t, expected, result, "HeapSort should work with strings")
}

func TestHeapSort_Float64(t *testing.T) {
	floats := []float64{3.14, 2.71, 1.41, 0.57}
	result, err := HeapSort(floats)
	require.NoError(t, err)
	expected := []float64{0.57, 1.41, 2.71, 3.14}
	assert.Equal(t, expected, result, "HeapSort should work with float64")
}

func TestHeapSort_NegativeNumbers(t *testing.T) {
	negatives := []int{-5, -1, -10, 0, 3, -3}
	result, err := HeapSort(negatives)
	require.NoError(t, err)
	expected := []int{-10, -5, -3, -1, 0, 3}
	assert.Equal(t, expected, result, "HeapSort should handle negative numbers")
}

func TestHeapSort_LargeDataset(t *testing.T) {
	// Create a large dataset
	size := 1000
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Intn(1000)
	}

	result, err := HeapSort(data)
	require.NoError(t, err)

	// Verify it's sorted
	assert.True(t, sort.IntsAreSorted(result), "Large dataset should be sorted")
	assert.Len(t, result, size, "Result should have same length as input")
}

func TestHeapSort_DoesNotModifyOriginal(t *testing.T) {
	original := []int{3, 1, 4, 1, 5}
	originalCopy := make([]int, len(original))
	copy(originalCopy, original)

	result, err := HeapSort(original)
	require.NoError(t, err)

	// Original should be unchanged
	assert.Equal(t, originalCopy, original, "HeapSort should not modify original slice")
	// Result should be different (sorted)
	assert.NotEqual(t, original, result, "Result should be different from original")
}

func TestHeapSort_CustomOrderedType(t *testing.T) {
	type Price float64

	prices := []Price{19.99, 9.99, 29.99, 4.99}
	result, err := HeapSort(prices)
	require.NoError(t, err)
	expected := []Price{4.99, 9.99, 19.99, 29.99}
	assert.Equal(t, expected, result, "HeapSort should work with custom ordered types")
}

func TestHeapSort_CorrectnessAgainstStandardLibrary(t *testing.T) {
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
		result, err := HeapSort(data)
		require.NoError(t, err)

		assert.Equal(t, expected, result, "HeapSort result should match standard library sort for dataset %d", i)
	}
}

// Benchmark tests
func BenchmarkHeapSort_Random100(b *testing.B) {
	data := make([]int, 100)
	for i := range data {
		data[i] = rand.Intn(1000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HeapSort(data)
	}
}

func BenchmarkHeapSort_Random1000(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = rand.Intn(10000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HeapSort(data)
	}
}

func BenchmarkHeapSort_Random10000(b *testing.B) {
	data := make([]int, 10000)
	for i := range data {
		data[i] = rand.Intn(100000)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HeapSort(data)
	}
}

func BenchmarkHeapSort_Sorted1000(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HeapSort(data)
	}
}

func BenchmarkHeapSort_Reverse1000(b *testing.B) {
	data := make([]int, 1000)
	for i := range data {
		data[i] = len(data) - i
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = HeapSort(data)
	}
}
