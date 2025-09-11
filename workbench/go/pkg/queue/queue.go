package queue

import "errors"

var (
	// ErrorQueueOverflow is returned when trying to enqueue to a full queue.
	ErrorQueueOverflow = errors.New("queue overflow")
	// ErrorQueueUnderflow is returned when trying to dequeue from an empty queue.
	ErrorQueueUnderflow = errors.New("queue underflow")
)

// Queue represents a generic FIFO (First In, First Out) circular queue using a fixed-size array.
// It stores values of type T directly and provides thread-unsafe operations.
// The zero value is not ready to use; use NewQueue to create a new queue.
//
// The queue has a fixed capacity determined at creation time and will return
// ErrorQueueOverflow when attempting to enqueue beyond capacity, or ErrorQueueUnderflow
// when attempting to dequeue from an empty queue.
//
// The circular buffer implementation efficiently reuses array space as items are
// enqueued and dequeued, preventing the need to shift elements.
//
// Time complexity:
//   - Enqueue: O(1)
//   - Dequeue: O(1)
//   - Peek: O(1)
//   - IsEmpty/IsFull/Size/Count: O(1)
//
// Space complexity: O(n) where n is the capacity.
type Queue[T any] struct {
	items []T
	size  int // maximum number of items the queue can hold
	count int // current number of items in the queue
	head  int
	tail  int
}

// NewQueue creates and returns a new Queue with the specified capacity.
// The size parameter must be greater than 0, otherwise the function will panic.
// The queue is initially empty and ready to use.
//
// The returned queue will have a fixed capacity and will not grow beyond the
// specified size. Use Enqueue to add items and Dequeue to remove them.
// The queue follows FIFO (First In, First Out) ordering.
//
// Parameters:
//   - size: The maximum number of items the queue can hold (must be > 0)
//
// Returns:
//   - A new Queue instance ready for use
//
// Panics:
//   - If size <= 0
//
// Example:
//
//	queue := NewQueue[int](10)    // Creates a queue that can hold 10 integers
//	queue := NewQueue[string](5)  // Creates a queue that can hold 5 strings
func NewQueue[T any](size int) *Queue[T] {
	if size <= 0 {
		panic("queue size must be greater than 0")
	}
	return &Queue[T]{
		items: make([]T, size),
		size:  size,
		count: 0,
		head:  0,
		tail:  0,
	}
}

// IsEmpty checks if the queue is empty.
// Returns true if there are no elements in the queue.
func (q *Queue[T]) IsEmpty() bool {
	return q.count == 0
}

// IsFull checks if the queue is full.
// Returns true if the queue has reached its maximum capacity.
func (q *Queue[T]) IsFull() bool {
	return q.count == q.size
}

// Enqueue adds an item to the rear of the queue.
// Returns ErrorQueueOverflow if the queue is full.
// The queue follows FIFO order, so this item will be the last to be dequeued.
//
// Example:
//
//	err := queue.Enqueue(42)
//	if err != nil {
//	    // handle queue overflow
//	}
func (q *Queue[T]) Enqueue(item T) error {
	if q.IsFull() {
		return ErrorQueueOverflow
	}
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.size
	q.count++
	return nil
}

// Dequeue removes and returns the front item from the queue.
// Returns ErrorQueueUnderflow if the queue is empty.
// The queue follows FIFO (First In, First Out) order.
//
// Example:
//
//	item, err := queue.Dequeue()
//	if err != nil {
//	    // handle queue underflow
//	} else {
//	    fmt.Println(item) // use the dequeued item
//	}
func (q *Queue[T]) Dequeue() (T, error) {
	var zero T
	if q.IsEmpty() {
		return zero, ErrorQueueUnderflow
	}
	item := q.items[q.head]
	q.items[q.head] = zero // Clear the slot
	q.head = (q.head + 1) % q.size
	q.count--
	return item, nil
}

// Peek returns the front item from the queue without removing it.
// Returns ErrorQueueUnderflow if the queue is empty.
// This operation does not modify the queue.
//
// Example:
//
//	item, err := queue.Peek()
//	if err != nil {
//	    // handle empty queue
//	} else {
//	    fmt.Println(item) // use the front item without removing it
//	}
func (q *Queue[T]) Peek() (T, error) {
	if q.IsEmpty() {
		var zero T
		return zero, ErrorQueueUnderflow
	}
	return q.items[q.head], nil
}

// Size returns the maximum capacity of the queue.
// This is the size that was specified when the queue was created.
func (q *Queue[T]) Size() int {
	return q.size
}

// Count returns the current number of items in the queue.
// This value ranges from 0 (empty) to Size() (full).
func (q *Queue[T]) Count() int {
	return q.count
}
