package pkg

import (
	"cmp"
	"errors"
)

// ErrorNodeIsNil is returned when an operation is attempted on a nil TreeNode.
var ErrorNodeIsNil = errors.New("node is nil")

// TreeNode represents a node in a binary search tree.
// It holds a generic data type T that must satisfy the cmp.Ordered interface,
// allowing for comparison between elements.
type TreeNode[T cmp.Ordered] struct {
	data   T
	parent *TreeNode[T]
	left   *TreeNode[T]
	right  *TreeNode[T]
	size   int
}

// NewTreeNode creates and returns a new TreeNode with the given data.
// The new node's size is initialized to 1.
func NewTreeNode[T cmp.Ordered](d T) *TreeNode[T] {
	return &TreeNode[T]{
		data: d,
		size: 1,
	}
}

// validateNode checks if a TreeNode is nil and returns an error if it is.
func validateNode[T cmp.Ordered](node *TreeNode[T]) error {
	if node == nil {
		return ErrorNodeIsNil
	}
	return nil
}

// InsertInOrder inserts a new data element into the binary search tree
// rooted at the current node. It maintains the binary search tree property:
// elements less than or equal to the current node's data go to the left,
// and elements greater than the current node's data go to the right.
// The size of each node in the path of insertion is incremented.
func (node *TreeNode[T]) InsertInOrder(data T) error {
	if err := validateNode(node); err != nil {
		return err
	}

	var err error
	if data <= node.data {
		if node.left == nil {
			err = node.setLeftChild(NewTreeNode(data))
		} else {
			err = node.left.InsertInOrder(data)
		}
	} else {
		if node.right == nil {
			err = node.setRightChild(NewTreeNode(data))
		} else {
			err = node.right.InsertInOrder(data)
		}
	}
	if err == nil {
		node.size++
	}
	return err
}

// SetLeftChild attaches a node as the left child of the current node.
// It also sets the parent of the left child to the current node.
func (node *TreeNode[T]) setLeftChild(left *TreeNode[T]) error {
	if err := validateNode(node); err != nil {
		return err
	}
	if err := validateNode(left); err != nil {
		return err
	}
	node.left = left
	left.parent = node
	return nil
}

// SetRightChild attaches a node as the right child of the current node.
// It also sets the parent of the right child to the current node.
func (node *TreeNode[T]) setRightChild(right *TreeNode[T]) error {
	if err := validateNode(node); err != nil {
		return err
	}
	if err := validateNode(right); err != nil {
		return err
	}
	node.right = right
	right.parent = node
	return nil
}

// Size returns the total number of nodes in the subtree rooted at the current node,
// including the node itself. Returns 0 if the node is nil.
func (node *TreeNode[T]) Size() int {
	if err := validateNode(node); err != nil {
		return 0
	}
	return node.size
}

// Find searches for a node with the given data in the subtree rooted at the current node.
// It returns a pointer to the found TreeNode or nil if the data is not found.
func (node *TreeNode[T]) Find(data T) *TreeNode[T] {
	if node == nil {
		return nil
	}
	if node.data == data {
		return node
	}
	if data < node.data {
		return node.left.Find(data)
	}
	return node.right.Find(data)
}
