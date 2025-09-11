package stack

import "errors"

var (
	// ErrorStackOverflow is returned when trying to push to a full stack.
	ErrorStackOverflow = errors.New("stack overflow")
	// ErrorStackUnderflow is returned when trying to pop from an empty stack.
	ErrorStackUnderflow = errors.New("stack underflow")
)

// Stack represents a generic LIFO (Last In, First Out) stack using a fixed-size array.
// It stores values of type T directly and provides thread-unsafe operations.
// The zero value is not ready to use; use NewStack to create a new stack.
//
// The stack has a fixed capacity determined at creation time and will return
// ErrorStackOverflow when attempting to push beyond capacity, or ErrorStackUnderflow
// when attempting to pop from an empty stack.
//
// Time complexity:
//   - Push: O(1)
//   - Pop: O(1)
//   - Peek: O(1)
//   - IsEmpty/IsFull/Size/Count: O(1)
//
// Space complexity: O(n) where n is the capacity.
type Stack[T any] struct {
	items []T // slice to store stack items
	size  int // maximum number of items the stack can hold
	count int // current number of items in the stack
}

// NewStack creates and returns a new Stack with the specified capacity.
// The size parameter must be greater than 0, otherwise the function will panic.
// The stack is initially empty and ready to use.
//
// The returned stack will have a fixed capacity and will not grow beyond the
// specified size. Use Push to add items and Pop to remove them.
//
// Parameters:
//   - size: The maximum number of items the stack can hold (must be > 0)
//
// Returns:
//   - A new Stack instance ready for use
//
// Panics:
//   - If size <= 0
//
// Example:
//
//	stack := NewStack[int](10) // Creates a stack that can hold 10 integers
//	stack := NewStack[string](5) // Creates a stack that can hold 5 strings
func NewStack[T any](size int) *Stack[T] {
	if size <= 0 {
		panic("stack size must be greater than 0")
	}
	return &Stack[T]{
		items: make([]T, size),
		size:  size,
		count: 0,
	}
}

// IsEmpty checks if the stack is empty.
// Returns true if there are no elements in the stack.
func (s *Stack[T]) IsEmpty() bool {
	return s.count == 0
}

// IsFull checks if the stack is full.
// Returns true if the stack has reached its maximum capacity.
func (s *Stack[T]) IsFull() bool {
	return s.count == s.size
}

// Push adds an item to the top of the stack.
// Returns ErrorStackOverflow if the stack is full.
//
// Example:
//
//	err := stack.Push(42)
//	if err != nil {
//	    // handle stack overflow
//	}
func (s *Stack[T]) Push(item T) error {
	if s.IsFull() {
		return ErrorStackOverflow
	}
	s.items[s.count] = item
	s.count++
	return nil
}

// Pop removes and returns the top item from the stack.
// Returns ErrorStackUnderflow if the stack is empty.
// The stack follows LIFO (Last In, First Out) order.
//
// Example:
//
//	item, err := stack.Pop()
//	if err != nil {
//	    // handle stack underflow
//	} else {
//	    fmt.Println(item) // use the popped item
//	}
func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if s.IsEmpty() {
		return zero, ErrorStackUnderflow
	}
	item := s.items[s.count-1]
	s.items[s.count-1] = zero // Clear the reference to prevent memory leaks
	s.count--
	return item, nil
}

// Peek returns the top item from the stack without removing it.
// Returns ErrorStackUnderflow if the stack is empty.
// This operation does not modify the stack.
//
// Example:
//
//	item, err := stack.Peek()
//	if err != nil {
//	    // handle empty stack
//	} else {
//	    fmt.Println(item) // use the top item without removing it
//	}
func (s *Stack[T]) Peek() (T, error) {
	if s.IsEmpty() {
		var zero T
		return zero, ErrorStackUnderflow
	}
	return s.items[s.count-1], nil
}

// Size returns the maximum capacity of the stack.
// This is the size that was specified when the stack was created.
func (s *Stack[T]) Size() int {
	return s.size
}

// Count returns the current number of items in the stack.
// This value ranges from 0 (empty) to Size() (full).
func (s *Stack[T]) Count() int {
	return s.count
}
