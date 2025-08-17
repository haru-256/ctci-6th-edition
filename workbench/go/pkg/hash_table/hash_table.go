// Package hash_table provides a thread-safe hash table implementation using chaining for collision resolution.
package hashtable

import (
	"encoding/binary"
	"errors"
	"hash"
	"hash/fnv"
	"math"
	"sync"

	l "github.com/haru-256/ctci-6th-edition/pkg/linked_list"
)

var (
	// ErrorAlreadyExists is returned when attempting to insert a value that already exists in the hash table.
	ErrorAlreadyExists = errors.New("value already exists in the hash table")
	// ErrorUnsupportedValueType is returned when attempting to hash a value of an unsupported type.
	ErrorUnsupportedValueType = errors.New("unsupported value type for hashing")
	// ErrorNodeNotFound is returned when a node is not found in the hash table.
	ErrorNodeNotFound = errors.New("node not found in the hash table")
)

// hasherPool is a pool of FNV-1a hashers to avoid allocations in getHash.
// This provides thread-safe access to reusable hash.Hash64 instances,
// improving performance by reducing garbage collection pressure.
var hasherPool = sync.Pool{
	New: func() any {
		return fnv.New64a()
	},
}

// HashChainTable implements a thread-safe hash table using chaining for collision resolution.
// It uses linked lists to handle hash collisions and supports generic types.
// The hash table is safe for concurrent access with read-write mutex protection.
type HashChainTable[T comparable] struct {
	// Table is an array of linked lists, where each bucket can contain multiple values
	Table []*l.LinkedList[T]
	// MaxSize is the maximum number of buckets in the hash table
MaxSize int
	// size tracks the total number of elements currently stored in the hash table
	size int
	// mu provides thread-safe access to the hash table
	mu sync.RWMutex
}

// NewHashChainTable creates and returns a new hash table with the specified maximum size.
// The maxSize parameter determines the number of buckets in the hash table.
// All buckets are initially empty (nil).
func NewHashChainTable[T comparable](maxSize int64) *HashChainTable[T] {
	if maxSize <= 0 {
		panic("hashtable: maxSize must be positive")
	}
	return &HashChainTable[T]{
		Table:   make([]*l.LinkedList[T], maxSize),
		MaxSize: maxSize,
		size:    0,
	}
}
}

// Size returns the total number of elements currently stored in the hash table.
// This method is thread-safe and uses a read lock for concurrent access.
func (table *HashChainTable[T]) Size() int {
	table.mu.RLock()
	defer table.mu.RUnlock()
	return table.size
}

// Insert adds a new value to the hash table.
// If the value already exists, it returns ErrorAlreadyExists.
// If the value type is not supported for hashing, it returns ErrorUnsupportedValueType.
// This method is thread-safe and uses a write lock for concurrent access.
func (table *HashChainTable[T]) Insert(value T) error {
	table.mu.Lock()
	defer table.mu.Unlock()

	hash, err := table.getHash(value)
	if err != nil {
		return err
	}

	index := hash % uint64(table.MaxSize)
	if table.Table[index] == nil {
		table.Table[index] = l.NewLinkedList[T]()
	}
	if table.Table[index].Search(value) != nil {
		return ErrorAlreadyExists
	}

	table.Table[index].Prepend(value)
	table.size++
	return nil
}

// Search looks for a value in the hash table and returns the corresponding node.
// If the value is found, it returns the node containing the value.
// If the value is not found, it returns nil for the node.
// If the value type is not supported for hashing, it returns ErrorUnsupportedValueType.
// This method is thread-safe and uses a read lock for concurrent access.
func (table *HashChainTable[T]) Search(value T) (*l.Node[T], error) {
	table.mu.RLock()
	defer table.mu.RUnlock()

	hash, err := table.getHash(value)
	if err != nil {
		return nil, err
	}

	index := hash % uint64(table.MaxSize)
	if table.Table[index] == nil {
		return nil, nil
	}

	return table.Table[index].Search(value), nil
}

// Delete removes a value from the hash table.
// If the value exists, it is removed and the size is decremented.
// If the value does not exist, the operation occurs ErrorNodeNotFound.
// If the value type is not supported for hashing, it returns ErrorUnsupportedValueType.
// This method is thread-safe and uses a write lock for concurrent access.
func (table *HashChainTable[T]) Delete(value T) error {
	table.mu.Lock()
	defer table.mu.Unlock()

	hash, err := table.getHash(value)
	if err != nil {
		return err
	}

	index := hash % uint64(table.MaxSize)
	if table.Table[index] == nil {
		return ErrorNodeNotFound
	}

	if err = table.Table[index].Delete(value); err != nil {
		return err
	}

	table.size--
	if table.Table[index].Head == nil { // if list is empty, remove bucket for garbage collection
		table.Table[index] = nil
	}
	return nil
}

// getHash computes and returns the FNV-1a hash of a given value.
// It uses a sync.Pool to reuse hasher objects, making it safe for concurrent use
// and avoiding allocations on each call.
// It supports int, float64, and string types. For other types,
// it returns ErrorUnsupportedValueType.
func (table *HashChainTable[T]) getHash(value T) (uint64, error) {
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
