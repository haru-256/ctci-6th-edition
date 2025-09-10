package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQueue(t *testing.T) {
	tests := []struct {
		name string
		size int
		test func(t *testing.T, q *Queue)
	}{
		{
			name: "creates empty queue with specified size",
			size: 5,
			test: func(t *testing.T, q *Queue) {
				require.NotNil(t, q)
				assert.True(t, q.IsEmpty())
				assert.False(t, q.IsFull())
				assert.Equal(t, 0, q.head)
				assert.Equal(t, 0, q.tail)
				assert.Equal(t, 5, q.size)
				assert.Equal(t, 5, len(q.items))
			},
		},
		{
			name: "creates queue with size 1",
			size: 1,
			test: func(t *testing.T, q *Queue) {
				require.NotNil(t, q)
				assert.True(t, q.IsEmpty())
				assert.True(t, q.IsFull()) // Size 1 queue is immediately full due to circular buffer
				assert.Equal(t, 1, q.size)
			},
		},
		{
			name: "creates large queue",
			size: 100,
			test: func(t *testing.T, q *Queue) {
				require.NotNil(t, q)
				assert.True(t, q.IsEmpty())
				assert.False(t, q.IsFull())
				assert.Equal(t, 100, q.size)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewQueue(tt.size)
			tt.test(t, queue)
		})
	}
}

func TestQueue_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Queue
		expected bool
	}{
		{
			name: "empty queue returns true",
			setup: func() *Queue {
				return NewQueue(5)
			},
			expected: true,
		},
		{
			name: "non-empty queue returns false",
			setup: func() *Queue {
				q := NewQueue(5)
				_ = q.Enqueue("item")
				return q
			},
			expected: false,
		},
		{
			name: "queue after enqueue and dequeue returns true",
			setup: func() *Queue {
				q := NewQueue(5)
				_ = q.Enqueue("item")
				_, _ = q.Dequeue()
				return q
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := tt.setup()
			result := queue.IsEmpty()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQueue_IsFull(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Queue
		expected bool
	}{
		{
			name: "empty queue returns false",
			setup: func() *Queue {
				return NewQueue(3)
			},
			expected: false,
		},
		{
			name: "partially filled queue returns false",
			setup: func() *Queue {
				q := NewQueue(3)
				_ = q.Enqueue("item1")
				return q
			},
			expected: false,
		},
		{
			name: "full queue returns true",
			setup: func() *Queue {
				q := NewQueue(3)
				_ = q.Enqueue("item1")
				_ = q.Enqueue("item2")
				return q
			},
			expected: true,
		},
		{
			name: "size 2 queue with one item is full",
			setup: func() *Queue {
				q := NewQueue(2)
				_ = q.Enqueue("item")
				return q
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := tt.setup()
			result := queue.IsFull()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestQueue_Enqueue(t *testing.T) {
	tests := []struct {
		name      string
		queueSize int
		items     []any
		checkFunc func(t *testing.T, q *Queue)
	}{
		{
			name:      "enqueue single item",
			queueSize: 5,
			items:     []any{"hello"},
			checkFunc: func(t *testing.T, q *Queue) {
				assert.False(t, q.IsEmpty())
				assert.False(t, q.IsFull())
				assert.Equal(t, 0, q.head)
				assert.Equal(t, 1, q.tail)
				assert.Equal(t, "hello", q.items[0])
			},
		},
		{
			name:      "enqueue multiple items",
			queueSize: 5,
			items:     []any{"first", "second", "third"},
			checkFunc: func(t *testing.T, q *Queue) {
				assert.False(t, q.IsEmpty())
				assert.False(t, q.IsFull())
				assert.Equal(t, 0, q.head)
				assert.Equal(t, 3, q.tail)
				assert.Equal(t, "first", q.items[0])
				assert.Equal(t, "second", q.items[1])
				assert.Equal(t, "third", q.items[2])
			},
		},
		{
			name:      "enqueue different types",
			queueSize: 5,
			items:     []any{42, "string", 3.14, true},
			checkFunc: func(t *testing.T, q *Queue) {
				assert.False(t, q.IsEmpty())
				assert.True(t, q.IsFull()) // 4 items in size 5 queue = full (circular buffer)
				assert.Equal(t, 4, q.tail)
				assert.Equal(t, 42, q.items[0])
				assert.Equal(t, "string", q.items[1])
				assert.Equal(t, 3.14, q.items[2])
				assert.Equal(t, true, q.items[3])
			},
		},
		{
			name:      "enqueue nil item",
			queueSize: 5,
			items:     []any{nil},
			checkFunc: func(t *testing.T, q *Queue) {
				assert.False(t, q.IsEmpty())
				assert.Equal(t, 1, q.tail)
				assert.Nil(t, q.items[0])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewQueue(tt.queueSize)

			for _, item := range tt.items {
				err := queue.Enqueue(item)
				require.NoError(t, err)
			}

			tt.checkFunc(t, queue)
		})
	}
}

func TestQueue_EnqueueOverflow(t *testing.T) {
	tests := []struct {
		name         string
		queueSize    int
		enqueueCount int
	}{
		{
			name:         "enqueue to exactly full queue",
			queueSize:    3,
			enqueueCount: 3,
		},
		{
			name:         "enqueue beyond queue capacity",
			queueSize:    2,
			enqueueCount: 5,
		},
		{
			name:         "enqueue to size 2 queue",
			queueSize:    2,
			enqueueCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewQueue(tt.queueSize)

			var err error
			for i := 0; i < tt.enqueueCount; i++ {
				err = queue.Enqueue(i)
				if i < tt.queueSize-1 { // -1 because circular queue leaves one spot empty
					require.NoError(t, err, "Enqueue %d should succeed", i)
				} else {
					require.Error(t, err, "Enqueue %d should fail due to overflow", i)
					assert.ErrorIs(t, err, ErrorQueueOverflow)
					assert.True(t, queue.IsFull(), "Queue should be full")
				}
			}
		})
	}
}

func TestQueue_Dequeue(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *Queue
		expectedItem  any
		expectedError error
		checkQueue    func(t *testing.T, q *Queue)
	}{
		{
			name: "dequeue from empty queue returns error",
			setup: func() *Queue {
				return NewQueue(5)
			},
			expectedItem:  nil,
			expectedError: ErrorQueueUnderflow,
			checkQueue: func(t *testing.T, q *Queue) {
				assert.True(t, q.IsEmpty())
			},
		},
		{
			name: "dequeue single item from queue",
			setup: func() *Queue {
				q := NewQueue(5)
				_ = q.Enqueue("item")
				return q
			},
			expectedItem:  "item",
			expectedError: nil,
			checkQueue: func(t *testing.T, q *Queue) {
				assert.True(t, q.IsEmpty())
				assert.Equal(t, 1, q.head)
				assert.Equal(t, 1, q.tail)
			},
		},
		{
			name: "dequeue from queue with multiple items returns first enqueued",
			setup: func() *Queue {
				q := NewQueue(5)
				_ = q.Enqueue("first")
				_ = q.Enqueue("second")
				_ = q.Enqueue("third")
				return q
			},
			expectedItem:  "first",
			expectedError: nil,
			checkQueue: func(t *testing.T, q *Queue) {
				assert.False(t, q.IsEmpty())
				assert.Equal(t, 1, q.head)
				assert.Equal(t, 3, q.tail)
				assert.Nil(t, q.items[0]) // Should be cleared
				assert.Equal(t, "second", q.items[1])
				assert.Equal(t, "third", q.items[2])
			},
		},
		{
			name: "dequeue nil item",
			setup: func() *Queue {
				q := NewQueue(5)
				_ = q.Enqueue(nil)
				return q
			},
			expectedItem:  nil,
			expectedError: nil,
			checkQueue: func(t *testing.T, q *Queue) {
				assert.True(t, q.IsEmpty())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := tt.setup()

			item, err := queue.Dequeue()

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedItem, item)
			tt.checkQueue(t, queue)
		})
	}
}

func TestQueue_Peek(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *Queue
		expectedItem  any
		expectedError error
		checkQueue    func(t *testing.T, q *Queue) // Peek should not modify the queue
	}{
		{
			name: "peek empty queue returns error",
			setup: func() *Queue {
				return NewQueue(5)
			},
			expectedItem:  nil,
			expectedError: ErrorQueueUnderflow,
			checkQueue: func(t *testing.T, q *Queue) {
				assert.True(t, q.IsEmpty())
				assert.Equal(t, 0, q.head)
				assert.Equal(t, 0, q.tail)
			},
		},
		{
			name: "peek single item queue",
			setup: func() *Queue {
				q := NewQueue(5)
				_ = q.Enqueue("item")
				return q
			},
			expectedItem:  "item",
			expectedError: nil,
			checkQueue: func(t *testing.T, q *Queue) {
				assert.False(t, q.IsEmpty())
				assert.Equal(t, 0, q.head)
				assert.Equal(t, 1, q.tail)
				assert.Equal(t, "item", q.items[0])
			},
		},
		{
			name: "peek queue with multiple items returns first enqueued",
			setup: func() *Queue {
				q := NewQueue(5)
				_ = q.Enqueue("first")
				_ = q.Enqueue("second")
				_ = q.Enqueue("third")
				return q
			},
			expectedItem:  "first",
			expectedError: nil,
			checkQueue: func(t *testing.T, q *Queue) {
				assert.False(t, q.IsEmpty())
				assert.Equal(t, 0, q.head)
				assert.Equal(t, 3, q.tail)
				assert.Equal(t, "first", q.items[0])
				assert.Equal(t, "second", q.items[1])
				assert.Equal(t, "third", q.items[2])
			},
		},
		{
			name: "peek nil item",
			setup: func() *Queue {
				q := NewQueue(5)
				_ = q.Enqueue(nil)
				return q
			},
			expectedItem:  nil,
			expectedError: nil,
			checkQueue: func(t *testing.T, q *Queue) {
				assert.False(t, q.IsEmpty())
				assert.Equal(t, 0, q.head)
				assert.Equal(t, 1, q.tail)
				assert.Nil(t, q.items[0])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := tt.setup()
			originalHead := queue.head
			originalTail := queue.tail

			item, err := queue.Peek()

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedItem, item)

			// Ensure peek doesn't modify the queue
			assert.Equal(t, originalHead, queue.head)
			assert.Equal(t, originalTail, queue.tail)
			tt.checkQueue(t, queue)
		})
	}
}

func TestQueue_CircularBehavior(t *testing.T) {
	tests := []struct {
		name       string
		operations func(t *testing.T, q *Queue)
	}{
		{
			name: "circular enqueue and dequeue",
			operations: func(t *testing.T, q *Queue) {
				// Fill the queue
				err := q.Enqueue("a")
				require.NoError(t, err)
				err = q.Enqueue("b")
				require.NoError(t, err)

				assert.True(t, q.IsFull())

				// Dequeue one item
				item, err := q.Dequeue()
				require.NoError(t, err)
				assert.Equal(t, "a", item)
				assert.False(t, q.IsFull())

				// Enqueue another item (should wrap around)
				err = q.Enqueue("c")
				require.NoError(t, err)
				assert.True(t, q.IsFull())

				// Verify order: b, c
				item, err = q.Dequeue()
				require.NoError(t, err)
				assert.Equal(t, "b", item)

				item, err = q.Dequeue()
				require.NoError(t, err)
				assert.Equal(t, "c", item)

				assert.True(t, q.IsEmpty())
			},
		},
		{
			name: "multiple wrap arounds",
			operations: func(t *testing.T, q *Queue) {
				items := []string{"a", "b", "c", "d", "e", "f", "g", "h"}

				for i, item := range items {
					if i < 2 {
						// Fill initially
						err := q.Enqueue(item)
						require.NoError(t, err)
					} else {
						// For subsequent items, dequeue one and enqueue one
						dequeued, err := q.Dequeue()
						require.NoError(t, err)
						assert.Equal(t, items[i-2], dequeued)

						err = q.Enqueue(item)
						require.NoError(t, err)
					}
				}

				// Dequeue remaining items
				for i := len(items) - 2; i < len(items); i++ {
					item, err := q.Dequeue()
					require.NoError(t, err)
					assert.Equal(t, items[i], item)
				}

				assert.True(t, q.IsEmpty())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewQueue(3) // Small queue to test circular behavior
			tt.operations(t, queue)
		})
	}
}

func TestQueue_SequentialOperations(t *testing.T) {
	tests := []struct {
		name       string
		operations func(t *testing.T, q *Queue)
	}{
		{
			name: "enqueue-dequeue-enqueue sequence",
			operations: func(t *testing.T, q *Queue) {
				// Enqueue first item
				err := q.Enqueue("first")
				require.NoError(t, err)
				assert.False(t, q.IsEmpty())

				// Peek should return first item
				item, err := q.Peek()
				require.NoError(t, err)
				assert.Equal(t, "first", item)

				// Dequeue should return first item
				item, err = q.Dequeue()
				require.NoError(t, err)
				assert.Equal(t, "first", item)
				assert.True(t, q.IsEmpty())

				// Enqueue second item
				err = q.Enqueue("second")
				require.NoError(t, err)
				assert.False(t, q.IsEmpty())

				// Peek should return second item
				item, err = q.Peek()
				require.NoError(t, err)
				assert.Equal(t, "second", item)
			},
		},
		{
			name: "multiple enqueue and dequeue operations",
			operations: func(t *testing.T, q *Queue) {
				items := []string{"a", "b", "c", "d", "e"}

				// Enqueue all items
				for _, item := range items {
					err := q.Enqueue(item)
					require.NoError(t, err)
				}

				// Queue should not be empty
				assert.False(t, q.IsEmpty())

				// Dequeue all items in FIFO order
				for _, expectedItem := range items {
					item, err := q.Dequeue()
					require.NoError(t, err)
					assert.Equal(t, expectedItem, item)
				}

				// Queue should be empty
				assert.True(t, q.IsEmpty())
			},
		},
		{
			name: "alternating enqueue and peek operations",
			operations: func(t *testing.T, q *Queue) {
				// Enqueue and peek multiple times - should always see first item
				err := q.Enqueue(1)
				require.NoError(t, err)
				item, err := q.Peek()
				require.NoError(t, err)
				assert.Equal(t, 1, item)

				err = q.Enqueue(2)
				require.NoError(t, err)
				item, err = q.Peek()
				require.NoError(t, err)
				assert.Equal(t, 1, item) // Still first item

				err = q.Enqueue(3)
				require.NoError(t, err)
				item, err = q.Peek()
				require.NoError(t, err)
				assert.Equal(t, 1, item) // Still first item

				// Queue should have 3 items but peek always returns first
				assert.False(t, q.IsEmpty())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewQueue(10)
			tt.operations(t, queue)
		})
	}
}

func TestQueue_ErrorConditions(t *testing.T) {
	tests := []struct {
		name      string
		operation func(t *testing.T, q *Queue)
	}{
		{
			name: "multiple dequeues from empty queue",
			operation: func(t *testing.T, q *Queue) {
				// First dequeue should return error
				item, err := q.Dequeue()
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorQueueUnderflow)
				assert.Nil(t, item)

				// Second dequeue should also return error
				item, err = q.Dequeue()
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorQueueUnderflow)
				assert.Nil(t, item)

				assert.True(t, q.IsEmpty())
			},
		},
		{
			name: "multiple peeks from empty queue",
			operation: func(t *testing.T, q *Queue) {
				// First peek should return error
				item, err := q.Peek()
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorQueueUnderflow)
				assert.Nil(t, item)

				// Second peek should also return error
				item, err = q.Peek()
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorQueueUnderflow)
				assert.Nil(t, item)

				assert.True(t, q.IsEmpty())
			},
		},
		{
			name: "dequeue until empty then continue dequeuing",
			operation: func(t *testing.T, q *Queue) {
				// Add some items
				err := q.Enqueue("a")
				require.NoError(t, err)
				err = q.Enqueue("b")
				require.NoError(t, err)

				// Dequeue all items
				_, err = q.Dequeue()
				require.NoError(t, err)
				_, err = q.Dequeue()
				require.NoError(t, err)

				assert.True(t, q.IsEmpty())

				// Try to dequeue from empty queue
				item, err := q.Dequeue()
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorQueueUnderflow)
				assert.Nil(t, item)
			},
		},
		{
			name: "enqueue until full then continue enqueuing",
			operation: func(t *testing.T, q *Queue) {
				// Fill the queue (size 3, but can only hold 2 items due to circular implementation)
				err := q.Enqueue("a")
				require.NoError(t, err)
				err = q.Enqueue("b")
				require.NoError(t, err)

				assert.True(t, q.IsFull())

				// Try to enqueue to full queue
				err = q.Enqueue("c")
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorQueueOverflow)

				// Try again
				err = q.Enqueue("d")
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorQueueOverflow)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			queue := NewQueue(3)
			tt.operation(t, queue)
		})
	}
}
