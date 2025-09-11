// Package queue provides a generic FIFO (First In, First Out) circular queue data structure.
//
// This package implements a fixed-capacity queue using a circular buffer approach for
// efficient memory usage and predictable performance. The queue is generic and
// can store any type T using Go's type parameters.
//
// Key Features:
//   - Generic implementation supporting any type T
//   - Fixed-capacity circular buffer for O(1) operations
//   - Thread-unsafe operations for maximum performance
//   - Comprehensive error handling with custom error types
//   - Zero-value safety with proper initialization checks
//
// The circular buffer implementation efficiently reuses array space as items are
// enqueued and dequeued, preventing the need to shift elements or reallocate memory.
// This makes it suitable for high-performance applications where predictable
// memory usage and consistent O(1) operation times are important.
//
// Performance Characteristics:
//   - Time Complexity: All operations (Enqueue, Dequeue, Peek, etc.) are O(1)
//   - Space Complexity: O(n) where n is the queue capacity
//   - Memory Usage: Fixed at creation time, no dynamic allocation during operations
//
// Thread Safety:
// This implementation is NOT thread-safe. If you need concurrent access,
// you must provide your own synchronization mechanisms.
//
// Example usage:
//
//	// Create a queue of integers with capacity 10
//	q := queue.NewQueue[int](10)
//
//	// Enqueue some values
//	if err := q.Enqueue(1); err != nil {
//	    log.Fatal(err) // Handle queue overflow
//	}
//	_ = q.Enqueue(2)
//	_ = q.Enqueue(3)
//
//	// Dequeue values (FIFO order)
//	val, err := q.Dequeue() // val = 1, err = nil
//	if err != nil {
//	    log.Fatal(err) // Handle queue underflow
//	}
//	val, _ = q.Dequeue()  // val = 2
//
//	// Peek at the front without removing
//	val, _ = q.Peek() // val = 3, item remains in queue
//
//	// Check queue state
//	fmt.Println("Empty:", q.IsEmpty()) // false
//	fmt.Println("Full:", q.IsFull())   // false
//	fmt.Println("Count:", q.Count())   // 1 (current items)
//	fmt.Println("Size:", q.Size())     // 10 (total capacity)
//
// Error Handling:
// The queue operations return specific errors for different failure conditions:
//   - ErrorQueueOverflow: Returned when trying to enqueue to a full queue
//   - ErrorQueueUnderflow: Returned when trying to dequeue from an empty queue
//
// These errors can be checked using errors.Is() for robust error handling:
//
//	if err := q.Enqueue(item); err != nil {
//	    if errors.Is(err, queue.ErrorQueueOverflow) {
//	        // Handle overflow specifically
//	    }
//	}
package queue
