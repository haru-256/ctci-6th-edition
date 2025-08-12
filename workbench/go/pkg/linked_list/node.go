package linked_list

type Node[T any] struct {
	Value T
	Prev  *Node[T]
	Next  *Node[T]
}

func NewNode[T any](value T) *Node[T] {
	return &Node[T]{Value: value}
}
