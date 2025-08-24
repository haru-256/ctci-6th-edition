package priorityqueue

import (
	"fmt"
	"time"

	"github.com/haru-256/ctci-6th-edition/pkg/heap"
)

// ErrNotFound is returned when trying to update an item that doesn't exist in the priority queue.
var ErrNotFound = fmt.Errorf("item not found")

// PriorityQueue is a generic priority queue implementation that maintains items in priority order.
// It uses a heap data structure internally for efficient insertion and removal operations.
//
// The type parameter T must be comparable to support item lookups for updates.
// Items are ordered by priority (higher numbers have higher precedence) and by insertion time
// for items with equal priority.
//
// Time complexities:
//   - Insert: O(log n)
//   - Pop: O(log n)
//   - Update: O(n) for finding the item + O(log n) for rebalancing
//
// Space complexity: O(n) where n is the number of items in the queue.
type PriorityQueue[T comparable] struct {
	heap *heap.Heap[Task[T]]
}

// NewPriorityQueue creates a new priority queue with the given comparison function.
//
// The comparison function should return:
//   - negative value if a should come before b
//   - positive value if a should come after b
//   - zero if a and b have equal priority
//
// For max-heap behavior (higher priorities first), use PriorityCmp[T].
//
// Example:
//
//	pq := NewPriorityQueue[string](PriorityCmp[string])
func NewPriorityQueue[T comparable](cmpFn func(a, b *Task[T]) int) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap: heap.NewHeap(cmpFn),
	}
}

// Insert adds a new item to the priority queue with the specified priority.
//
// The item will be positioned in the queue according to its priority and insertion time.
// Higher priority numbers indicate higher precedence. Items with the same priority
// are ordered by insertion time (earlier items first).
//
// Time complexity: O(log n)
//
// Example:
//
//	pq.Insert("urgent task", 10)
//	pq.Insert("normal task", 5)
func (pq *PriorityQueue[T]) Insert(item T, priority int) {
	task := NewTask(priority, item)
	pq.heap.Insert(task)
}

// Pop removes and returns the highest priority item from the queue.
//
// Returns the task with the highest priority, or an error if the queue is empty.
// For items with equal priority, the item that was inserted first is returned.
//
// Time complexity: O(log n)
//
// Example:
//
//	task, err := pq.Pop()
//	if err != nil {
//		// Handle empty queue
//		return
//	}
//	fmt.Printf("Processing: %s (priority: %d)\n", task.Value, task.Priority)
func (pq *PriorityQueue[T]) Pop() (*Task[T], error) {
	task, err := pq.heap.Pop()
	if err != nil {
		return nil, err
	}
	return task, nil
}

// Update changes the priority of an existing item in the queue.
//
// Searches for the item with the given value and updates its priority.
// The item's position in the queue will be adjusted accordingly.
// If the item is not found, returns ErrNotFound.
// If the new priority is the same as the current priority, no operation is performed.
//
// Time complexity: O(n) for searching + O(log n) for rebalancing
//
// Example:
//
//	err := pq.Update("task1", 15) // Change priority of "task1" to 15
//	if err == ErrNotFound {
//		fmt.Println("Task not found in queue")
//	}
func (pq *PriorityQueue[T]) Update(item T, priority int) error {
	var targetTask *Task[T]
	var targetIdx int
	for idx, task := range pq.heap.GetItems() {
		if task.Value == item {
			targetTask = task
			targetIdx = idx
			break
		}
	}
	if targetTask == nil {
		return ErrNotFound
	}
	if targetTask.Priority == priority {
		return nil
	}

	toLargerThan := targetTask.Priority < priority
	toLessThan := targetTask.Priority > priority

	targetTask.Priority = priority

	if toLessThan {
		// Priority increased (was less, now higher), move up towards root
		pq.heap.DownHeap(targetIdx)
	} else if toLargerThan {
		// Priority decreased (was more, now lower), move down towards leaves
		pq.heap.UpHeap(targetIdx)
	}
	return nil
}

// Task represents an item in the priority queue with associated metadata.
//
// Each task contains:
//   - Priority: The priority level (higher numbers = higher priority)
//   - Time: Timestamp when the task was created (used for tie-breaking)
//   - Value: The actual data being stored in the queue
//
// Tasks are compared first by priority, then by creation time for items
// with equal priority (earlier tasks have higher precedence).
type Task[T comparable] struct {
	Priority int       // Priority level (higher numbers have higher precedence)
	Time     time.Time // Creation timestamp for tie-breaking
	Value    T         // The actual data stored in the task
}

// NewTask creates a new task with the specified priority and value.
//
// The task's creation time is automatically set to the current time,
// which is used for ordering items with equal priority.
//
// Example:
//
//	task := NewTask(5, "process this item")
//	fmt.Printf("Task: %s, Priority: %d, Created: %s\n",
//		task.Value, task.Priority, task.Time)
func NewTask[T comparable](priority int, value T) Task[T] {
	return Task[T]{
		Priority: priority,
		Time:     time.Now(),
		Value:    value,
	}
}

// PriorityCmp is a comparison function for tasks that implements max-heap behavior.
//
// Returns:
//   - negative value if task 'a' should come before task 'b' (higher priority)
//   - positive value if task 'a' should come after task 'b' (lower priority)
//   - zero if tasks have equal priority and time
//
// Comparison logic:
//  1. Higher priority numbers come first (max-heap)
//  2. For equal priorities, earlier creation times come first
//  3. Tasks with identical priority and time are considered equal
//
// This function can be passed to NewPriorityQueue to create a max-heap priority queue.
//
// Example:
//
//	pq := NewPriorityQueue[string](PriorityCmp[string])
//	// This creates a priority queue where higher priority numbers are served first
func PriorityCmp[T comparable](a, b *Task[T]) int {
	if a.Priority < b.Priority {
		return -1
	} else if a.Priority > b.Priority {
		return 1
	}
	// If priorities are equal, compare by time
	if a.Time.Before(b.Time) {
		return 1
	} else if a.Time.After(b.Time) {
		return -1
	}
	return 0
}
