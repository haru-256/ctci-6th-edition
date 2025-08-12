// Package linked_list provides a generic doubly linked list implementation.
package linked_list

import "errors"

// ErrorNodeIsNil is returned when an operation is attempted on a nil Node.
var (
	ErrorNodeIsNil = errors.New("node is nil")
	// ErrorNodeNotFound is returned when a delete operation cannot find the target node.
	ErrorNodeNotFound = errors.New("node not found")
)

// LinkedList represents a generic doubly linked list.
// It maintains pointers to both the head and tail nodes for efficient operations.
// The list supports any comparable type T for its values.
type LinkedList[T comparable] struct {
	Head *Node[T]
	Tail *Node[T]
}

// NewLinkedList creates and returns a new empty LinkedList.
// The list is initialized with nil head and tail pointers.
func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// Search traverses the list from head to tail looking for a node with the specified value.
// It returns a pointer to the first node found with the matching value, or nil if not found.
// The search performs a linear traversal with O(n) time complexity.
func (l *LinkedList[T]) Search(value T) *Node[T] {
	current := l.Head
	for current != nil {
		if current.Value == value {
			return current
		}
		current = current.Next
	}
	return nil
}

// Prepend adds a new node with the specified value to the beginning of the list.
// If the list is empty, the new node becomes both head and tail.
// Otherwise, the new node is inserted before the current head and becomes the new head.
// This operation has O(1) time complexity.
func (l *LinkedList[T]) Prepend(value T) {
	newNode := NewNode(value)
	if l.Head == nil { // if list is empty
		l.Head = newNode
		l.Tail = newNode
	} else {
		newNode.Next = l.Head
		l.Head.Prev = newNode
		l.Head = newNode
	}
}

// Insert adds a new node with the specified value immediately after the given node.
// The 'after' parameter must not be nil, or ErrorNodeIsNil will be returned.
// If 'after' is the current tail, the new node becomes the new tail.
// This operation maintains all doubly-linked relationships and has O(1) time complexity.
func (l *LinkedList[T]) Insert(value T, after *Node[T]) error {
	if after == nil {
		return ErrorNodeIsNil
	}
	newNode := NewNode(value)
	if after.Next != nil {
		after.Next.Prev = newNode
		newNode.Next = after.Next
	} else {
		l.Tail = newNode
	}
	after.Next = newNode
	newNode.Prev = after

	return nil
}

// Delete removes the first node with the specified value from the list.
// It returns ErrorNodeNotFound if no node with the given value exists.
// When a node is deleted, all links are properly updated to maintain list integrity.
// If the deleted node was the head or tail, those pointers are updated accordingly.
// This operation has O(n) time complexity due to the search phase.
func (l *LinkedList[T]) Delete(value T) error {
	node := l.Search(value)
	if node == nil {
		return ErrorNodeNotFound
	}
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		l.Head = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		l.Tail = node.Prev
	}
	return nil
}
