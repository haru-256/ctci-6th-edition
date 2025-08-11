// Package tree provides a generic binary search tree implementation.
package tree

import (
	"cmp"
	"errors"
)

// ErrorNodeIsNil is returned when an operation is attempted on a nil Node.
var ErrorNodeIsNil = errors.New("node is nil")

// ErrorNodeNotFound is returned when a delete operation cannot find the target node.
var ErrorNodeNotFound = errors.New("node not found")

// Node represents a node in a binary search tree.
// It holds a generic key `K` that must be an ordered type, and a generic value `V`
// that must be a comparable type.
type Node[K cmp.Ordered, V comparable] struct {
	key    K
	value  V
	parent *Node[K, V]
	left   *Node[K, V]
	right  *Node[K, V]
}

// NewNode creates and returns a new Node with the given key and value.
func NewNode[K cmp.Ordered, V comparable](key K, value V) *Node[K, V] {
	return &Node[K, V]{
		key:   key,
		value: value,
	}
}

// insertInOrder inserts a new node into the binary search tree rooted at the current node.
// It maintains the binary search tree property: nodes with smaller or equal keys go to the left,
// and nodes with greater keys go to the right.
func (node *Node[K, V]) insertInOrder(key K, value V) error {
	if node == nil {
		return ErrorNodeIsNil
	}

	var err error
	if key <= node.key {
		if node.left == nil {
			err = node.setLeftChild(NewNode(key, value))
		} else {
			err = node.left.insertInOrder(key, value)
		}
	} else {
		if node.right == nil {
			err = node.setRightChild(NewNode(key, value))
		} else {
			err = node.right.insertInOrder(key, value)
		}
	}
	return err
}

// delete removes a node with the specified key and value from the binary search tree.
// It returns the new root of the subtree after deletion and an error if the operation fails.
// For nodes with two children, it replaces the node with its in-order successor.
// Returns ErrorNodeNotFound if the target node is not found.
func (node *Node[K, V]) delete(key K, value V) (*Node[K, V], error) {
	if node == nil {
		return nil, ErrorNodeNotFound
	}

	if key < node.key {
		_node, err := node.left.delete(key, value)
		if err != nil {
			return nil, err
		}
		node.left = _node
		if node.left != nil {
			node.left.parent = node
		}
	} else if key > node.key {
		_node, err := node.right.delete(key, value)
		if err != nil {
			return nil, err
		}
		node.right = _node
		if node.right != nil {
			node.right.parent = node
		}
	} else if node.value == value { // key == node.key
		// Delete this node
		if node.left == nil {
			return node.right, nil
		} else if node.right == nil {
			return node.left, nil
		} else {
			minNode, err := node.right.findMin()
			if err != nil {
				return nil, err
			}
			node.key = minNode.key
			node.value = minNode.value
			_node, err := node.right.delete(minNode.key, minNode.value)
			if err != nil {
				return nil, err
			}
			node.right = _node
			return node, nil
		}
	} else { // key == node.key but value is different
		// Same key but different value, search in left subtree
		_node, err := node.left.delete(key, value)
		if err != nil {
			return nil, err
		}
		node.left = _node
	}
	return node, nil
}

// findMin finds and returns the node with the minimum key in the subtree rooted at this node.
// It traverses left children until reaching the leftmost node.
// Returns ErrorNodeIsNil if called on a nil node.
func (node *Node[K, V]) findMin() (*Node[K, V], error) {
	if node == nil {
		return nil, ErrorNodeIsNil
	}
	current := node
	for current.left != nil {
		current = current.left
	}
	return current, nil
}

// setLeftChild attaches a node as the left child of the current node.
// It also sets the parent of the left child to the current node.
func (node *Node[K, V]) setLeftChild(left *Node[K, V]) error {
	if node == nil || left == nil {
		return ErrorNodeIsNil
	}
	node.left = left
	left.parent = node
	return nil
}

// setRightChild attaches a node as the right child of the current node.
// It also sets the parent of the right child to the current node.
func (node *Node[K, V]) setRightChild(right *Node[K, V]) error {
	if node == nil || right == nil {
		return ErrorNodeIsNil
	}
	node.right = right
	right.parent = node
	return nil
}

// find searches for a node with the given key and value in the subtree rooted at the current node.
// It returns a pointer to the found Node or nil if it is not found.
func (node *Node[K, V]) find(key K, value V) *Node[K, V] {
	if node == nil {
		return nil
	}

	if key < node.key {
		return node.left.find(key, value)
	}
	if key > node.key {
		return node.right.find(key, value)
	}

	// key == node.key
	if node.value == value {
		return node
	}

	// Since insertInOrder places equal keys in the left subtree,
	// we continue the search there for other nodes with the same key.
	return node.left.find(key, value)
}
