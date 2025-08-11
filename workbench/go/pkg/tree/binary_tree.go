// Package tree provides a generic binary search tree implementation.
package tree

import (
	"encoding/binary"
	"errors"
	"hash"
	"hash/fnv"
	"math"
)

// ErrorUnsupportedValueType is returned when getHash is called with a value
// type that cannot be converted to a byte slice for hashing.
var ErrorUnsupportedValueType = errors.New("unsupported value type for hashing")

// BinaryTree represents a binary search tree.
// The keys are uint64 hashes of the values `V`, which must be comparable.
type BinaryTree[V comparable] struct {
	root   *Node[uint64, V]
	size   int
	hasher hash.Hash64
}

// NewBinaryTree creates and returns a new BinaryTree, optionally with a root node.
func NewBinaryTree[V comparable]() (*BinaryTree[V], error) {
	tree := &BinaryTree[V]{root: nil, size: 0, hasher: fnv.New64a()}
	return tree, nil
}

// Size returns the total number of nodes in the tree.
func (tree *BinaryTree[V]) Size() int {
	return tree.size
}

// InsertInOrder inserts a new value into the binary search tree.
// It calculates a hash of the value to use as the key, then inserts the
// new node while maintaining the binary search tree property.
// The tree's size is incremented on successful insertion.
func (tree *BinaryTree[V]) InsertInOrder(value V) error {
	key, err := tree.getHash(value)
	if err != nil {
		return err
	}
	node := NewNode(key, value)

	if tree.root == nil {
		tree.root = node
		tree.size++
		return nil
	}

	err = tree.root.insertInOrder(key, value)
	if err == nil {
		tree.size++
	}
	return err
}

// Find searches for a node with the given value in the tree.
// It returns a pointer to the found Node or nil if the value is not found.
func (tree *BinaryTree[V]) Find(value V) (*Node[uint64, V], error) {
	if tree.root == nil {
		return nil, ErrorNodeIsNil
	}
	key, err := tree.getHash(value)
	if err != nil {
		return nil, err
	}
	return tree.root.find(key, value), nil
}

// getHash computes and returns the FNV-1a hash of a given value.
// It supports int, float64, and string types. For other types,
// it returns an ErrorUnsupportedValueType.
func (tree *BinaryTree[V]) getHash(value V) (uint64, error) {
	// Reset the hasher for a new computation.
	tree.hasher.Reset()

	// Convert the value to a byte slice based on its type.
	var b []byte
	switch v := any(value).(type) {
	case int:
		var buf [binary.MaxVarintLen64]byte
		n := binary.PutVarint(buf[:], int64(v))
		b = buf[:n]
	case float64:
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], math.Float64bits(v))
		b = buf[:]
	case string:
		b = []byte(v)
	default:
		return 0, ErrorUnsupportedValueType
	}

	// Write the byte slice to the hasher.
	// This Write is guaranteed not to return an error.
	tree.hasher.Write(b)

	// Return the computed 64-bit hash.
	return tree.hasher.Sum64(), nil
}
