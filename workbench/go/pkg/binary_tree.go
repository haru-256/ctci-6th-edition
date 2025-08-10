package pkg

import (
	"cmp"
	"errors"
)

var NodeIsNilError = errors.New("node is nil")

type TreeNode[T cmp.Ordered] struct {
	data   T
	parent *TreeNode[T]
	left   *TreeNode[T]
	right  *TreeNode[T]
	size   int
}

func NewTreeNode[T cmp.Ordered](d T) *TreeNode[T] {
	return &TreeNode[T]{
		data: d,
		size: 1,
	}
}

func validateNode[T cmp.Ordered](node *TreeNode[T]) error {
	if node == nil {
		return NodeIsNilError
	}
	return nil
}

func (node *TreeNode[T]) InsertInOrder(data T) error {
	if err := validateNode(node); err != nil {
		return err
	}

	if data <= node.data {
		if node.left == nil {
			node.SetLeftChild(NewTreeNode(data))
		} else {
			node.left.InsertInOrder(data)
		}
	} else {
		if node.right == nil {
			node.SetRightChild(NewTreeNode(data))
		} else {
			node.right.InsertInOrder(data)
		}
	}
	node.size++
	return nil
}

func (node *TreeNode[T]) SetLeftChild(left *TreeNode[T]) error {
	if err := validateNode(node); err != nil {
		return err
	}
	if err := validateNode(left); err != nil {
		return err
	}
	node.left = left
	if left != nil {
		left.parent = node
	}
	return nil
}

func (node *TreeNode[T]) SetRightChild(right *TreeNode[T]) error {
	if err := validateNode(node); err != nil {
		return err
	}
	if err := validateNode(right); err != nil {
		return err
	}
	node.right = right
	if right != nil {
		right.parent = node
	}
	return nil
}

func (node *TreeNode[T]) Size() int {
	if err := validateNode(node); err != nil {
		return 0
	}
	return node.size
}

func (node *TreeNode[T]) Find(data T) *TreeNode[T] {
	if err := validateNode(node); err != nil {
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
