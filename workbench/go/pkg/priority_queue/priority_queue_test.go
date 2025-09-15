package priorityqueue

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/haru-256/ctci-6th-edition/pkg/heap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

// Helper types for test data
type testItem struct {
	value    string
	priority int
}

type testUpdate struct {
	item     string
	priority int
}

// Helper function to create a priority queue with test items
func setupPriorityQueue(items []testItem) *PriorityQueue[string] {
	pq := NewPriorityQueue(PriorityCmp[string])
	for _, item := range items {
		pq.Insert(item.value, item.priority)
	}
	return pq
}

// Helper function to verify popping order
func verifyPopOrder(t *testing.T, pq *PriorityQueue[string], expectedOrder []testItem) {
	t.Helper()
	for i, expected := range expectedOrder {
		task, err := pq.Pop()
		require.NoError(t, err, "Pop %d should not return error", i)
		require.Equal(t, expected.value, task.Value, "Pop %d value mismatch", i)
		require.Equal(t, expected.priority, task.Priority, "Pop %d priority mismatch", i)
	}
}

func TestNewPriorityQueue(t *testing.T) {
	tests := []struct {
		name  string
		setup func() (*PriorityQueue[string], error)
		want  func(*testing.T, *PriorityQueue[string])
	}{
		{
			name: "create string priority queue",
			setup: func() (*PriorityQueue[string], error) {
				return NewPriorityQueue(PriorityCmp[string]), nil
			},
			want: func(t *testing.T, pq *PriorityQueue[string]) {
				require.NotNil(t, pq)
				require.NotNil(t, pq.heap)
				require.Equal(t, 0, pq.heap.Size())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq, err := tt.setup()
			require.NoError(t, err)
			tt.want(t, pq)
		})
	}
}

func TestPriorityQueue_Insert(t *testing.T) {
	tests := []struct {
		name     string
		items    []testItem
		wantSize int
	}{
		{
			name: "insert single element",
			items: []testItem{
				{"task1", 10},
			},
			wantSize: 1,
		},
		{
			name: "insert multiple elements",
			items: []testItem{
				{"task1", 10},
				{"task2", 5},
				{"task3", 15},
				{"task4", 1},
			},
			wantSize: 4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := setupPriorityQueue(tt.items)
			require.Equal(t, tt.wantSize, pq.heap.Size())
		})
	}
}

func TestPriorityQueue_Pop(t *testing.T) {
	tests := []struct {
		name          string
		setupItems    []testItem
		expectedOrder []testItem
		wantError     error
	}{
		{
			name:          "pop from empty queue",
			setupItems:    nil,
			expectedOrder: nil,
			wantError:     heap.ErrorIsEmpty,
		},
		{
			name: "pop in priority order",
			setupItems: []testItem{
				{"low", 10},
				{"high", 30},
				{"medium", 20},
			},
			expectedOrder: []testItem{
				{"high", 30},
				{"medium", 20},
				{"low", 10},
			},
			wantError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := setupPriorityQueue(tt.setupItems)

			if tt.wantError != nil {
				// Test error case
				_, err := pq.Pop()
				require.Equal(t, tt.wantError, err)
				return
			}

			// Test successful pops and verify order
			verifyPopOrder(t, pq, tt.expectedOrder)

			// Queue should be empty now
			require.Equal(t, 0, pq.heap.Size())

			// Another pop should return error
			_, err := pq.Pop()
			require.Equal(t, heap.ErrorIsEmpty, err)
		})
	}
}

func TestPriorityQueue_Update(t *testing.T) {
	tests := []struct {
		name           string
		setupItems     []testItem
		updateItem     string
		updatePriority int
		expectedOrder  []testItem
		wantError      error
	}{
		{
			name: "update to higher priority",
			setupItems: []testItem{
				{"task1", 10},
				{"task2", 5},
				{"task3", 15},
			},
			updateItem:     "task2",
			updatePriority: 20,
			expectedOrder: []testItem{
				{"task2", 20},
				{"task3", 15},
				{"task1", 10},
			},
			wantError: nil,
		},
		{
			name: "update to lower priority",
			setupItems: []testItem{
				{"task1", 10},
				{"task2", 15},
				{"task3", 5},
			},
			updateItem:     "task2",
			updatePriority: 3,
			expectedOrder: []testItem{
				{"task1", 10},
				{"task3", 5},
				{"task2", 3},
			},
			wantError: nil,
		},
		{
			name: "update same priority",
			setupItems: []testItem{
				{"task1", 10},
			},
			updateItem:     "task1",
			updatePriority: 10,
			expectedOrder: []testItem{
				{"task1", 10},
			},
			wantError: nil,
		},
		{
			name: "update non-existent item",
			setupItems: []testItem{
				{"task1", 10},
				{"task2", 5},
			},
			updateItem:     "nonexistent",
			updatePriority: 1,
			expectedOrder:  nil,
			wantError:      ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := setupPriorityQueue(tt.setupItems)

			// Update
			err := pq.Update(tt.updateItem, tt.updatePriority)

			if tt.wantError != nil {
				require.Equal(t, tt.wantError, err)
				return
			}

			require.NoError(t, err)

			// Verify order
			verifyPopOrder(t, pq, tt.expectedOrder)
		})
	}
}

func TestTask_NewTask(t *testing.T) {
	tests := []struct {
		name     string
		priority int
		value    string
		want     func(*testing.T, Task[string], time.Time)
	}{
		{
			name:     "create task with integer priority",
			priority: 5,
			value:    "test",
			want: func(t *testing.T, task Task[string], beforeTime time.Time) {
				require.Equal(t, 5, task.Priority)
				require.Equal(t, "test", task.Value)
				require.True(t, task.Time.After(beforeTime) || task.Time.Equal(beforeTime))
			},
		},
		{
			name:     "create task with zero priority",
			priority: 0,
			value:    "zero",
			want: func(t *testing.T, task Task[string], beforeTime time.Time) {
				require.Equal(t, 0, task.Priority)
				require.Equal(t, "zero", task.Value)
				require.True(t, task.Time.After(beforeTime) || task.Time.Equal(beforeTime))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			task := NewTask(tt.priority, tt.value)
			tt.want(t, task, now)
		})
	}
}

func TestPriorityCmp(t *testing.T) {
	tests := []struct {
		name     string
		task1    Task[string]
		task2    Task[string]
		expected int
	}{
		{
			name:     "higher priority number comes first",
			task1:    Task[string]{Priority: 5, Time: time.Now(), Value: "task1"},
			task2:    Task[string]{Priority: 10, Time: time.Now(), Value: "task2"},
			expected: -1,
		},
		{
			name:     "lower priority number comes second",
			task1:    Task[string]{Priority: 10, Time: time.Now(), Value: "task1"},
			task2:    Task[string]{Priority: 5, Time: time.Now(), Value: "task2"},
			expected: 1,
		},
		{
			name: "equal priorities, earlier time comes first",
			task1: Task[string]{
				Priority: 5,
				Time:     time.Now(),
				Value:    "task1",
			},
			task2: Task[string]{
				Priority: 5,
				Time:     time.Now().Add(time.Second),
				Value:    "task2",
			},
			expected: 1,
		},
		{
			name: "equal priorities, later time comes second",
			task1: Task[string]{
				Priority: 5,
				Time:     time.Now().Add(time.Second),
				Value:    "task1",
			},
			task2: Task[string]{
				Priority: 5,
				Time:     time.Now(),
				Value:    "task2",
			},
			expected: -1,
		},
		{
			name: "completely equal tasks",
			task1: Task[string]{
				Priority: 5,
				Time:     time.Now(),
				Value:    "task1",
			},
			task2: Task[string]{
				Priority: 5,
				Time:     time.Now(),
				Value:    "task2",
			},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Make sure times are the same for the equal case
			if tt.name == "completely equal tasks" {
				now := time.Now()
				tt.task1.Time = now
				tt.task2.Time = now
			}

			result := PriorityCmp(&tt.task1, &tt.task2)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestPriorityQueue_WithDifferentTypes(t *testing.T) {
	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "integer values",
			test: func(t *testing.T) {
				intPQ := NewPriorityQueue(PriorityCmp[int])

				items := []struct {
					value    int
					priority int
				}{
					{100, 3},
					{200, 5},
					{300, 1},
				}

				for _, item := range items {
					intPQ.Insert(item.value, item.priority)
				}

				expected := []struct {
					value    int
					priority int
				}{
					{200, 5},
					{100, 3},
					{300, 1},
				}

				for i, exp := range expected {
					task, err := intPQ.Pop()
					require.NoError(t, err)
					require.Equal(t, exp.value, task.Value, "Pop %d value mismatch", i)
					require.Equal(t, exp.priority, task.Priority, "Pop %d priority mismatch", i)
				}
			},
		},
		{
			name: "custom job type",
			test: func(t *testing.T) {
				type Job struct {
					ID   int
					Name string
				}

				pq := NewPriorityQueue(PriorityCmp[Job])

				jobs := []struct {
					job      Job
					priority int
				}{
					{Job{ID: 1, Name: "job1"}, 30},
					{Job{ID: 2, Name: "job2"}, 40},
					{Job{ID: 3, Name: "job3"}, 20},
				}

				for _, j := range jobs {
					pq.Insert(j.job, j.priority)
				}

				expected := []struct {
					job      Job
					priority int
				}{
					{Job{ID: 2, Name: "job2"}, 40},
					{Job{ID: 1, Name: "job1"}, 30},
					{Job{ID: 3, Name: "job3"}, 20},
				}

				for i, exp := range expected {
					task, err := pq.Pop()
					require.NoError(t, err)
					require.Equal(t, exp.job, task.Value, "Pop %d job mismatch", i)
					require.Equal(t, exp.priority, task.Priority, "Pop %d priority mismatch", i)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

func TestPriorityQueue_TimeOrdering(t *testing.T) {
	tests := []struct {
		name          string
		setupTasks    []Task[string]
		expectedOrder []string
	}{
		{
			name: "same priority different times",
			setupTasks: []Task[string]{
				{Priority: 5, Time: time.Now(), Value: "first"},
				{Priority: 5, Time: time.Now().Add(time.Millisecond), Value: "second"},
				{Priority: 5, Time: time.Now().Add(2 * time.Millisecond), Value: "third"},
			},
			expectedOrder: []string{"first", "second", "third"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewPriorityQueue(PriorityCmp[string])

			// Insert tasks manually to control timing
			for _, task := range tt.setupTasks {
				require.NoError(t, pq.heap.Insert(task))
			}

			// Verify order
			for i, expectedValue := range tt.expectedOrder {
				task, err := pq.Pop()
				require.NoError(t, err, "Pop %d should not return error", i)
				require.Equal(t, expectedValue, task.Value, "Pop %d value mismatch", i)
			}
		})
	}
}

func TestPriorityQueue_LargeDataset(t *testing.T) {
	tests := []struct {
		name       string
		priorities []int
		values     []int
		verifyFunc func(*testing.T, []int)
	}{
		{
			name:       "random priorities descending order",
			priorities: []int{50, 10, 30, 20, 40, 5, 15, 25, 35, 45},
			values:     []int{500, 100, 300, 200, 400, 50, 150, 250, 350, 450},
			verifyFunc: func(t *testing.T, poppedPriorities []int) {
				// Verify priorities are in descending order (highest first)
				for i := 1; i < len(poppedPriorities); i++ {
					require.GreaterOrEqual(t, poppedPriorities[i-1], poppedPriorities[i],
						"Priority at index %d (%d) should be >= priority at index %d (%d)",
						i-1, poppedPriorities[i-1], i, poppedPriorities[i])
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := NewPriorityQueue(PriorityCmp[int])

			// Insert items
			for i := 0; i < len(tt.priorities); i++ {
				pq.Insert(tt.values[i], tt.priorities[i])
			}

			require.Equal(t, len(tt.priorities), pq.heap.Size())

			// Pop all items and collect priorities
			var poppedPriorities []int
			for pq.heap.Size() > 0 {
				task, err := pq.Pop()
				require.NoError(t, err)
				poppedPriorities = append(poppedPriorities, task.Priority)
			}

			tt.verifyFunc(t, poppedPriorities)
		})
	}
}

func TestPriorityQueue_MultipleUpdates(t *testing.T) {
	tests := []struct {
		name          string
		setupItems    []testItem
		updates       []testUpdate
		expectedOrder []testItem
	}{
		{
			name: "multiple priority updates",
			setupItems: []testItem{
				{"A", 10},
				{"B", 20},
				{"C", 30},
			},
			updates: []testUpdate{
				{"C", 35}, // Move C to highest priority
				{"A", 5},  // Move A to lowest priority
				{"B", 25}, // Keep B in middle
			},
			expectedOrder: []testItem{
				{"C", 35},
				{"B", 25},
				{"A", 5},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pq := setupPriorityQueue(tt.setupItems)

			// Apply updates
			for _, update := range tt.updates {
				err := pq.Update(update.item, update.priority)
				require.NoError(t, err, "Update %s to priority %d should not fail", update.item, update.priority)
			}

			// Verify final order
			verifyPopOrder(t, pq, tt.expectedOrder)
		})
	}
}

// Benchmark tests
func BenchmarkPriorityQueue_Insert(b *testing.B) {
	pq := NewPriorityQueue(PriorityCmp[int])

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pq.Insert(i, i%100)
	}
}

func BenchmarkPriorityQueue_Pop(b *testing.B) {
	pq := NewPriorityQueue(PriorityCmp[int])

	// Pre-populate
	for i := 0; i < b.N; i++ {
		pq.Insert(i, i%100)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pq.Pop()
		require.NoError(b, err)
	}
}

func BenchmarkPriorityQueue_Update(b *testing.B) {
	pq := NewPriorityQueue[int](PriorityCmp[int])

	// Pre-populate with 1000 items
	for i := 0; i < 1000; i++ {
		pq.Insert(i, i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Update random items
		item := i % 1000
		newPriority := (i * 17) % 100
		err := pq.Update(item, newPriority)
		require.NoError(b, err)
	}
}

// Additional debug tests for development
func TestPriorityQueue_DebugHeapBehavior(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping debug test in short mode")
	}

	tests := []struct {
		name string
		test func(*testing.T)
	}{
		{
			name: "heap behavior validation",
			test: func(t *testing.T) {
				// Test the heap directly to understand its behavior
				heap := heap.NewHeap(PriorityCmp[string])

				task1 := Task[string]{Priority: 10, Value: "task1", Time: time.Now()}
				task2 := Task[string]{Priority: 5, Value: "task2", Time: time.Now()}
				task3 := Task[string]{Priority: 15, Value: "task3", Time: time.Now()}

				require.NoError(t, heap.Insert(task1))
				require.NoError(t, heap.Insert(task2))
				require.NoError(t, heap.Insert(task3))

				t.Logf("Initial heap:")
				for i, task := range heap.GetItems() {
					t.Logf("  [%d] %s priority %d", i, task.Value, task.Priority)
				}

				// Manually change priority and test UpHeap
				heap.GetItems()[1].Priority = 20 // Change task at index 1 to priority 20
				t.Logf("After changing priority at index 1 to 20:")
				for i, task := range heap.GetItems() {
					t.Logf("  [%d] %s priority %d", i, task.Value, task.Priority)
				}

				require.NoError(t, heap.UpHeap(1))
				t.Logf("After UpHeap(1):")
				for i, task := range heap.GetItems() {
					t.Logf("  [%d] %s priority %d", i, task.Value, task.Priority)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.test(t)
		})
	}
}

// TestPriorityQueue_ConcurrentInsertPop tests thread safety of Insert and Pop operations
// by running multiple goroutines concurrently using errgroup
func TestPriorityQueue_ConcurrentInsertPop(t *testing.T) {
	pq := NewPriorityQueue(PriorityCmp[string])

	const numGoroutines = 10
	const numOpsPerGoroutine = 50

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	// Launch goroutines that insert values
	for i := 0; i < numGoroutines; i++ {
		goroutineID := i
		g.Go(func() error {
			for j := 0; j < numOpsPerGoroutine; j++ {
				task := fmt.Sprintf("task-%d-%d", goroutineID, j)
				priority := goroutineID*numOpsPerGoroutine + j
				pq.Insert(task, priority)
			}
			return nil
		})
	}

	// Launch goroutines that pop values
	for i := 0; i < numGoroutines; i++ {
		g.Go(func() error {
			for j := 0; j < numOpsPerGoroutine; j++ {
				// Try to pop, but handle empty queue gracefully
				_, _ = pq.Pop() // Ignore errors as queue might be empty
			}
			return nil
		})
	}

	// Wait for all operations to complete
	err := g.Wait()
	require.NoError(t, err, "All goroutines should complete without error")

	// Verify no deadlocks occurred by checking priority queue is still functional
	pq.Insert("test-task", 999)
	task, err := pq.Pop()
	require.NoError(t, err, "Priority queue should still be functional after concurrent operations")
	require.NotNil(t, task, "Should be able to pop after concurrent operations")
	assert.Equal(t, "test-task", task.Value, "Inserted value should be popped")
	assert.Equal(t, 999, task.Priority, "Priority should be preserved")
}

// TestPriorityQueue_ConcurrentUpdate tests thread safety of Update operations
// by running multiple goroutines concurrently performing priority updates
func TestPriorityQueue_ConcurrentUpdate(t *testing.T) {
	pq := NewPriorityQueue(PriorityCmp[string])

	// Pre-populate the queue with test items
	testItems := []string{"task-1", "task-2", "task-3", "task-4", "task-5"}
	for i, item := range testItems {
		pq.Insert(item, i+1) // Initial priorities: 1, 2, 3, 4, 5
	}

	const numGoroutines = 10
	const numUpdatesPerGoroutine = 20

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	// Launch goroutines that perform concurrent updates
	for i := 0; i < numGoroutines; i++ {
		goroutineID := i
		g.Go(func() error {
			for j := 0; j < numUpdatesPerGoroutine; j++ {
				// Randomly update existing items
				itemIndex := (goroutineID + j) % len(testItems)
				item := testItems[itemIndex]
				newPriority := goroutineID*100 + j

				// Update may fail if item doesn't exist (due to concurrent pops), that's OK
				_ = pq.Update(item, newPriority)
			}
			return nil
		})
	}

	// Also run some concurrent inserts and pops to create realistic load
	for i := 0; i < 3; i++ {
		goroutineID := i
		g.Go(func() error {
			for j := 0; j < 10; j++ {
				// Insert new items
				newItem := fmt.Sprintf("new-task-%d-%d", goroutineID, j)
				pq.Insert(newItem, j+50)

				// Try to pop (might fail if queue is empty)
				_, _ = pq.Pop()
			}
			return nil
		})
	}

	// Wait for all operations to complete
	err := g.Wait()
	require.NoError(t, err, "All concurrent operations should complete without error")

	// Verify the priority queue is still functional
	pq.Insert("final-test", 1000)
	task, err := pq.Pop()
	require.NoError(t, err, "Priority queue should still be functional after concurrent updates")
	require.NotNil(t, task, "Should be able to pop after concurrent operations")
	assert.Equal(t, "final-test", task.Value, "Highest priority item should be popped first")
	assert.Equal(t, 1000, task.Priority, "Priority should be preserved")
}

// TestPriorityQueue_MixedConcurrentOperations tests realistic concurrent usage patterns
// combining Insert, Pop, and Update operations from multiple goroutines
func TestPriorityQueue_MixedConcurrentOperations(t *testing.T) {
	pq := NewPriorityQueue(PriorityCmp[string])

	const numWorkers = 8
	const numOperations = 100

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	// Start producer and consumer goroutines
	startProducers(g, pq, numWorkers/2, numOperations)
	startConsumers(g, pq, numWorkers/2, numOperations)

	// Wait for all operations to complete
	err := g.Wait()
	require.NoError(t, err, "All mixed concurrent operations should complete without error")

	// Verify final state functionality
	verifyPriorityQueueFunctionality(t, pq)
}

// startProducers launches producer goroutines that insert and update tasks
func startProducers(g *errgroup.Group, pq *PriorityQueue[string], numProducers, numOperations int) {
	for i := 0; i < numProducers; i++ {
		workerID := i
		g.Go(func() error {
			return runProducerWorker(pq, workerID, numOperations)
		})
	}
}

// runProducerWorker executes producer operations for a single worker
func runProducerWorker(pq *PriorityQueue[string], workerID, numOperations int) error {
	for j := 0; j < numOperations; j++ {
		task := fmt.Sprintf("producer-%d-task-%d", workerID, j)
		priority := workerID*1000 + j
		pq.Insert(task, priority)

		// Occasionally update an existing task's priority
		if j%10 == 0 && j > 0 {
			oldTask := fmt.Sprintf("producer-%d-task-%d", workerID, j-5)
			newPriority := priority + 5000 // Higher priority
			_ = pq.Update(oldTask, newPriority)
		}
	}
	return nil
}

// startConsumers launches consumer goroutines that pop and process tasks
func startConsumers(g *errgroup.Group, pq *PriorityQueue[string], numConsumers, numOperations int) {
	for i := 0; i < numConsumers; i++ {
		g.Go(func() error {
			return runConsumerWorker(pq, numOperations)
		})
	}
}

// runConsumerWorker executes consumer operations for a single worker
func runConsumerWorker(pq *PriorityQueue[string], numOperations int) error {
	processedCount := 0
	for processedCount < numOperations {
		if task, err := pq.Pop(); err == nil {
			// Simulate processing time variation
			if processedCount%20 == 0 {
				// Occasionally re-insert a task with different priority
				newTask := fmt.Sprintf("reprocessed-%s", task.Value)
				newPriority := task.Priority - 100 // Lower priority for reprocessing
				pq.Insert(newTask, newPriority)
			}
			processedCount++
		}
		// If queue is empty, continue trying (producers might still be working)
	}
	return nil
}

// verifyPriorityQueueFunctionality tests that the priority queue maintains correct behavior
func verifyPriorityQueueFunctionality(t *testing.T, pq *PriorityQueue[string]) {
	// Insert some test items with known priorities
	testPriorities := []int{100, 500, 50, 750, 25}
	for i, priority := range testPriorities {
		pq.Insert(fmt.Sprintf("final-test-%d", i), priority)
	}

	// Pop items and verify they come out in priority order (highest first)
	poppedPriorities := collectTestPriorities(pq, len(testPriorities))

	// Verify priorities are in descending order (max-heap behavior)
	require.Equal(t, len(testPriorities), len(poppedPriorities), "Should pop all test items")
	verifyDescendingOrder(t, poppedPriorities)
}

// collectTestPriorities pops test items and returns their priorities
func collectTestPriorities(pq *PriorityQueue[string], expectedCount int) []int {
	var poppedPriorities []int
	for len(poppedPriorities) < expectedCount {
		if task, popErr := pq.Pop(); popErr == nil {
			if len(task.Value) >= 10 && task.Value[:10] == "final-test" { // Only check our test items
				poppedPriorities = append(poppedPriorities, task.Priority)
			}
		} else {
			break // No more items
		}
	}
	return poppedPriorities
}

// verifyDescendingOrder checks that priorities are in descending order
func verifyDescendingOrder(t *testing.T, priorities []int) {
	for i := 1; i < len(priorities); i++ {
		assert.GreaterOrEqual(t, priorities[i-1], priorities[i],
			"Priorities should be in descending order (index %d: %d >= %d)",
			i, priorities[i-1], priorities[i])
	}
}
