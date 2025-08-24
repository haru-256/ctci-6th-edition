/*
Package priorityqueue provides a generic priority queue implementation using a heap data structure.

The priority queue maintains items in order of their priority, where higher priority numbers
indicate higher precedence (max-heap behavior). Items with the same priority are ordered by
their insertion time, with earlier items having higher precedence.

# Features

- Generic implementation supporting any comparable type
- Max-heap behavior (higher priority numbers are served first)
- Time-based tie-breaking for items with equal priority
- Efficient O(log n) insertion and removal operations
- O(n) priority update operations with O(log n) rebalancing
- Thread-safe for read operations

# Performance Characteristics

- Insert: O(log n)
- Pop: O(log n)
- Update: O(n) for search + O(log n) for rebalancing
- Space: O(n)

# Basic Usage

	// Create a priority queue for strings
	pq := priorityqueue.NewPriorityQueue[string](priorityqueue.PriorityCmp[string])

	// Insert items with priorities
	pq.Insert("low priority task", 1)
	pq.Insert("high priority task", 10)
	pq.Insert("medium priority task", 5)

	// Pop items in priority order (highest first)
	for {
		task, err := pq.Pop()
		if err != nil {
			break // Queue is empty
		}
		fmt.Printf("Processing: %s (priority: %d)\n", task.Value, task.Priority)
	}

# Advanced Usage

	// Create a priority queue for custom types
	type Job struct {
		ID   int
		Name string
	}

	jobQueue := priorityqueue.NewPriorityQueue[Job](priorityqueue.PriorityCmp[Job])

	// Insert jobs with different priorities
	jobQueue.Insert(Job{ID: 1, Name: "backup"}, 3)
	jobQueue.Insert(Job{ID: 2, Name: "critical-update"}, 10)
	jobQueue.Insert(Job{ID: 3, Name: "maintenance"}, 1)

	// Update priority of an existing job
	err := jobQueue.Update(Job{ID: 1, Name: "backup"}, 8)
	if err != nil {
		fmt.Printf("Job not found: %v\n", err)
	}

# Custom Comparison Functions

For different ordering behaviors, you can provide custom comparison functions:

	// Min-heap behavior (lower numbers first)
	minHeapCmp := func(a, b *priorityqueue.Task[string]) int {
		if a.Priority > b.Priority {
			return -1
		} else if a.Priority < b.Priority {
			return 1
		}
		// Time-based tie breaking
		if a.Time.Before(b.Time) {
			return 1
		} else if a.Time.After(b.Time) {
			return -1
		}
		return 0
	}

	minPQ := priorityqueue.NewPriorityQueue[string](minHeapCmp)

# Thread Safety

This implementation is thread-safe for concurrent read operations but requires external
synchronization for concurrent modifications. Use sync.Mutex or sync.RWMutex to protect
concurrent access when multiple goroutines are modifying the queue.

	var mu sync.Mutex

	// Safe concurrent insertion
	mu.Lock()
	pq.Insert("task", 5)
	mu.Unlock()

# Error Handling

The package defines specific errors for different failure conditions:

- ErrNotFound: Returned when trying to update an item that doesn't exist in the queue
- heap.ErrorIsEmpty: Returned when trying to pop from an empty queue

Always check for errors when calling Pop() and Update() methods.
*/
package priorityqueue
