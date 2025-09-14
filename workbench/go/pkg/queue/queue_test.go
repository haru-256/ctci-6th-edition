package queue

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
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
	require.NoError(t, q.Enqueue(item))
	assert.False(t, q.IsEmpty())

	_, err := q.Dequeue()
	require.NoError(t, err)
	assert.True(t, q.IsEmpty())
}

func TestIsFull(t *testing.T) {
	q := NewQueue[int](3)
	assert.False(t, q.IsFull())

	item1 := 1
	require.NoError(t, q.Enqueue(item1))
	assert.False(t, q.IsFull())

	item2 := 2
	require.NoError(t, q.Enqueue(item2))
	assert.False(t, q.IsFull())

	item3 := 3
	require.NoError(t, q.Enqueue(item3))
	assert.True(t, q.IsFull())
}

func TestEnqueue(t *testing.T) {
	t.Run("successful enqueue", func(t *testing.T) {
		q := NewQueue[int](3)

		item1 := 10
		err := q.Enqueue(item1)
		assert.NoError(t, err)
		assert.False(t, q.IsEmpty())

		item2 := 20
		err = q.Enqueue(item2)
		assert.NoError(t, err)

		item3 := 30
		err = q.Enqueue(item3)
		assert.NoError(t, err)
		assert.True(t, q.IsFull())
	})

	t.Run("enqueue to full queue", func(t *testing.T) {
		q := NewQueue[int](2)

		item1 := 1
		require.NoError(t, q.Enqueue(item1))
		item2 := 2
		require.NoError(t, q.Enqueue(item2))
		assert.True(t, q.IsFull())

		item3 := 3
		err := q.Enqueue(item3)
		assert.Equal(t, ErrorQueueOverflow, err)
	})
}

func TestDequeue(t *testing.T) {
	t.Run("successful dequeue", func(t *testing.T) {
		q := NewQueue[int](3)

		item1 := 10
		item2 := 20
		item3 := 30

		require.NoError(t, q.Enqueue(item1))
		require.NoError(t, q.Enqueue(item2))
		require.NoError(t, q.Enqueue(item3))

		// Dequeue in FIFO order
		dequeued, err := q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, 10, dequeued)

		dequeued, err = q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, 20, dequeued)

		dequeued, err = q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, 30, dequeued)

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

		require.NoError(t, q.Enqueue(item1))
		require.NoError(t, q.Enqueue(item2))

		// Peek should return the first item without removing it
		peeked, err := q.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 10, peeked)

		// Queue should remain unchanged
		assert.False(t, q.IsEmpty())

		// Peek again should return the same item
		peeked, err = q.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 10, peeked)

		// Dequeue should return the same item that was peeked
		dequeued, err := q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, 10, dequeued)
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
	require.NoError(t, q.Enqueue(item1))
	require.NoError(t, q.Enqueue(item2))
	require.NoError(t, q.Enqueue(item3))
	assert.True(t, q.IsFull())

	// Dequeue one item
	dequeued, err := q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 1, dequeued)
	assert.False(t, q.IsFull())

	// Enqueue another item (should wrap around)
	item4 := 4
	require.NoError(t, q.Enqueue(item4))

	// Dequeue remaining items to verify order
	dequeued, err = q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 2, dequeued)

	dequeued, err = q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 3, dequeued)

	dequeued, err = q.Dequeue()
	assert.NoError(t, err)
	assert.Equal(t, 4, dequeued)

	assert.True(t, q.IsEmpty())
}

func FuzzQueue_EnqueueDequeue(f *testing.F) {
	f.Add(1)
	f.Add(42)
	f.Add(-7)
	f.Fuzz(func(t *testing.T, x int) {
		q := NewQueue[int](10)
		require.NoError(t, q.Enqueue(x))
		got, err := q.Dequeue()
		require.NoError(t, err)
		assert.Equal(t, x, got)
		assert.True(t, q.IsEmpty())
	})
}

func TestQueue_ConcurrentEnqueueDequeue(t *testing.T) {
	q := NewQueue[int](100)
	var g errgroup.Group

	g.Go(func() error {
		for i := 0; i < 100; i++ {
			if err := q.Enqueue(i); err != nil {
				return err
			}
		}
		return nil
	})
	g.Go(func() error {
		for i := 0; i < 100; i++ {
			_, err := q.Dequeue()
			if err != nil && err != ErrorQueueUnderflow {
				return err
			}
		}
		return nil
	})
	require.NoError(t, g.Wait(), "errgroup should not return error")
	// Not strictly thread-safe, but should not panic or deadlock
}

func TestQueue_ConcurrentStress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	q := NewQueue[int](1000)
	var g errgroup.Group
	const numGoroutines = 10
	const itemsPerGoroutine = 100

	// Multiple producers
	for i := 0; i < numGoroutines; i++ {
		i := i
		g.Go(func() error {
			for j := 0; j < itemsPerGoroutine; j++ {
				if err := q.Enqueue(i*itemsPerGoroutine + j); err != nil {
					return err
				}
			}
			return nil
		})
	}

	// Multiple consumers
	for i := 0; i < numGoroutines; i++ {
		g.Go(func() error {
			for j := 0; j < itemsPerGoroutine; j++ {
				_, err := q.Dequeue()
				if err != nil && err != ErrorQueueUnderflow {
					return err
				}
			}
			return nil
		})
	}

	require.NoError(t, g.Wait(), "stress test should not fail")
}

func TestQueue_ZeroAndPointerValues(t *testing.T) {
	q := NewQueue[*int](3)
	var nilPtr *int
	require.NoError(t, q.Enqueue(nilPtr))
	v, err := q.Dequeue()
	require.NoError(t, err)
	assert.Nil(t, v)

	q2 := NewQueue[int](2)
	require.NoError(t, q2.Enqueue(0))
	v2, err := q2.Dequeue()
	require.NoError(t, err)
	assert.Equal(t, 0, v2)
}

func TestGenericTypes(t *testing.T) {
	t.Run("string queue", func(t *testing.T) {
		q := NewQueue[string](2)

		str1 := "hello"
		str2 := "world"

		require.NoError(t, q.Enqueue(str1))
		require.NoError(t, q.Enqueue(str2))

		dequeued, err := q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, "hello", dequeued)

		dequeued, err = q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, "world", dequeued)
	})

	t.Run("struct queue", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		q := NewQueue[Person](2)

		person1 := Person{Name: "Alice", Age: 30}
		person2 := Person{Name: "Bob", Age: 25}

		require.NoError(t, q.Enqueue(person1))
		require.NoError(t, q.Enqueue(person2))

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
		require.NoError(t, q.Enqueue(item))
		assert.True(t, q.IsFull())

		// Try to enqueue another item
		item2 := 43
		err := q.Enqueue(item2)
		assert.Equal(t, ErrorQueueOverflow, err)

		// Dequeue the item
		dequeued, err := q.Dequeue()
		assert.NoError(t, err)
		assert.Equal(t, 42, dequeued)
		assert.True(t, q.IsEmpty())
	})

	t.Run("multiple enqueue/dequeue cycles", func(t *testing.T) {
		q := NewQueue[int](3)

		// First cycle
		for i := 1; i <= 3; i++ {
			item := i
			require.NoError(t, q.Enqueue(item))
		}

		for i := 1; i <= 3; i++ {
			dequeued, err := q.Dequeue()
			assert.NoError(t, err)
			assert.Equal(t, i, dequeued)
		}
		assert.True(t, q.IsEmpty())

		// Second cycle
		for i := 4; i <= 6; i++ {
			item := i
			require.NoError(t, q.Enqueue(item))
		}

		for i := 4; i <= 6; i++ {
			dequeued, err := q.Dequeue()
			assert.NoError(t, err)
			assert.Equal(t, i, dequeued)
		}
		assert.True(t, q.IsEmpty())
	})
}

// Example functions for godoc

func ExampleNewQueue() {
	// Create a queue of integers with capacity 5
	q := NewQueue[int](5)
	fmt.Println("Created queue with capacity:", q.Size())
	fmt.Println("Queue is empty:", q.IsEmpty())
}

func ExampleQueue_Enqueue() {
	q := NewQueue[string](3)

	err := q.Enqueue("first")
	fmt.Println("Enqueue error:", err)

	err = q.Enqueue("second")
	fmt.Println("Enqueue error:", err)

	fmt.Println("Queue count:", q.Count())
	fmt.Println("Queue is full:", q.IsFull())

	// Output:
	// Enqueue error: <nil>
	// Enqueue error: <nil>
	// Queue count: 2
	// Queue is full: false
}

func ExampleQueue_Dequeue() {
	q := NewQueue[int](3)
	err := q.Enqueue(10)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = q.Enqueue(20)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = q.Enqueue(30)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Dequeue items in FIFO order
	item, err := q.Dequeue()
	fmt.Println("Dequeued:", item, "Error:", err)

	item, err = q.Dequeue()
	fmt.Println("Dequeued:", item, "Error:", err)

	fmt.Println("Queue count:", q.Count())

	// Output:
	// Dequeued: 10 Error: <nil>
	// Dequeued: 20 Error: <nil>
	// Queue count: 1
}

func ExampleQueue_Peek() {
	q := NewQueue[string](3)
	err := q.Enqueue("hello")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	err = q.Enqueue("world")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Peek at the front item without removing it
	item, err := q.Peek()
	fmt.Println("Peeked:", item, "Error:", err)
	fmt.Println("Queue count after peek:", q.Count())

	// Dequeue the front item
	item, err = q.Dequeue()
	fmt.Println("Dequeued:", item, "Error:", err)
	fmt.Println("Queue count after dequeue:", q.Count())

	// Output:
	// Peeked: hello Error: <nil>
	// Queue count after peek: 2
	// Dequeued: hello Error: <nil>
	// Queue count after dequeue: 1
}

func ExampleQueue_circularBehavior() {
	q := NewQueue[int](3)

	// Fill the queue
	if err := q.Enqueue(1); err != nil {
		fmt.Println("Error:", err)
	}
	if err := q.Enqueue(2); err != nil {
		fmt.Println("Error:", err)
	}
	if err := q.Enqueue(3); err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Queue is full:", q.IsFull())

	// Remove one item
	item, _ := q.Dequeue()
	fmt.Println("Dequeued:", item)

	// Add another item (demonstrates circular behavior)
	if err := q.Enqueue(4); err != nil {
		fmt.Println("Error:", err)
	}
	fmt.Println("Queue count:", q.Count())

	// Dequeue remaining items
	for !q.IsEmpty() {
		item, _ = q.Dequeue()
		fmt.Println("Dequeued:", item)
	}

	// Output:
	// Queue is full: true
	// Dequeued: 1
	// Queue count: 3
	// Dequeued: 2
	// Dequeued: 3
	// Dequeued: 4
}

// Benchmark functions

func BenchmarkQueueEnqueue(b *testing.B) {
	q := NewQueue[int](b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if err := q.Enqueue(i); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQueueDequeue(b *testing.B) {
	q := NewQueue[int](b.N)
	// Pre-fill the queue
	for i := 0; i < b.N; i++ {
		if err := q.Enqueue(i); err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := q.Dequeue(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQueueEnqueueDequeue(b *testing.B) {
	q := NewQueue[int](1000) // Fixed size queue
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Enqueue
		if err := q.Enqueue(i); err != nil {
			// If queue is full, dequeue one item and try again
			if _, err2 := q.Dequeue(); err2 != nil {
				b.Fatal(err2)
			}
			if err = q.Enqueue(i); err != nil {
				b.Fatal(err)
			}
		}

		// Occasionally dequeue to prevent queue from filling up
		if i%10 == 0 && !q.IsEmpty() {
			if _, err := q.Dequeue(); err != nil {
				b.Fatal(err)
			}
		}
	}
}

func BenchmarkQueuePeek(b *testing.B) {
	q := NewQueue[int](10)
	if err := q.Enqueue(42); err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := q.Peek(); err != nil {
			b.Fatal(err)
		}
	}
}
