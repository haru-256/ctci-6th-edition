package trietree

import (
	"errors"
	"sync"
)

var (
	ErrKeyNotFound = errors.New("key not found in trie tree")
)

type TrieTree[K comparable, V any] struct {
	root *node[K, V]
	mu   sync.RWMutex
}

type node[K comparable, V any] struct {
	children map[K]*node[K, V]
	value    V    // Value associated with the key ending at this node
	isEnd    bool // Indicates if this node marks the end of a key
}

func NewTrieTree[K comparable, V any]() *TrieTree[K, V] {
	return &TrieTree[K, V]{
		root: &node[K, V]{
			children: make(map[K]*node[K, V]),
			isEnd:    false,
		},
	}
}

func (t *TrieTree[K, V]) Insert(key []K, value V) {
	t.mu.Lock()
	defer t.mu.Unlock()

	current := t.root
	for _, k := range key {
		if _, exists := current.children[k]; !exists {
			current.children[k] = &node[K, V]{
				children: make(map[K]*node[K, V]),
				isEnd:    false,
			}
		}
		current = current.children[k]
	}
	current.value = value
	current.isEnd = true
}

func (t *TrieTree[K, V]) Search(key []K) (V, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	current := t.root
	for _, k := range key {
		if _, exists := current.children[k]; !exists {
			var zero V
			return zero, false
		}
		current = current.children[k]
	}
	if current.isEnd {
		return current.value, true
	}
	var zero V
	return zero, false
}

func (t *TrieTree[K, V]) StartsWith(key []K) bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	current := t.root
	for _, k := range key {
		if _, exists := current.children[k]; !exists {
			return false
		}
		current = current.children[k]
	}
	return true
}

func (t *TrieTree[K, V]) Delete(key []K) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if len(key) == 0 {
		return ErrKeyNotFound
	}

	// First, verify the key exists
	current := t.root
	for _, k := range key {
		if _, exists := current.children[k]; !exists {
			return ErrKeyNotFound
		}
		current = current.children[k]
	}
	if !current.isEnd {
		return ErrKeyNotFound
	}

	// Mark the end node as no longer being an end
	current.isEnd = false
	var zero V
	current.value = zero

	// Recursively delete nodes that are no longer needed
	t.deleteRecursive(t.root, key, 0)
	return nil
}

// deleteRecursive removes nodes that are no longer needed
func (t *TrieTree[K, V]) deleteRecursive(current *node[K, V], key []K, index int) bool {
	if index == len(key) {
		// We've reached the end of the key
		return !current.isEnd && len(current.children) == 0
	}

	k := key[index]
	child, exists := current.children[k]
	if !exists {
		return false
	}

	// Recursively check if the child should be deleted
	shouldDelete := t.deleteRecursive(child, key, index+1)

	if shouldDelete {
		delete(current.children, k)
	}

	// Return true if this node should be deleted
	// (it's not an end node and has no children)
	return !current.isEnd && len(current.children) == 0
}

func (t *TrieTree[K, V]) Size() int {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.sizeRecursive(t.root)
}

// DFS (Depth-First Search) to count the number of keys in the trie
func (t *TrieTree[K, V]) sizeRecursive(current *node[K, V]) int {
	if current == nil {
		return 0
	}
	size := 0
	if current.isEnd {
		size++
	}
	for _, child := range current.children {
		size += t.sizeRecursive(child)
	}
	return size
}

func (t *TrieTree[K, V]) IsEmpty() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return !t.root.isEnd && len(t.root.children) == 0
}

func (t *TrieTree[K, V]) Keys() [][]K {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var results [][]K
	var currentKey []K
	t.collectKeys(t.root, currentKey, &results)
	return results
}

func (t *TrieTree[K, V]) KeysWithPrefix(prefix []K) ([][]K, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	current := t.root
	for _, k := range prefix {
		if _, exists := current.children[k]; !exists {
			return nil, ErrKeyNotFound
		}
		current = current.children[k]
	}
	var results [][]K
	var currentKey []K
	t.collectKeys(current, append(currentKey, prefix...), &results)
	return results, nil
}

func (t *TrieTree[K, V]) collectKeys(current *node[K, V], currentKey []K, results *[][]K) {
	if current.isEnd {
		// Make a copy of currentKey to avoid mutation issues
		keyCopy := make([]K, len(currentKey))
		copy(keyCopy, currentKey)
		*results = append(*results, keyCopy)
	}
	for k, child := range current.children {
		// append function can return a reference to the same underlying array, so we need to be careful
		nextKey := make([]K, len(currentKey)+1)
		copy(nextKey, currentKey)
		nextKey[len(currentKey)] = k
		t.collectKeys(child, nextKey, results)
	}
}
