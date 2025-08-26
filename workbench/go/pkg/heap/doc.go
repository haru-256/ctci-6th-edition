/*
Package heap provides a generic heap implementation with flexible comparison functions.

The heap maintains the heap property based on a user-provided comparison function, allowing
for both max heap and min heap behavior. The implementation supports any type and provides
efficient insertion, extraction, and heap operations.

# Features

- Generic implementation supporting any type with custom comparison functions
- Configurable as max heap or min heap through comparison functions
- Efficient heap operations with logarithmic time complexity
- Heap sort functionality for in-place sorting
- BuildMaxHeap and BuildMinHeap operations for converting arbitrary arrays
- Helper functions for heap index calculations
- Memory-efficient array-based storage
- Convenience functions for common ordered types

# Performance Characteristics

- Insert: O(log n)
- Pop (extract top): O(log n)
- Peek (view top): O(1)
- BuildHeap: O(n)
- HeapSort: O(n log n)
- Space: O(n)

All operations maintain the heap property efficiently through up-heap and down-heap operations.

# Basic Usage

	// Create a new max heap for integers
	maxHeap := heap.NewMaxHeap[int]()

	// Insert elements
	maxHeap.Insert(10)
	maxHeap.Insert(30)
	maxHeap.Insert(20)
	maxHeap.Insert(40)

	// View the maximum element without removing it
	max, err := maxHeap.Peek()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Maximum: %d\n", *max)

	// Extract elements in descending order (max heap)
	for maxHeap.Size() > 0 {
		item, err := maxHeap.Pop()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Popped: %d\n", *item)
	}

# Min Heap Usage

	// Create a new min heap for integers
	minHeap := heap.NewMinHeap[int]()

	// Insert elements
	minHeap.Insert(30)
	minHeap.Insert(10)
	minHeap.Insert(20)
	minHeap.Insert(5)

	// Extract elements in ascending order (min heap)
	for minHeap.Size() > 0 {
		item, err := minHeap.Pop()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Popped: %d\n", *item)
	}
	// Output: 5, 10, 20, 30
	# Custom Comparison Functions

	// Define a custom comparison function for max heap behavior
	maxCmp := func(a, b *int) int {
		if *a > *b {
			return 1
		} else if *a < *b {
			return -1
		}
		return 0
	}

	// Create heap with custom comparison
	customHeap := heap.NewHeap(maxCmp)
	customHeap.Insert(15)
	customHeap.Insert(25)
	customHeap.Insert(5)

	// Custom types with comparison
	type Person struct {
		Name string
		Age  int
	}

	// Compare by age (older person has higher priority)
	personCmp := func(a, b *Person) int {
		if a.Age > b.Age {
			return 1
		} else if a.Age < b.Age {
			return -1
		}
		return 0
	}

	personHeap := heap.NewHeap(personCmp)
	personHeap.Insert(Person{Name: "Alice", Age: 30})
	personHeap.Insert(Person{Name: "Bob", Age: 25})
	personHeap.Insert(Person{Name: "Charlie", Age: 35})

	// Charlie (age 35) will be at the top

# Building Heap from Existing Data

	// Create heap from existing array
	values := []int{15, 10, 20, 8, 25}
	valuePtrs := make([]*int, len(values))
	for i, v := range values {
		valuePtrs[i] = &v
	}

	// Build max heap in O(n) time
	maxHeap := heap.BuildMaxHeap(valuePtrs)

	// Build min heap in O(n) time
	minHeap := heap.BuildMinHeap(valuePtrs)

	// Verify heap property
	max, err := maxHeap.Peek()
	if err == nil {
		fmt.Printf("Maximum after BuildMaxHeap: %d\n", *max) // Should be 25
	}

	min, err := minHeap.Peek()
	if err == nil {
		fmt.Printf("Minimum after BuildMinHeap: %d\n", *min) // Should be 8
	}# Heap Sort

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
