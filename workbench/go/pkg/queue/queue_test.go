package queue

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewQueue(t *testing.T) {
	t.Run("valid size", func(t *testing.T) {
		q := NewQueue[int](5)
		assert.NotNil(t, q)
		assert.True(t, q.IsEmpty())
		assert.False(t, q.IsFull())
	})

	t.Run("panic on invalid size", func(t *testing.T) {
		assert.Panics(t, func() {
			NewQueue[int](0)
		})
		assert.Panics(t, func() {
			NewQueue[int](-1)
		})
	})
}

func TestIsEmpty(t *testing.T) {
	q := NewQueue[int](3)
	assert.True(t, q.IsEmpty())

	item := 1
	require.NoError(t, q.Enqueue(&item))
	assert.False(t, q.IsEmpty())

	_, err := q.Dequeue()
	require.NoError(t, err)
	assert.True(t, q.IsEmpty())
}

func TestIsFull(t *testing.T) {
	q := NewQueue[int](3)
	assert.False(t, q.IsFull())

	item1 := 1
	require.NoError(t, q.Enqueue(&item1))
	assert.False(t, q.IsFull())

	item2 := 2
	require.NoError(t, q.Enqueue(&item2))
	assert.False(t, q.IsFull())

	item3 := 3
	require.NoError(t, q.Enqueue(&item3))
	assert.True(t, q.IsFull())
}

func TestEnqueue(t *testing.T) {
	t.Run("successful enqueue", func(t *testing.T) {
		q := NewQueue[int](3)

		item1 := 10
		err := q.Enqueue(&item1)
		assert.NoError(t, err)
		assert.False(t, q.IsEmpty())

		item2 := 20
		err = q.Enqueue(&item2)
		assert.NoError(t, err)

		item3 := 30
		err = q.Enqueue(&item3)
		assert.NoError(t, err)
		assert.True(t, q.IsFull())
	})

	t.Run("enqueue to full queue", func(t *testing.T) {
		q := NewQueue[int](2)

		item1 := 1
		require.NoError(t, q.Enqueue(&item1))
		item2 := 2
		require.NoError(t, q.Enqueue(&item2))
		assert.True(t, q.IsFull())

		item3 := 3
		err := q.Enqueue(&item3)
		assert.Equal(t, ErrorQueueOverflow, err)
	})
}

func TestDequeue(t *testing.T) {
	t.Run("successful dequeue", func(t *testing.T) {
		q := NewQueue[int](3)

		item1 := 10
		item2 := 20
		item3 := 30

		require.NoError(t, q.Enqueue(&item1))
		require.NoError(t, q.Enqueue(&item2))
		require.NoError(t, q.Enqueue(&item3))

		// Dequeue in FIFO order
		dequeued, err := q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, 10, *dequeued)

		dequeued, err = q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, 20, *dequeued)

		dequeued, err = q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, 30, *dequeued)

		assert.True(t, q.IsEmpty())
	})

	t.Run("dequeue from empty queue", func(t *testing.T) {
		q := NewQueue[int](3)

		_, err := q.Dequeue()
		assert.Equal(t, ErrorQueueUnderflow, err)
	})
}

func TestPeek(t *testing.T) {
	t.Run("successful peek", func(t *testing.T) {
		q := NewQueue[int](3)

		item1 := 10
		item2 := 20

		require.NoError(t, q.Enqueue(&item1))
		require.NoError(t, q.Enqueue(&item2))

		// Peek should return the first item without removing it
		peeked, err := q.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 10, *peeked)

		// Queue should remain unchanged
		assert.False(t, q.IsEmpty())

		// Peek again should return the same item
		peeked, err = q.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 10, *peeked)

		// Dequeue should return the same item that was peeked
		dequeued, err := q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, 10, *dequeued)
	})

	t.Run("peek empty queue", func(t *testing.T) {
		q := NewQueue[int](3)

		_, err := q.Peek()
		assert.Equal(t, ErrorQueueUnderflow, err)
	})
}

func TestCircularBehavior(t *testing.T) {
	q := NewQueue[int](3)

	// Fill the queue
	item1 := 1
	item2 := 2
	item3 := 3
	require.NoError(t, q.Enqueue(&item1))
	require.NoError(t, q.Enqueue(&item2))
	require.NoError(t, q.Enqueue(&item3))
	assert.True(t, q.IsFull())

	// Dequeue one item
	dequeued, err := q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 1, *dequeued)
	assert.False(t, q.IsFull())

	// Enqueue another item (should wrap around)
	item4 := 4
	require.NoError(t, q.Enqueue(&item4))

	// Dequeue remaining items to verify order
	dequeued, err = q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 2, *dequeued)

	dequeued, err = q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 3, *dequeued)

	dequeued, err = q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 4, *dequeued)

	assert.True(t, q.IsEmpty())
}

func TestGenericTypes(t *testing.T) {
	t.Run("string queue", func(t *testing.T) {
		q := NewQueue[string](2)

		str1 := "hello"
		str2 := "world"

		require.NoError(t, q.Enqueue(&str1))
		require.NoError(t, q.Enqueue(&str2))

		dequeued, err := q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, "hello", *dequeued)

		dequeued, err = q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, "world", *dequeued)
	})

	t.Run("struct queue", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		q := NewQueue[Person](2)

		person1 := Person{Name: "Alice", Age: 30}
		person2 := Person{Name: "Bob", Age: 25}

		require.NoError(t, q.Enqueue(&person1))
		require.NoError(t, q.Enqueue(&person2))

		dequeued, err := q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, "Alice", dequeued.Name)
		assert.Equal(t, 30, dequeued.Age)

		dequeued, err = q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, "Bob", dequeued.Name)
		assert.Equal(t, 25, dequeued.Age)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("single element queue", func(t *testing.T) {
		q := NewQueue[int](1)

		item := 42
		require.NoError(t, q.Enqueue(&item))
		assert.True(t, q.IsFull())

		// Try to enqueue another item
		item2 := 43
		err := q.Enqueue(&item2)
		assert.Equal(t, ErrorQueueOverflow, err)

		// Dequeue the item
		dequeued, err := q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, 42, *dequeued)
		assert.True(t, q.IsEmpty())
	})

	t.Run("multiple enqueue/dequeue cycles", func(t *testing.T) {
		q := NewQueue[int](3)

		// First cycle
		for i := 1; i <= 3; i++ {
			item := i
			require.NoError(t, q.Enqueue(&item))
		}

		for i := 1; i <= 3; i++ {
			dequeued, err := q.Dequeue()
			assert.NoError(t, err)
			assert.Equal(t, i, *dequeued)
		}
		assert.True(t, q.IsEmpty())

		// Second cycle
		for i := 4; i <= 6; i++ {
			item := i
			require.NoError(t, q.Enqueue(&item))
		}

		for i := 4; i <= 6; i++ {
			dequeued, err := q.Dequeue()
			assert.NoError(t, err)
			assert.Equal(t, i, *dequeued)
		}
		assert.True(t, q.IsEmpty())
	})
}
