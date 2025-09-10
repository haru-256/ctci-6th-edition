// Package queue provides a generic circular queue data structure.
package queue

import "errors"

var (
	// ErrorQueueOverflow is returned when trying to enqueue to a full queue.
	ErrorQueueOverflow = errors.New("queue overflow")
	// ErrorQueueUnderflow is returned when trying to dequeue from an empty queue.
	ErrorQueueUnderflow = errors.New("queue underflow")
)

// Queue represents a generic circular queue using a fixed-size array.
type Queue[T any] struct {
	items []*T
	size  int
	head  int
	tail  int
}

// NewQueue creates and returns a new Queue with the specified capacity.
// The queue uses a circular buffer implementation for efficient space utilization.
func NewQueue[T any](size int) *Queue[T] {
	if size <= 0 {
		panic("queue size must be greater than 0")
	}
	return &Queue[T]{
		items: make([]*T, size),
		size:  size,
		head:  0,
		tail:  0,
	}
}

// IsEmpty checks if the queue is empty.
// Returns true if there are no elements in the queue.
func (q *Queue[T]) IsEmpty() bool {
	return q.head == q.tail && q.items[q.head] == nil
}

// IsFull checks if the queue is full.
// Returns true if the queue has reached its maximum capacity.
func (q *Queue[T]) IsFull() bool {
	if q.IsEmpty() {
		return false
	}
	return q.tail%q.size == q.head
}

// Enqueue adds an item to the rear of the queue.
// Returns ErrorQueueOverflow if the queue is full.
func (q *Queue[T]) Enqueue(item *T) error {
	if q.IsFull() {
		return ErrorQueueOverflow
	}
	q.items[q.tail] = item
	q.tail = (q.tail + 1) % q.size
	return nil
}

// Dequeue removes and returns the front item from the queue.
// Returns ErrorQueueUnderflow if the queue is empty.
func (q *Queue[T]) Dequeue() (*T, error) {
	if q.IsEmpty() {
		return nil, ErrorQueueUnderflow
	}
	item := q.items[q.head]
	q.items[q.head] = nil
	q.head = (q.head + 1) % q.size
	return item, nil
}

// Peek returns the front item from the queue without removing it.
// Returns ErrorQueueUnderflow if the queue is empty.
func (q *Queue[T]) Peek() (*T, error) {
	if q.IsEmpty() {
		return nil, ErrorQueueUnderflow
	}
	return q.items[q.head], nil
}
