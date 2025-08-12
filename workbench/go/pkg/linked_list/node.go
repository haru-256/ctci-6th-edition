// Package linked_list provides a generic doubly linked list implementation.
package linked_list

// Node represents a single node in a doubly linked list.
// It contains a value of generic type T and pointers to the previous and next nodes.
// The type T can be any type since nodes are used internally by LinkedList which requires comparable types.
type Node[T any] struct {
	Value T        // The data stored in this node
	Prev  *Node[T] // Pointer to the previous node in the list, nil if this is the head
	Next  *Node[T] // Pointer to the next node in the list, nil if this is the tail
}

// NewNode creates and returns a new Node with the specified value.
// The new node is initialized with nil Prev and Next pointers, making it suitable
// for insertion into a list or use as a standalone node.
func NewNode[T any](value T) *Node[T] {
	return &Node[T]{Value: value}
}
