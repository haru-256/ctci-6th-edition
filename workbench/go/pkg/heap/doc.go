/*
Package heap provides a generic max heap implementation with comprehensive heap operations.

The max heap maintains the heap property where parent nodes are always greater than or equal
to their children, ensuring the maximum element is always at the root. The implementation
supports any ordered key type and allows associated values to be stored with each key.

# Features

- Generic implementation supporting any ordered key type (int, float64, string, etc.)
- Max heap property with efficient insertion and extraction
- Associated values can be stored with keys
- Heap sort functionality for in-place sorting
- BuildMaxHeap operation for converting arbitrary arrays
- Helper functions for heap index calculations
- Memory-efficient array-based storage

# Performance Characteristics

- Insert: O(log n)
- Pop (extract max): O(log n)
- Peek (view max): O(1)
- BuildMaxHeap: O(n)
- HeapSort: O(n log n)
- Space: O(n)

All operations maintain the heap property efficiently through up-heap and down-heap operations.

# Basic Usage

	// Create a new max heap for integers with string values
	heap := heap.NewMaxHeap[int, string]()

	// Insert key-value pairs
	heap.Insert(10, "ten")
	heap.Insert(30, "thirty")
	heap.Insert(20, "twenty")
	heap.Insert(40, "forty")

	// View the maximum element without removing it
	maxNode, err := heap.Max()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Maximum: %d -> %s\n", maxNode.Key, maxNode.Value)

	// Extract elements in descending order
	for heap.Size() > 0 {
		node, err := heap.Pop()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Popped: %d -> %s\n", node.Key, node.Value)
	}

# Advanced Usage with Custom Types

	// Define a custom ordered type
	type Priority int

	const (
		Low    Priority = 1
		Medium Priority = 5
		High   Priority = 10
	)

	// Create heap for custom priority system
	taskHeap := heap.NewMaxHeap[Priority, string]()

	// Insert tasks with priorities
	taskHeap.Insert(Medium, "Process emails")
	taskHeap.Insert(High, "Critical bug fix")
	taskHeap.Insert(Low, "Documentation update")
	taskHeap.Insert(High, "Security patch")

	// Process tasks in priority order (highest first)
	fmt.Println("Processing tasks:")
	for taskHeap.Size() > 0 {
		task, err := taskHeap.Pop()
		if err != nil {
			break
		}
		fmt.Printf("Priority %d: %s\n", task.Key, task.Value)
	}

# Building Heap from Existing Data

	// Start with an arbitrary collection
	heap := heap.NewMaxHeap[int, string]()
	
	// Manually populate with unordered data
	data := []struct{ key int; value string }{
		{15, "fifteen"}, {10, "ten"}, {20, "twenty"}, {8, "eight"}, {25, "twenty-five"},
	}
	
	// Add elements without maintaining heap property
	for _, item := range data {
		heap.Insert(item.key, item.value)
	}
	
	// Convert to proper max heap in O(n) time
	heap.BuildMaxHeap(heap)
	
	// Verify heap property
	max, err := heap.Max()
	if err == nil {
		fmt.Printf("Maximum after BuildMaxHeap: %d\n", max.Key) // Should be 25
	}

# Heap Sort

	// Create and populate a heap
	heap := heap.NewMaxHeap[int, string]()
	values := []int{64, 34, 25, 12, 22, 11, 90}
	
	for _, v := range values {
		heap.Insert(v, fmt.Sprintf("value_%d", v))
	}
	
	fmt.Printf("Original heap size: %d\n", heap.Size())
	
	// Sort the heap in-place (ascending order)
	heap.HeapSort(heap)
	
	// After sorting, elements are in ascending order in the underlying array
	// Note: The heap structure is destroyed after sorting
	fmt.Println("Heap sorted successfully")

# Heap Index Calculations

The package provides utility functions for heap index calculations:

	// For any element at index i:
	parentIndex := heap.Parent(i)    // Parent at (i-1)/2
	leftChild := heap.Left(i)        // Left child at 2*i+1  
	rightChild := heap.Right(i)      // Right child at 2*(i+1)
	
	// Example: For element at index 2
	fmt.Printf("Element at index 2:\n")
	fmt.Printf("  Parent index: %d\n", heap.Parent(2))    // 0
	fmt.Printf("  Left child index: %d\n", heap.Left(2))   // 5
	fmt.Printf("  Right child index: %d\n", heap.Right(2)) // 6

# Priority Queue Implementation

The heap can be used as a priority queue:

	type Task struct {
		ID          int
		Description string
	}

	// Higher numbers = higher priority
	pq := heap.NewMaxHeap[int, Task]()

	// Add tasks with priorities
	pq.Insert(5, Task{1, "Normal task"})
	pq.Insert(10, Task{2, "Important task"})
	pq.Insert(1, Task{3, "Low priority task"})
	pq.Insert(10, Task{4, "Another important task"})

	// Process tasks by priority
	for pq.Size() > 0 {
		taskNode, err := pq.Pop()
		if err != nil {
			break
		}
		task := taskNode.Value
		fmt.Printf("Processing (Priority %d): %s\n", 
			taskNode.Key, task.Description)
	}

# Memory Management

The heap implementation is designed for efficient memory usage:

	heap := heap.NewMaxHeap[int, string]()
	
	// Insert many elements
	for i := 0; i < 1000; i++ {
		heap.Insert(i, fmt.Sprintf("item_%d", i))
	}
	
	// Extract all elements (heap will shrink automatically)
	for heap.Size() > 0 {
		_, err := heap.Pop()
		if err != nil {
			break
		}
	}
	
	// Heap is now empty and ready for reuse
	fmt.Printf("Final heap size: %d\n", heap.Size()) // 0

Key memory features:
- Array automatically shrinks when elements are removed
- Popped elements are set to nil to prevent memory leaks
- No additional memory overhead beyond the stored elements
- Efficient reallocation strategies for dynamic sizing

# Error Handling

The heap operations can return errors in specific conditions:

	heap := heap.NewMaxHeap[int, string]()
	
	// Attempting to pop from empty heap
	_, err := heap.Pop()
	if errors.Is(err, heap.ErrorIsEmpty) {
		fmt.Println("Cannot pop from empty heap")
	}
	
	// Attempting to peek empty heap
	_, err = heap.Max()
	if errors.Is(err, heap.ErrorIsEmpty) {
		fmt.Println("Cannot peek empty heap")
	}
	
	// Always check errors for robustness
	node, err := heap.Pop()
	if err != nil {
		log.Printf("Pop failed: %v", err)
		return
	}
	fmt.Printf("Successfully popped: %d\n", node.Key)

# Thread Safety

The heap implementation is not thread-safe. For concurrent access, use external synchronization:

	heap := heap.NewMaxHeap[int, string]()
	var mu sync.Mutex
	
	// Safe concurrent insertions
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			mu.Lock()
			heap.Insert(value, fmt.Sprintf("item_%d", value))
			mu.Unlock()
		}(i)
	}
	wg.Wait()
	
	// Safe concurrent extractions
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			mu.Lock()
			node, err := heap.Pop()
			mu.Unlock()
			if err == nil {
				fmt.Printf("Extracted: %d\n", node.Key)
			}
		}()
	}
	wg.Wait()

# Implementation Details

The max heap uses a complete binary tree stored in an array:
- Root element is at index 0 (maximum element)
- For element at index i:
  - Parent is at index (i-1)/2
  - Left child is at index 2*i+1
  - Right child is at index 2*i+2
- Up-heap operation maintains heap property after insertion
- Down-heap operation maintains heap property after extraction
- BuildMaxHeap uses bottom-up heapification for O(n) construction
*/
package heap
