// Package stack provides a generic stack data structure.
package stack

import "errors"

var (
	// ErrorStackUnderflow is returned when trying to pop from an empty stack.
	ErrorStackOverflow  = errors.New("stack overflow")
	ErrorStackUnderflow = errors.New("stack underflow")
)

// Stack represents a generic stack.
type Stack struct {
	items []any
	size  int
	top   int
}

// NewStack creates and returns a new Stack.
func NewStack(size int) *Stack {
	return &Stack{
		items: make([]any, size),
		size:  size,
		top:   0,
	}
}

// IsEmpty checks if the stack is empty.
func (s *Stack) IsEmpty() bool {
	return s.top == 0
}

// Push adds an item to the top of the stack.
func (s *Stack) Push(item any) error {
	if s.top >= s.size {
		return ErrorStackOverflow
	}
	s.items[s.top] = item
	s.top++
	return nil
}

// Pop removes and returns the top item from the stack.
// It returns an error if the stack is empty.
func (s *Stack) Pop() (any, error) {
	if s.IsEmpty() {
		return nil, ErrorStackUnderflow
	}
	item := s.items[s.top-1]
	s.items[s.top-1] = nil // Clear the reference
	s.top--
	return item, nil
}

// Peek returns the top item from the stack without removing it.
// It returns an error if the stack is empty.
func (s *Stack) Peek() (any, error) {
	if s.IsEmpty() {
		return nil, ErrorStackUnderflow
	}
	return s.items[s.top-1], nil
}
