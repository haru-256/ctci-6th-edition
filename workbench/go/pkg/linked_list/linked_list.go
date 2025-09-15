package linked_list

import (
	"errors"
	"sync"
)

// ErrorNodeIsNil is returned when an operation is attempted on a nil Node.
var (
	ErrorNodeIsNil = errors.New("node is nil")
	// ErrorNodeNotFound is returned when a delete operation cannot find the target node.
	ErrorNodeNotFound = errors.New("node not found")
)

// LinkedList represents a generic doubly linked list.
// The linked list is thread-safe and can be safely used by multiple goroutines concurrently.
// Read operations use RWMutex.RLock() to allow concurrent reads,
// while write operations use RWMutex.Lock() for exclusive access.
// No external synchronization is required for concurrent use.
type LinkedList[T comparable] struct {
	head  *Node[T]     // Points to the first node in the list
	tail  *Node[T]     // Points to the last node in the list
	mutex sync.RWMutex // Protects list operations for thread safety
}

// NewLinkedList creates and returns a new empty doubly linked list.
// The returned list is ready for use and properly initialized.
func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{}
}

// Head returns the first node in the list or nil if the list is empty.
// This method is thread-safe and uses read locking.
func (l *LinkedList[T]) Head() *Node[T] {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.head
}

// Tail returns the last node in the list or nil if the list is empty.
// This method is thread-safe and uses read locking.
func (l *LinkedList[T]) Tail() *Node[T] {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	return l.tail
}

// Search traverses the list from head to tail looking for a node with the specified value.
// It returns a pointer to the first node found with the matching value, or nil if not found.
// The search performs a linear traversal with O(n) time complexity.
// This method is thread-safe and uses RLock to allow concurrent read access.
func (l *LinkedList[T]) Search(value T) *Node[T] {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	current := l.head
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
// This operation has O(1) time complexity and is thread-safe using exclusive locking.
func (l *LinkedList[T]) Prepend(value T) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	newNode := NewNode(value)
	if l.head == nil { // if list is empty
		l.head = newNode
		l.tail = newNode
	} else {
		newNode.Next = l.head
		l.head.Prev = newNode
		l.head = newNode
	}
}

// Insert adds a new node with the specified value immediately after the given node.
// The 'after' parameter must not be nil, or ErrorNodeIsNil will be returned.
// If 'after' is the current tail, the new node becomes the new tail.
// This operation maintains all doubly-linked relationships and has O(1) time complexity.
// This method is thread-safe using exclusive locking.
func (l *LinkedList[T]) Insert(value T, after *Node[T]) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if after == nil {
		return ErrorNodeIsNil
	}
	newNode := NewNode(value)
	if after.Next != nil {
		after.Next.Prev = newNode
		newNode.Next = after.Next
	} else {
		l.tail = newNode
	}
	after.Next = newNode
	newNode.Prev = after

	return nil
}

// search is an internal helper method that searches for a node without locking.
// This method assumes the caller already holds the appropriate lock.
func (l *LinkedList[T]) search(value T) *Node[T] {
	current := l.head
	for current != nil {
		if current.Value == value {
			return current
		}
		current = current.Next
	}
	return nil
}

// Delete removes the first node with the specified value from the list.
// It returns ErrorNodeNotFound if no node with the given value exists.
// When a node is deleted, all links are properly updated to maintain list integrity.
// If the deleted node was the head or tail, those pointers are updated accordingly.
// This operation has O(n) time complexity due to the search phase.
// This method is thread-safe using exclusive locking.
func (l *LinkedList[T]) Delete(value T) error {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	node := l.search(value)
	if node == nil {
		return ErrorNodeNotFound
	}
	if node.Prev != nil {
		node.Prev.Next = node.Next
	} else {
		l.head = node.Next
	}
	if node.Next != nil {
		node.Next.Prev = node.Prev
	} else {
		l.tail = node.Prev
	}

	// Help GC by breaking references from the deleted node.
	node.Prev = nil
	node.Next = nil

	return nil
}
