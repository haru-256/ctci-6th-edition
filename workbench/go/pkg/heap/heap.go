// Package heap provides a generic heap implementation.
// It supports any comparable type and allows storing elements.
// The heap can be configured as either a max heap or min heap depending on the comparison function.
//
// Thread Safety:
// This implementation is thread-safe and can be used concurrently by multiple goroutines.
// All public methods use appropriate mutex locking:
//   - Read operations (Peek, Size, GetItems) use RWMutex.RLock() for concurrent reads
//   - Write operations (Insert, Pop, UpHeap, DownHeap) use RWMutex.Lock() for exclusive access
//   - Internal methods (upHeap, downHeap) do not acquire locks and should only be called
//     when the caller already holds the appropriate lock to avoid deadlocks
package heap

import (
	"cmp"
	"errors"
	"sync"
)

// ErrorIsEmpty is returned when attempting to perform operations on an empty heap.
var ErrorIsEmpty = errors.New("heap is empty")

// Heap represents a heap data structure that maintains elements in heap order.
// The heap property is determined by the comparison function provided at creation.
// For a max heap, the comparison function should return positive when a > b.
// For a min heap, the comparison function should return positive when a < b.
// T can be any type.
//
// Thread Safety:
// The Heap is thread-safe for concurrent use by multiple goroutines.
// It uses sync.RWMutex to coordinate access:
// - Multiple readers can access read-only operations (Peek, Size, GetItems) concurrently
// - Write operations (Insert, Pop, UpHeap, DownHeap) acquire exclusive locks
// - The mutex prevents race conditions and ensures heap consistency across goroutines
type Heap[T any] struct {
	items []*T
	cmpFn func(a, b *T) int
	mu    sync.RWMutex
}

// NewHeap creates and returns a new empty heap.
// The heap is initialized with zero elements and ready for use.
// The cmpFn parameter determines the heap type:
// - For max heap: return positive when a > b, negative when a < b, zero when equal
// - For min heap: return positive when a < b, negative when a > b, zero when equal
func NewHeap[T any](cmpFn func(a, b *T) int) *Heap[T] {
	return &Heap[T]{
		items: []*T{},
		cmpFn: cmpFn,
	}
}

// Insert adds a new element to the heap.
// The element is inserted at the end and then moved up to maintain the heap property.
// Time complexity: O(log n) where n is the number of elements in the heap.
func (heap *Heap[T]) Insert(item T) {
	heap.mu.Lock()
	defer heap.mu.Unlock()

	heap.items = append(heap.items, &item)
	heap.upHeap(len(heap.items) - 1)
}

// swap exchanges the elements at indices i and j in the heap.
func (h *Heap[T]) swap(i, j int) {
	h.items[i], h.items[j] = h.items[j], h.items[i]
}

func (h *Heap[T]) GetItems() []*T {
	h.mu.RLock()
	defer h.mu.RUnlock()

	itemsCopy := make([]*T, len(h.items))
	copy(itemsCopy, h.items)
	return itemsCopy
}

// upHeap moves the element at the given index up the heap until the heap property is satisfied.
// This is an internal method that doesn't acquire locks - it should only be called
// when the caller already holds the appropriate lock.
func (h *Heap[T]) upHeap(index int) {
	for {
		parentIndex := Parent(index)
		// Stop if we reach the root, or if parent satisfies heap property relative to current element
		if index == 0 || h.cmpFn(h.items[parentIndex], h.items[index]) >= 0 {
			break
		}
		// Swap with parent
		h.swap(parentIndex, index)
		// Move up to parent's position and continue
		index = parentIndex
	}
}

// UpHeap moves the element at the given index up the heap until the heap property is satisfied.
// This is the public method that acquires the necessary lock and calls the internal upHeap method.
// This is used after inserting a new element to maintain the heap property.
func (h *Heap[T]) UpHeap(index int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if index < 0 || index >= len(h.items) {
		return // Or return an error
	}
	h.upHeap(index)
}

// downHeap moves the element at the given index down the heap until the heap property is satisfied.
// This is an internal method that doesn't acquire locks - it should only be called
// when the caller already holds the appropriate lock.
// This method uses the current heap size.
func (h *Heap[T]) downHeap(index int) {
	h.downHeapWithSize(index, len(h.items))
}

// DownHeap moves the element at the given index down the heap until the heap property is satisfied.
// This is the public method that acquires the necessary lock and calls the internal downHeap method.
// This method uses the current heap size.
func (h *Heap[T]) DownHeap(index int) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.downHeap(index)
}

// downHeapWithSize moves the element at the given index down the heap until the heap property is satisfied.
// The heapSize parameter allows limiting the effective heap size, which is useful during heap sort.
func (h *Heap[T]) downHeapWithSize(index int, heapSize int) {
	// NOTE: The recursive implementation of downHeapWithSize is clear, but an iterative version can be more performant by avoiding function call overhead and eliminates the risk of stack overflow on extremely deep heaps. An iterative approach is often preferred for heap operations in production-grade code.
	for {
		l := Left(index)
		r := Right(index)
		largest := index

		if l < heapSize && h.cmpFn(h.items[l], h.items[largest]) > 0 {
			largest = l
		}
		if r < heapSize && h.cmpFn(h.items[r], h.items[largest]) > 0 {
			largest = r
		}

		if largest == index {
			break // Heap property is satisfied.
		}

		h.swap(index, largest)
		index = largest // Continue sifting down from the new position.
	}
}

// Pop removes and returns the top element from the heap.
// For a max heap, this returns the maximum element.
// For a min heap, this returns the minimum element.
// After removing the root, the heap property is restored by moving the last element
// to the root and performing a down-heap operation.
// Returns ErrorIsEmpty if the heap is empty.
// Time complexity: O(log n) where n is the number of elements in the heap.
func (h *Heap[T]) Pop() (*T, error) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if len(h.items) == 0 {
		return nil, ErrorIsEmpty
	}

	// Get the root (top element)
	top := h.items[0]
	lastIndex := len(h.items) - 1

	// Move the last element to the root
	h.items[0] = h.items[lastIndex]
	// Reduce the slice length by one
	h.items[lastIndex] = nil // Avoid memory leak by setting to nil for garbage collection
	h.items = h.items[:lastIndex]

	// Restore heap property by moving the new root down (down-heap)
	if len(h.items) > 0 {
		h.downHeap(0)
	}

	return top, nil
}

// Peek returns the top element from the heap without removing it.
// For a max heap, this returns the maximum element.
// For a min heap, this returns the minimum element.
// This allows peeking at the top element without modifying the heap.
// Returns ErrorIsEmpty if the heap is empty.
// Time complexity: O(1).
func (h *Heap[T]) Peek() (*T, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if len(h.items) == 0 {
		return nil, ErrorIsEmpty
	}
	return h.items[0], nil
}

// Size returns the number of elements currently in the heap.
// Time complexity: O(1).
func (h *Heap[T]) Size() int {
	h.mu.RLock()
	defer h.mu.RUnlock()

	return len(h.items)
}

// Left returns the index of the left child of the element at index i.
// In a binary heap stored as an array, the left child of element at index i
// is located at index 2*i + 1.
func Left(i int) int {
	return 2*i + 1
}

// Right returns the index of the right child of the element at index i.
// In a binary heap stored as an array, the right child of element at index i
// is located at index 2*(i + 1).
func Right(i int) int {
	return 2 * (i + 1)
}

// Parent returns the index of the parent of the element at index i.
// In a binary heap stored as an array, the parent of element at index i
// is located at index (i - 1) / 2.
func Parent(i int) int {
	return (i - 1) / 2
}

// BuildHeap converts an arbitrary array into a heap.
// This function performs the "heapify" operation by calling downHeap
// on all non-leaf nodes, starting from the last parent node and working upwards.
// The heap property (max or min) is determined by the comparison function.
// Time complexity: O(n) where n is the number of elements in the heap.
func BuildHeap[T cmp.Ordered](arr []*T, cmpFn func(a, b *T) int) *Heap[T] {
	heap := NewHeap(cmpFn)
	heap.items = arr
	size := heap.Size()
	for i := size/2 - 1; i >= 0; i-- {
		heap.downHeapWithSize(i, size)
	}
	return heap
}

func BuildMaxHeap[T cmp.Ordered](arr []*T) *Heap[T] {
	return BuildHeap(arr, maxCmp[T])
}

func BuildMinHeap[T cmp.Ordered](arr []*T) *Heap[T] {
	return BuildHeap(arr, minCmp[T])
}

// HeapSort sorts the elements in the heap using the heap sort algorithm.
// The algorithm works by:
// 1. Building a heap from the input array
// 2. Repeatedly extracting the top element and placing it at the end
// 3. Reducing the heap size and restoring the heap property
// After sorting, the heap's underlying array will contain elements in order
// determined by the comparison function (ascending for max heap, descending for min heap).
// Time complexity: O(n log n) where n is the number of elements.
// Space complexity: O(1) as it sorts in-place.
//
//	func HeapSort[T any](heap *Heap[T]) {
//		BuildHeap(heap)
//		heapSize := heap.Size()
//		for i := heapSize - 1; i > 0; i-- {
//			heap.swap(0, i)
//			heapSize--
//			heap.downHeapWithSize(0, heapSize)
//		}
//	}
func HeapSort[T cmp.Ordered](arr []*T) []*T {
	heap := BuildMaxHeap(arr)
	heapSize := heap.Size()
	for i := heapSize - 1; i > 0; i-- {
		heap.swap(0, i)
		heapSize--
		heap.downHeapWithSize(0, heapSize)
	}
	return heap.items
}

// Node represents a key-value pair stored in the heap.
// K is the key type used for comparison and maintaining heap order.
// V is the value type associated with each key.
type Node[K cmp.Ordered, V any] struct {
	Key   K // The key used for heap ordering
	Value V // The associated value
}

// NewMaxHeap creates a new max heap for ordered types.
// This is a convenience function for creating max heaps with ordered types.
func NewMaxHeap[T cmp.Ordered]() *Heap[T] {
	return NewHeap(maxCmp[T])
}

// NewMinHeap creates a new min heap for ordered types.
// This is a convenience function for creating min heaps with ordered types.
func NewMinHeap[T cmp.Ordered]() *Heap[T] {
	return NewHeap(minCmp[T])
}

func maxCmp[T cmp.Ordered](a, b *T) int {
	if *a > *b {
		return 1
	} else if *a < *b {
		return -1
	}
	return 0
}

func minCmp[T cmp.Ordered](a, b *T) int {
	if *a < *b {
		return 1
	} else if *a > *b {
		return -1
	}
	return 0
}
