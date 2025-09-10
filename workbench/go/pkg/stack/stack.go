// Package stack provides a generic LIFO (Last In, First Out) stack data structure.
package stack

import "errors"

var (
	// ErrorStackOverflow is returned when trying to push to a full stack.
	ErrorStackOverflow = errors.New("stack overflow")
	// ErrorStackUnderflow is returned when trying to pop from an empty stack.
	ErrorStackUnderflow = errors.New("stack underflow")
)

// Stack represents a generic LIFO (Last In, First Out) stack using a fixed-size array.
// It stores pointers to values of type T and provides thread-unsafe operations.
// The zero value is not ready to use; use NewStack to create a new stack.
type Stack[T any] struct {
	items []*T // slice to store pointers to stack items
	size  int  // maximum number of items the stack can hold
	top   int  // index of the next available position (also represents current count)
}

// NewStack creates and returns a new Stack with the specified capacity.
// The size parameter must be greater than 0, otherwise the function will panic.
// The stack is initially empty and ready to use.
//
// Example:
//
//	stack := NewStack[int](10) // Creates a stack that can hold 10 integers
func NewStack[T any](size int) *Stack[T] {
	if size <= 0 {
		panic("stack size must be greater than 0")
	}
	return &Stack[T]{
		items: make([]*T, size),
		size:  size,
		top:   0,
	}
}

// IsEmpty checks if the stack is empty.
// Returns true if there are no elements in the stack.
func (s *Stack[T]) IsEmpty() bool {
	return s.top == 0
}

// IsFull checks if the stack is full.
// Returns true if the stack has reached its maximum capacity.
func (s *Stack[T]) IsFull() bool {
	return s.top >= s.size
}

// Push adds a pointer to an item to the top of the stack.
// Returns ErrorStackOverflow if the stack is full.
//
// Example:
//
//	item := 42
//	err := stack.Push(&item)
//	if err != nil {
//	    // handle stack overflow
//	}
func (s *Stack[T]) Push(item *T) error {
	if s.IsFull() {
		return ErrorStackOverflow
	}
	s.items[s.top] = item
	s.top++
	return nil
}

// Pop removes and returns a pointer to the top item from the stack.
// Returns ErrorStackUnderflow if the stack is empty.
// The stack follows LIFO (Last In, First Out) order.
//
// Example:
//
//	item, err := stack.Pop()
//	if err != nil {
//	    // handle stack underflow
//	} else {
//	    fmt.Println(*item) // use the popped item
//	}
func (s *Stack[T]) Pop() (*T, error) {
	if s.IsEmpty() {
		return nil, ErrorStackUnderflow
	}
	item := s.items[s.top-1]
	s.items[s.top-1] = nil // Clear the reference to prevent memory leaks
	s.top--
	return item, nil
}

// Peek returns a pointer to the top item from the stack without removing it.
// Returns ErrorStackUnderflow if the stack is empty.
// This operation does not modify the stack.
//
// Example:
//
//	item, err := stack.Peek()
//	if err != nil {
//	    // handle empty stack
//	} else {
//	    fmt.Println(*item) // use the top item without removing it
//	}
func (s *Stack[T]) Peek() (*T, error) {
	if s.IsEmpty() {
		return nil, ErrorStackUnderflow
	}
	return s.items[s.top-1], nil
}
