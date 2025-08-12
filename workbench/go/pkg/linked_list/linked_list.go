package linked_list

import "errors"

var (
	ErrorNodeIsNil    = errors.New("node is nil")
	ErrorNodeNotFound = errors.New("node not found")
)

type LinkedList[T comparable] struct {
	Head *Node[T]
	Tail *Node[T]
}

func NewLinkedList[T comparable]() *LinkedList[T] {
	return &LinkedList[T]{}
}

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
