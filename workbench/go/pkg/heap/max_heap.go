// Package heap provides a generic max heap implementation.
// It supports any comparable key type and allows storing associated values.
package heap

import (
	"cmp"
	"errors"
)

// ErrorIsEmpty is returned when attempting to perform operations on an empty heap.
var ErrorIsEmpty = errors.New("heap is empty")

// MaxHeap represents a max heap data structure that maintains elements in heap order.
// The heap property ensures that the parent node is always greater than or equal to its children.
// K must be an ordered type (supporting comparison operations).
// V must be a comparable type (supporting equality operations).
type MaxHeap[K cmp.Ordered, V comparable] struct {
	values []*Node[K, V]
}

// NewMaxHeap creates and returns a new empty max heap.
// The heap is initialized with zero elements and ready for use.
func NewMaxHeap[K cmp.Ordered, V comparable]() *MaxHeap[K, V] {
	return &MaxHeap[K, V]{
		values: []*Node[K, V]{},
	}
}

// Insert adds a new key-value pair to the heap.
// The element is inserted at the end and then moved up to maintain the heap property.
// Time complexity: O(log n) where n is the number of elements in the heap.
func (heap *MaxHeap[K, V]) Insert(key K, value V) {
	heap.values = append(heap.values, &Node[K, V]{Key: key, Value: value})
	heap.upHeap(heap.Size() - 1)
}

// swap exchanges the elements at indices i and j in the heap.
func (h *MaxHeap[K, V]) swap(i, j int) {
	h.values[i], h.values[j] = h.values[j], h.values[i]
}

// upHeap moves the element at the given index up the heap until the heap property is satisfied.
// This is used after inserting a new element to maintain the max heap property.
func (h *MaxHeap[K, V]) upHeap(index int) {
	for {
		parentIndex := Parent(index)
		// Stop if we reach the root, or if parent is greater than or equal to current element
		if index == 0 || h.values[parentIndex].Key >= h.values[index].Key {
			break
		}
		// Swap with parent
		h.swap(parentIndex, index)
		// Move up to parent's position and continue
		index = parentIndex
	}
}

// downHeap moves the element at the given index down the heap until the heap property is satisfied.
// This method uses the current heap size.
func (h *MaxHeap[K, V]) downHeap(index int) {
	h.downHeapWithSize(index, h.Size())
}

// downHeapWithSize moves the element at the given index down the heap until the heap property is satisfied.
// The heapSize parameter allows limiting the effective heap size, which is useful during heap sort.
func (h *MaxHeap[K, V]) downHeapWithSize(index int, heapSize int) {
	// NOTE: The recursive implementation of downHeapWithSize is clear, but an iterative version can be more performant by avoiding function call overhead and eliminates the risk of stack overflow on extremely deep heaps. An iterative approach is often preferred for heap operations in production-grade code.
	for {
		l := Left(index)
		r := Right(index)
		largest := index

		if l < heapSize && h.values[l].Key > h.values[largest].Key {
			largest = l
		}
		if r < heapSize && h.values[r].Key > h.values[largest].Key {
			largest = r
		}

		if largest == index {
			break // Heap property is satisfied.
		}

		h.swap(index, largest)
		index = largest // Continue sifting down from the new position.
	}
}

// Pop removes and returns the maximum element (root) from the heap.
// After removing the root, the heap property is restored by moving the last element
// to the root and performing a down-heap operation.
// Returns ErrorIsEmpty if the heap is empty.
// Time complexity: O(log n) where n is the number of elements in the heap.
func (h *MaxHeap[K, V]) Pop() (*Node[K, V], error) {
	if len(h.values) == 0 {
		return nil, ErrorIsEmpty
	}

	// Get the root (maximum value)
	max := h.values[0]
	lastIndex := len(h.values) - 1

	// Move the last element to the root
	h.values[0] = h.values[lastIndex]
	// Reduce the slice length by one
	h.values[lastIndex] = nil // Avoid memory leak 長さ（Length）」を lastIndex に変更するだけです。元となっている配列や、スライスのCapacityには何の変化がないため、nilを設定してガーベジコレクションの対象にします。
	h.values = h.values[:lastIndex]

	// Restore heap property by moving the new root down (down-heap)
	if len(h.values) > 0 {
		h.downHeap(0)
	}

	return max, nil
}

// Max returns the maximum element (root) from the heap without removing it.
// This allows peeking at the largest element without modifying the heap.
// Returns ErrorIsEmpty if the heap is empty.
// Time complexity: O(1).
func (h *MaxHeap[K, V]) Max() (*Node[K, V], error) {
	if len(h.values) == 0 {
		return nil, ErrorIsEmpty
	}
	return h.values[0], nil
}

// Size returns the number of elements currently in the heap.
// Time complexity: O(1).
func (h *MaxHeap[K, V]) Size() int {
	return len(h.values)
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

// BuildMaxHeap converts an arbitrary array into a max heap.
// This function performs the "heapify" operation by calling downHeap
// on all non-leaf nodes, starting from the last parent node and working upwards.
// Time complexity: O(n) where n is the number of elements in the heap.
func BuildMaxHeap[K cmp.Ordered, V comparable](heap *MaxHeap[K, V]) {
	for i := heap.Size()/2 - 1; i >= 0; i-- {
		heap.downHeapWithSize(i, heap.Size())
	}
}

// HeapSort sorts the elements in the heap in ascending order using the heap sort algorithm.
// The algorithm works by:
// 1. Building a max heap from the input array
// 2. Repeatedly extracting the maximum element (root) and placing it at the end
// 3. Reducing the heap size and restoring the heap property
// After sorting, the heap's underlying array will contain elements in ascending order.
// Time complexity: O(n log n) where n is the number of elements.
// Space complexity: O(1) as it sorts in-place.
func HeapSort[K cmp.Ordered, V comparable](heap *MaxHeap[K, V]) {
	BuildMaxHeap(heap)
	heapSize := heap.Size()
	for i := heapSize - 1; i > 0; i-- {
		heap.swap(0, i)
		heapSize--
		heap.downHeapWithSize(0, heapSize)
	}
}

// Node represents a key-value pair stored in the heap.
// K is the key type used for comparison and maintaining heap order.
// V is the value type associated with each key.
type Node[K cmp.Ordered, V comparable] struct {
	Key   K // The key used for heap ordering
	Value V // The associated value
}
