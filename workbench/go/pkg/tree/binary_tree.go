// Package tree provides a generic binary search tree implementation.
package tree

import (
	"encoding/binary"
	"errors"
	"hash"
	"hash/fnv"
	"math"
	"sync"
)

// ErrorUnsupportedValueType is returned when getHash is called with a value
// type that cannot be converted to a byte slice for hashing.
var ErrorUnsupportedValueType = errors.New("unsupported value type for hashing")

// hasherPool is a pool of FNV-1a hashers to avoid allocations in getHash.
// This provides thread-safe access to reusable hash.Hash64 instances,
// improving performance by reducing garbage collection pressure.
var hasherPool = sync.Pool{
	New: func() any {
		return fnv.New64a()
	},
}

// BinaryTree represents a thread-safe binary search tree.
// The keys are uint64 hashes of the values `V`, which must be comparable.
// It uses FNV-1a hashing to convert values to keys and maintains BST properties.
type BinaryTree[V comparable] struct {
	root *Node[uint64, V]
	size int
	mu   sync.RWMutex
}

// NewBinaryTree creates and returns a new empty BinaryTree.
// The tree is initialized with no root node and zero size.
func NewBinaryTree[V comparable]() (*BinaryTree[V], error) {
	return &BinaryTree[V]{
		root: nil,
		size: 0,
	}, nil
}

// Size returns the total number of nodes currently in the tree.
// This method is thread-safe.
func (tree *BinaryTree[V]) Size() int {
	tree.mu.RLock()
	defer tree.mu.RUnlock()
	return tree.size
}

// InsertInOrder inserts a new value into the binary search tree.
// It calculates a hash of the value to use as the key, then inserts the
// new node while maintaining the binary search tree property.
// The tree's size is incremented on successful insertion.
// This method is thread-safe.
func (tree *BinaryTree[V]) InsertInOrder(value V) error {
	tree.mu.Lock()
	defer tree.mu.Unlock()

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

// Delete removes a node with the specified value from the binary search tree.
// It returns the hash key of the deleted value and an error if the operation fails.
// If the tree is empty, it returns ErrorNodeIsNil.
// If the node is not found, it returns the key but doesn't modify the tree.
// The tree's size is decremented only on successful deletion.
// This method is thread-safe.
func (tree *BinaryTree[V]) Delete(value V) (uint64, error) {
	tree.mu.Lock()
	defer tree.mu.Unlock()

	if tree.root == nil {
		return 0, ErrorNodeIsNil
	}

	key, err := tree.getHash(value)
	if err != nil {
		return 0, err
	}

	_root, err := tree.root.delete(key, value)
	if err != nil {
		if errors.Is(err, ErrorNodeNotFound) {
			// If the node is not found, the size is not changed
			return key, nil
		}
		return 0, err // Other errors
	}
	tree.root = _root
	// Decrement size only on successful deletion
	tree.size--
	return key, nil
}

// Find searches for a node with the given value in the tree.
// It returns a pointer to the found Node or nil if the value is not found.
// If the tree is empty, it returns ErrorNodeIsNil.
// This method is thread-safe and uses a read lock for concurrent access.
func (tree *BinaryTree[V]) Find(value V) (*Node[uint64, V], error) {
	tree.mu.RLock()
	defer tree.mu.RUnlock()

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
// It uses a sync.Pool to reuse hasher objects, making it safe for concurrent use
// and avoiding allocations on each call.
// It supports int, float64, and string types. For other types,
// it returns ErrorUnsupportedValueType.
func (tree *BinaryTree[V]) getHash(value V) (uint64, error) {
	// Get a hasher from the pool and defer returning it
	hasher := hasherPool.Get().(hash.Hash64)
	defer hasherPool.Put(hasher)
	hasher.Reset()

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
	hasher.Write(b)

	// Return the computed 64-bit hash.
	return hasher.Sum64(), nil
}
