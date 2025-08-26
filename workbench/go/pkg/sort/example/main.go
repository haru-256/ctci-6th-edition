package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/haru-256/ctci-6th-edition/pkg/sort"
)

func main() {
	fmt.Println("Sort Package Example")
	fmt.Println("===================")

	// Example 1: Sorting integers
	fmt.Println("\n1. Sorting integers:")
	numbers := []int{64, 34, 25, 12, 22, 11, 90}
	fmt.Printf("Original: %v\n", numbers)

	heapSorted := sort.HeapSort(numbers)
	fmt.Printf("HeapSort: %v\n", heapSorted)

	quickSorted := sort.QuickSort(numbers)
	fmt.Printf("QuickSort: %v\n", quickSorted)

	// Example 2: Sorting strings
	fmt.Println("\n2. Sorting strings:")
	words := []string{"banana", "apple", "cherry", "date", "elderberry"}
	fmt.Printf("Original: %v\n", words)

	sortedWords := sort.HeapSort(words)
	fmt.Printf("Sorted: %v\n", sortedWords)

	// Example 3: Sorting floats
	fmt.Println("\n3. Sorting float64 values:")
	prices := []float64{19.99, 9.99, 29.99, 4.99, 15.50}
	fmt.Printf("Original: %v\n", prices)

	sortedPrices := sort.QuickSort(prices)
	fmt.Printf("Sorted: %v\n", sortedPrices)

	// Example 4: Custom ordered type
	fmt.Println("\n4. Custom ordered type (Temperature):")
	type Temperature float64
	temps := []Temperature{98.6, 32.0, 212.0, 100.0, 0.0}
	fmt.Printf("Original: %v\n", temps)

	sortedTemps := sort.HeapSort(temps)
	fmt.Printf("Sorted: %v\n", sortedTemps)

	// Example 5: Performance comparison
	fmt.Println("\n5. Performance comparison with large dataset:")
	size := 10000
	data := make([]int, size)
	for i := 0; i < size; i++ {
		data[i] = rand.Intn(size)
	}

	// Test HeapSort
	heapData := make([]int, len(data))
	copy(heapData, data)
	start := time.Now()
	sort.HeapSort(heapData)
	heapTime := time.Since(start)

	// Test QuickSort
	quickData := make([]int, len(data))
	copy(quickData, data)
	start = time.Now()
	sort.QuickSort(quickData)
	quickTime := time.Since(start)

	fmt.Printf("Dataset size: %d elements\n", size)
	fmt.Printf("HeapSort time: %v\n", heapTime)
	fmt.Printf("QuickSort time: %v\n", quickTime)

	if quickTime < heapTime {
		fmt.Printf("QuickSort was %.2fx faster\n", float64(heapTime)/float64(quickTime))
	} else {
		fmt.Printf("HeapSort was %.2fx faster\n", float64(quickTime)/float64(heapTime))
	}

	// Example 6: Demonstrating that original slices are not modified
	fmt.Println("\n6. Original slice preservation:")
	original := []int{3, 1, 4, 1, 5}
	fmt.Printf("Before sorting: %v\n", original)

	sorted := sort.HeapSort(original)
	fmt.Printf("After HeapSort - Original: %v, Sorted: %v\n", original, sorted)
}
