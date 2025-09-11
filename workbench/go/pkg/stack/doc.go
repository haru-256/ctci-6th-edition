// Package stack provides a generic LIFO (Last In, First Out) stack data structure.
//
// This package implements a fixed-capacity stack using an array-based approach for
// efficient memory usage and predictable performance. The stack is generic and
// can store any type T.
//
// Example usage:
//
//	// Create a stack of integers with capacity 10
//	s := stack.NewStack[int](10)
//
//	// Push some values
//	_ = s.Push(1)
//	_ = s.Push(2)
//	_ = s.Push(3)
//
//	// Pop values (LIFO order)
//	val, _ := s.Pop() // val = 3
//	val, _ = s.Pop()  // val = 2
//
//	// Check stack state
//	fmt.Println("Empty:", s.IsEmpty()) // false
//	fmt.Println("Count:", s.Count())   // 1
//	fmt.Println("Size:", s.Size())     // 10
package stack
