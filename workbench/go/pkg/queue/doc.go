// Package queue provides a generic FIFO (First In, First Out) circular queue data structure.
//
// This package implements a fixed-capacity queue using a circular buffer approach for
// efficient memory usage and predictable performance. The queue is generic and
// can store any type T.
//
// Example usage:
//
//	// Create a queue of integers with capacity 10
//	q := queue.NewQueue[int](10)
//
//	// Enqueue some values
//	_ = q.Enqueue(1)
//	_ = q.Enqueue(2)
//	_ = q.Enqueue(3)
//
//	// Dequeue values (FIFO order)
//	val, _ := q.Dequeue() // val = 1
//	val, _ = q.Dequeue()  // val = 2
//
//	// Peek at the front without removing
//	val, _ = q.Peek() // val = 3
//
//	// Check queue state
//	fmt.Println("Empty:", q.IsEmpty()) // false
//	fmt.Println("Full:", q.IsFull())   // false
//	fmt.Println("Count:", q.Count())   // 1
//	fmt.Println("Size:", q.Size())     // 10
package queue
