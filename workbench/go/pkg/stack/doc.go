// Package stack provides a generic LIFO (Last In, First Out) stack data structure.
//
// This package implements a fixed-capacity stack using an array-based approach for
// efficient memory usage and predictable performance. The stack is generic and
// can store any type T using Go's type parameters.
//
// Key Features:
//   - Generic implementation supporting any type T
//   - Fixed-capacity array-based storage for O(1) operations
//   - Thread-unsafe operations for maximum performance
//   - Comprehensive error handling with custom error types
//   - Zero-value safety with proper initialization checks
//   - Value-based API (no pointers required)
//
// The array-based implementation provides predictable memory usage and consistent
// performance characteristics, making it suitable for applications where memory
// efficiency and predictable operation times are important.
//
// Performance Characteristics:
//   - Time Complexity: All operations (Push, Pop, Peek, etc.) are O(1)
//   - Space Complexity: O(n) where n is the stack capacity
//   - Memory Usage: Fixed at creation time, no dynamic allocation during operations
//
// Thread Safety:
// This implementation is NOT thread-safe. If you need concurrent access,
// you must provide your own synchronization mechanisms.
//
// Example usage:
//
//	// Create a stack of integers with capacity 10
//	s := stack.NewStack[int](10)
//
//	// Push some values
//	if err := s.Push(1); err != nil {
//	    log.Fatal(err) // Handle stack overflow
//	}
//	_ = s.Push(2)
//	_ = s.Push(3)
//
//	// Pop values (LIFO order)
//	val, err := s.Pop() // val = 3, err = nil
//	if err != nil {
//	    log.Fatal(err) // Handle stack underflow
//	}
//	val, _ = s.Pop()  // val = 2
//
//	// Peek at the top without removing
//	val, _ = s.Peek() // val = 1, item remains on stack
//
//	// Check stack state
//	fmt.Println("Empty:", s.IsEmpty()) // false
//	fmt.Println("Full:", s.IsFull())   // false
//	fmt.Println("Count:", s.Count())   // 1 (current items)
//	fmt.Println("Size:", s.Size())     // 10 (total capacity)
//
// Error Handling:
// The stack operations return specific errors for different failure conditions:
//   - ErrorStackOverflow: Returned when trying to push to a full stack
//   - ErrorStackUnderflow: Returned when trying to pop from an empty stack
//
// These errors can be checked using errors.Is() for robust error handling:
//
//	if err := s.Push(item); err != nil {
//	    if errors.Is(err, stack.ErrorStackOverflow) {
//	        // Handle overflow specifically
//	    }
//	}
package stack
