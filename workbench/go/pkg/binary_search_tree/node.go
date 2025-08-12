// package binary_search_tree provides a generic binary search tree implementation.
package binary_search_tree

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
// It performs a recursive search and deletion while maintaining BST properties and parent relationships.
// The method handles four cases:
// 1. Target is in left subtree (key < node.key)
// 2. Target is in right subtree (key > node.key)
// 3. Same key but different value (searches left subtree for duplicates)
// 4. Exact match found (delegates to deleteCurrentNode for removal)
// Returns the new root of the subtree after deletion and ErrorNodeNotFound if target not found.
func (node *Node[K, V]) delete(key K, value V) (*Node[K, V], error) {
	if node == nil {
		return nil, ErrorNodeNotFound
	}

	switch {
	case key < node.key:
		newLeft, err := node.left.delete(key, value)
		if err != nil {
			return nil, err
		}
		return node.updateChild(newLeft, true), nil

	case key > node.key:
		newRight, err := node.right.delete(key, value)
		if err != nil {
			return nil, err
		}
		return node.updateChild(newRight, false), nil

	case node.value != value:
		// Same key but different value, search in left subtree
		newLeft, err := node.left.delete(key, value)
		if err != nil {
			return nil, err
		}
		return node.updateChild(newLeft, true), nil

	default:
		// Found the node to delete (key == node.key && value == node.value)
		return node.deleteCurrentNode()
	}
}

// updateChild updates either the left or right child of the current node and maintains parent relationships.
// If isLeft is true, updates the left child; otherwise updates the right child.
// The parent pointer of the new child is automatically set to point to the current node.
// Returns the current node to support method chaining.
func (node *Node[K, V]) updateChild(child *Node[K, V], isLeft bool) *Node[K, V] {
	if isLeft {
		node.left = child
	} else {
		node.right = child
	}
	if child != nil {
		child.parent = node
	}
	return node
}

// deleteCurrentNode handles the deletion of the current node when it matches the target key and value.
// It implements the three standard BST deletion cases:
// 1. Node with no children (leaf): simply return nil
// 2. Node with one child: return the child to replace this node
// 3. Node with two children: replace with in-order successor and delete the successor
// This method assumes the current node is the target to be deleted.
func (node *Node[K, V]) deleteCurrentNode() (*Node[K, V], error) {
	// Case 1: No left child
	if node.left == nil {
		return node.right, nil
	}

	// Case 2: No right child
	if node.right == nil {
		return node.left, nil
	}

	// Case 3: Two children - replace with in-order successor
	successor, err := node.right.findMin()
	if err != nil {
		return nil, err
	}

	// Copy successor's data to current node
	node.key = successor.key
	node.value = successor.value

	// Delete the successor from right subtree
	newRight, err := node.right.delete(successor.key, successor.value)
	if err != nil {
		return nil, err
	}

	node.updateChild(newRight, false)

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
