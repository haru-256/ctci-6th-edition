package priorityqueue

import (
	"fmt"
	"time"

	"github.com/haru-256/ctci-6th-edition/pkg/heap"
)

var ErrNotFound = fmt.Errorf("item not found")

type PriorityQueue[T comparable] struct {
	heap *heap.Heap[Task[T]]
}

func NewPriorityQueue[T comparable](cmpFn func(a, b *Task[T]) int) *PriorityQueue[T] {
	return &PriorityQueue[T]{
		heap: heap.NewHeap(cmpFn),
	}
}

func (pq *PriorityQueue[T]) Insert(item T, priority int) {
	task := NewTask(priority, item)
	pq.heap.Insert(task)
}

func (pq *PriorityQueue[T]) Pop() (*Task[T], error) {
	task, err := pq.heap.Pop()
	if err != nil {
		return nil, err
	}
	return task, nil
}

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

type Task[T comparable] struct {
	Priority int
	Time     time.Time
	Value    T
}

func NewTask[T comparable](priority int, value T) Task[T] {
	return Task[T]{
		Priority: priority,
		Time:     time.Now(),
		Value:    value,
	}
}

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
