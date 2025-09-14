package stack

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestNewStack(t *testing.T) {
	t.Run("valid size", func(t *testing.T) {
		s := NewStack[int](10)
		assert.NotNil(t, s)
		assert.True(t, s.IsEmpty())
		assert.False(t, s.IsFull())
		assert.Equal(t, 10, s.Size())
		assert.Equal(t, 0, s.Count())
	})

	t.Run("panic on invalid size", func(t *testing.T) {
		assert.Panics(t, func() {
			NewStack[int](0)
		})
		assert.Panics(t, func() {
			NewStack[int](-1)
		})
	})
}

func TestIsEmpty(t *testing.T) {
	s := NewStack[int](5)
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Count())

	item := 42
	require.NoError(t, s.Push(item))
	assert.False(t, s.IsEmpty())
	assert.Equal(t, 1, s.Count())

	_, err := s.Pop()
	require.NoError(t, err)
	assert.True(t, s.IsEmpty())
	assert.Equal(t, 0, s.Count())
}

func TestIsFull(t *testing.T) {
	s := NewStack[int](2)
	assert.False(t, s.IsFull())
	assert.Equal(t, 0, s.Count())

	item1 := 1
	require.NoError(t, s.Push(item1))
	assert.False(t, s.IsFull())
	assert.Equal(t, 1, s.Count())

	item2 := 2
	require.NoError(t, s.Push(item2))
	assert.True(t, s.IsFull())
	assert.Equal(t, 2, s.Count())
	assert.Equal(t, s.Count(), s.Size())
}

func TestPush(t *testing.T) {
	t.Run("successful push", func(t *testing.T) {
		s := NewStack[string](3)

		str1 := "hello"
		err := s.Push(str1)
		assert.NoError(t, err)
		assert.False(t, s.IsEmpty())

		str2 := "world"
		err = s.Push(str2)
		assert.NoError(t, err)
		assert.False(t, s.IsFull())

		str3 := "test"
		err = s.Push(str3)
		assert.NoError(t, err)
		assert.True(t, s.IsFull())
	})

	t.Run("push to full stack", func(t *testing.T) {
		s := NewStack[int](2)

		item1 := 1
		require.NoError(t, s.Push(item1))
		item2 := 2
		require.NoError(t, s.Push(item2))
		assert.True(t, s.IsFull())

		item3 := 3
		err := s.Push(item3)
		assert.Equal(t, ErrorStackOverflow, err)
	})
}

func TestPop(t *testing.T) {
	t.Run("successful pop", func(t *testing.T) {
		s := NewStack[int](3)

		// Push items in order
		item1 := 10
		item2 := 20
		item3 := 30

		require.NoError(t, s.Push(item1))
		require.NoError(t, s.Push(item2))
		require.NoError(t, s.Push(item3))

		// Pop in LIFO order
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 30, popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 20, popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 10, popped)

		assert.True(t, s.IsEmpty())
	})

	t.Run("pop from empty stack", func(t *testing.T) {
		s := NewStack[int](3)

		_, err := s.Pop()
		assert.Equal(t, ErrorStackUnderflow, err)
	})
}

func TestPeek(t *testing.T) {
	t.Run("successful peek", func(t *testing.T) {
		s := NewStack[int](3)

		item1 := 10
		item2 := 20

		require.NoError(t, s.Push(item1))
		require.NoError(t, s.Push(item2))

		// Peek should return the top item without removing it
		peeked, err := s.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 20, peeked)

		// Stack should remain unchanged
		assert.False(t, s.IsEmpty())

		// Peek again should return the same item
		peeked, err = s.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 20, peeked)

		// Pop should return the same item that was peeked
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 20, popped)
	})

	t.Run("peek empty stack", func(t *testing.T) {
		s := NewStack[int](3)

		_, err := s.Peek()
		assert.Equal(t, ErrorStackUnderflow, err)
	})
}

func TestLIFOBehavior(t *testing.T) {
	s := NewStack[string](5)

	items := []string{"first", "second", "third", "fourth", "fifth"}

	// Push all items
	for _, item := range items {
		require.NoError(t, s.Push(item))
	}

	assert.True(t, s.IsFull())

	// Pop all items in reverse order (LIFO)
	for i := len(items) - 1; i >= 0; i-- {
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, items[i], popped)
	}

	assert.True(t, s.IsEmpty())
}

func TestGenericTypes(t *testing.T) {
	t.Run("string stack", func(t *testing.T) {
		s := NewStack[string](2)

		str1 := "hello"
		str2 := "world"

		require.NoError(t, s.Push(str1))
		require.NoError(t, s.Push(str2))

		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, "world", popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, "hello", popped)
	})

	t.Run("struct stack", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		s := NewStack[Person](2)

		person1 := Person{Name: "Alice", Age: 30}
		person2 := Person{Name: "Bob", Age: 25}

		require.NoError(t, s.Push(person1))
		require.NoError(t, s.Push(person2))

		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, "Bob", popped.Name)
		assert.Equal(t, 25, popped.Age)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, "Alice", popped.Name)
		assert.Equal(t, 30, popped.Age)
	})

	t.Run("interface{} stack", func(t *testing.T) {
		s := NewStack[interface{}](3)

		var items []interface{}
		items = append(items, 42)
		items = append(items, "string")
		items = append(items, 3.14)

		for _, item := range items {
			require.NoError(t, s.Push(item))
		}

		// Pop in reverse order
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 3.14, popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, "string", popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 42, popped)
	})
}

func FuzzStack_PushPop(f *testing.F) {
	f.Add(1)
	f.Add(42)
	f.Add(-7)
	f.Fuzz(func(t *testing.T, x int) {
		s := NewStack[int](10)
		require.NoError(t, s.Push(x))
		got, err := s.Pop()
		require.NoError(t, err)
		assert.Equal(t, x, got)
		assert.True(t, s.IsEmpty())
	})
}

func TestStack_ConcurrentPushPop(t *testing.T) {
	s := NewStack[int](100)
	var g errgroup.Group

	g.Go(func() error {
		for i := 0; i < 100; i++ {
			if err := s.Push(i); err != nil && err != ErrorStackOverflow {
				return err
			}
		}
		return nil
	})
	g.Go(func() error {
		for i := 0; i < 100; i++ {
			_, err := s.Pop()
			if err != nil && err != ErrorStackUnderflow {
				return err
			}
		}
		return nil
	})
	require.NoError(t, g.Wait(), "errgroup should not return error")
	// Not strictly thread-safe, but should not panic or deadlock
}

func TestStack_ConcurrentStress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	s := NewStack[int](1000)
	var g errgroup.Group
	const numGoroutines = 10
	const itemsPerGoroutine = 100

	// Multiple pushers
	for i := 0; i < numGoroutines; i++ {
		i := i
		g.Go(func() error {
			for j := 0; j < itemsPerGoroutine; j++ {
				if err := s.Push(i*itemsPerGoroutine + j); err != nil {
					return err
				}
			}
			return nil
		})
	}

	// Multiple poppers
	for i := 0; i < numGoroutines; i++ {
		g.Go(func() error {
			for j := 0; j < itemsPerGoroutine; j++ {
				_, err := s.Pop()
				if err != nil && err != ErrorStackUnderflow {
					return err
				}
			}
			return nil
		})
	}

	require.NoError(t, g.Wait(), "stress test should not fail")
}

func TestStack_ZeroAndPointerValues(t *testing.T) {
	s := NewStack[*int](3)
	var nilPtr *int
	require.NoError(t, s.Push(nilPtr))
	v, err := s.Pop()
	require.NoError(t, err)
	assert.Nil(t, v)

	s2 := NewStack[int](2)
	require.NoError(t, s2.Push(0))
	v2, err := s2.Pop()
	require.NoError(t, err)
	assert.Equal(t, 0, v2)
}

func TestEdgeCases(t *testing.T) {
	t.Run("single element stack", func(t *testing.T) {
		s := NewStack[int](1)

		item := 42
		require.NoError(t, s.Push(item))
		assert.True(t, s.IsFull())

		// Try to push another item
		item2 := 43
		err := s.Push(item2)
		assert.Equal(t, ErrorStackOverflow, err)

		// Pop the item
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 42, popped)
		assert.True(t, s.IsEmpty())
	})

	t.Run("multiple push/pop cycles", func(t *testing.T) {
		s := NewStack[int](3)

		// First cycle
		for i := 1; i <= 3; i++ {
			item := i
			require.NoError(t, s.Push(item))
		}

		for i := 3; i >= 1; i-- {
			popped, err := s.Pop()
			assert.NoError(t, err)
			assert.Equal(t, i, popped)
		}
		assert.True(t, s.IsEmpty())

		// Second cycle
		for i := 4; i <= 6; i++ {
			item := i
			require.NoError(t, s.Push(item))
		}

		for i := 6; i >= 4; i-- {
			popped, err := s.Pop()
			assert.NoError(t, err)
			assert.Equal(t, i, popped)
		}
		assert.True(t, s.IsEmpty())
	})

	t.Run("alternating operations", func(t *testing.T) {
		s := NewStack[int](5)

		// Push, peek, push, pop pattern
		item1 := 1
		require.NoError(t, s.Push(item1))

		peeked, err := s.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 1, peeked)

		item2 := 2
		require.NoError(t, s.Push(item2))

		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 2, popped)

		// Stack should still have the first item
		assert.False(t, s.IsEmpty())

		peeked, err = s.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 1, peeked)
	})
}

func TestErrorHandling(t *testing.T) {
	t.Run("multiple operations on empty stack", func(t *testing.T) {
		s := NewStack[int](3)

		// Multiple pops should all fail
		for i := 0; i < 3; i++ {
			_, err := s.Pop()
			assert.Equal(t, ErrorStackUnderflow, err)
		}

		// Multiple peeks should all fail
		for i := 0; i < 3; i++ {
			_, err := s.Peek()
			assert.Equal(t, ErrorStackUnderflow, err)
		}

		assert.True(t, s.IsEmpty())
	})

	t.Run("multiple operations on full stack", func(t *testing.T) {
		s := NewStack[int](2)

		// Fill the stack
		item1 := 1
		item2 := 2
		require.NoError(t, s.Push(item1))
		require.NoError(t, s.Push(item2))
		assert.True(t, s.IsFull())

		// Multiple pushes should all fail
		for i := 3; i <= 5; i++ {
			item := i
			err := s.Push(item)
			assert.Equal(t, ErrorStackOverflow, err)
		}

		// Stack should still be full with original items
		assert.True(t, s.IsFull())

		// Verify original items are still there
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 2, popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 1, popped)
	})
}

func TestSize(t *testing.T) {
	t.Run("returns correct capacity", func(t *testing.T) {
		sizes := []int{1, 5, 10, 100}

		for _, size := range sizes {
			s := NewStack[int](size)
			assert.Equal(t, size, s.Size())

			// Size should remain constant regardless of operations
			item := 42
			if size > 0 {
				require.NoError(t, s.Push(item))
				assert.Equal(t, size, s.Size())

				_, err := s.Pop()
				require.NoError(t, err)
				assert.Equal(t, size, s.Size())
			}
		}
	})
}

func TestCount(t *testing.T) {
	t.Run("tracks current item count", func(t *testing.T) {
		s := NewStack[int](5)

		// Initially empty
		assert.Equal(t, 0, s.Count())
		assert.True(t, s.IsEmpty())

		// Add items and verify count increases
		for i := 1; i <= 5; i++ {
			item := i * 10
			require.NoError(t, s.Push(item))
			assert.Equal(t, i, s.Count())
		}

		assert.True(t, s.IsFull())
		assert.Equal(t, 5, s.Count())

		// Remove items and verify count decreases
		for i := 4; i >= 0; i-- {
			_, err := s.Pop()
			require.NoError(t, err)
			assert.Equal(t, i, s.Count())
		}

		assert.True(t, s.IsEmpty())
		assert.Equal(t, 0, s.Count())
	})

	t.Run("count consistency with operations", func(t *testing.T) {
		s := NewStack[string](3)

		// Test alternating push/pop operations
		str1 := "first"
		require.NoError(t, s.Push(str1))
		assert.Equal(t, 1, s.Count())

		str2 := "second"
		require.NoError(t, s.Push(str2))
		assert.Equal(t, 2, s.Count())

		// Pop one
		_, err := s.Pop()
		require.NoError(t, err)
		assert.Equal(t, 1, s.Count())

		// Push another
		str3 := "third"
		require.NoError(t, s.Push(str3))
		assert.Equal(t, 2, s.Count())

		// Verify IsEmpty/IsFull consistency with Count
		assert.Equal(t, s.Count() == 0, s.IsEmpty())
		assert.Equal(t, s.Count() == s.Size(), s.IsFull())
	})
}

func TestSizeAndCountRelationship(t *testing.T) {
	s := NewStack[int](10)

	// Count should always be <= Size
	assert.LessOrEqual(t, s.Count(), s.Size())

	// Fill stack partially
	for i := 0; i < 7; i++ {
		item := i
		require.NoError(t, s.Push(item))
		assert.LessOrEqual(t, s.Count(), s.Size())
	}

	// Count should equal Size when full
	for i := 7; i < 10; i++ {
		item := i
		require.NoError(t, s.Push(item))
	}
	assert.Equal(t, s.Count(), s.Size())
	assert.True(t, s.IsFull())

	// Empty the stack
	for s.Count() > 0 {
		_, err := s.Pop()
		require.NoError(t, err)
		assert.LessOrEqual(t, s.Count(), s.Size())
	}

	assert.Equal(t, 0, s.Count())
	assert.True(t, s.IsEmpty())
}

func TestZeroValues(t *testing.T) {
	t.Run("int stack with zero values", func(t *testing.T) {
		s := NewStack[int](3)

		// Push zero values
		require.NoError(t, s.Push(0))
		require.NoError(t, s.Push(0))

		assert.Equal(t, 2, s.Count())
		assert.False(t, s.IsEmpty())

		// Pop zero values
		val, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 0, val)

		val, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 0, val)

		assert.True(t, s.IsEmpty())
	})

	t.Run("string stack with empty strings", func(t *testing.T) {
		s := NewStack[string](2)

		require.NoError(t, s.Push(""))
		require.NoError(t, s.Push("non-empty"))

		val, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, "non-empty", val)

		val, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, "", val)
	})
}

func TestStackWithPointers(t *testing.T) {
	s := NewStack[*int](3)

	val1 := 42
	val2 := 84
	var nilVal *int = nil

	// Push pointers including nil
	require.NoError(t, s.Push(&val1))
	require.NoError(t, s.Push(nilVal))
	require.NoError(t, s.Push(&val2))

	// Pop in LIFO order
	popped, err := s.Pop()
	assert.NoError(t, err)
	assert.Equal(t, &val2, popped)
	assert.Equal(t, 84, *popped)

	popped, err = s.Pop()
	assert.NoError(t, err)
	assert.Nil(t, popped)

	popped, err = s.Pop()
	assert.NoError(t, err)
	assert.Equal(t, &val1, popped)
	assert.Equal(t, 42, *popped)
}

// Benchmark tests
func BenchmarkStackPush(b *testing.B) {
	s := NewStack[int](b.N)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Push(i)
	}
}

func BenchmarkStackPop(b *testing.B) {
	s := NewStack[int](b.N)
	for i := 0; i < b.N; i++ {
		_ = s.Push(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = s.Pop()
	}
}

func BenchmarkStackPushPop(b *testing.B) {
	s := NewStack[int](1000) // Fixed size to avoid overflow
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = s.Push(i % 1000)
		if i%2 == 1 {
			_, _ = s.Pop()
		}
	}
}

func BenchmarkStackPeek(b *testing.B) {
	s := NewStack[int](1)
	_ = s.Push(42)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = s.Peek()
	}
}

// Example functions for documentation
func ExampleNewStack() {
	// Create a stack of integers with capacity 3
	stack := NewStack[int](3)

	// Push some values
	_ = stack.Push(1)
	_ = stack.Push(2)
	_ = stack.Push(3)

	// Check if full
	fmt.Println("Full:", stack.IsFull())
	fmt.Println("Count:", stack.Count())

	// Output:
	// Full: true
	// Count: 3
}

func ExampleStack_Push() {
	stack := NewStack[string](2)

	err := stack.Push("first")
	fmt.Println("Error:", err)

	err = stack.Push("second")
	fmt.Println("Error:", err)

	// This will cause an overflow
	err = stack.Push("third")
	fmt.Println("Error:", err)

	// Output:
	// Error: <nil>
	// Error: <nil>
	// Error: stack overflow
}

func ExampleStack_Pop() {
	stack := NewStack[int](3)
	_ = stack.Push(10)
	_ = stack.Push(20)
	_ = stack.Push(30)

	// Pop values in LIFO order
	val, _ := stack.Pop()
	fmt.Println("Popped:", val)

	val, _ = stack.Pop()
	fmt.Println("Popped:", val)

	val, _ = stack.Pop()
	fmt.Println("Popped:", val)

	// Output:
	// Popped: 30
	// Popped: 20
	// Popped: 10
}

func ExampleStack_Peek() {
	stack := NewStack[string](3)
	_ = stack.Push("bottom")
	_ = stack.Push("middle")
	_ = stack.Push("top")

	// Peek doesn't remove the item
	val, _ := stack.Peek()
	fmt.Println("Peeked:", val)
	fmt.Println("Count after peek:", stack.Count())

	// Pop actually removes it
	val, _ = stack.Pop()
	fmt.Println("Popped:", val)
	fmt.Println("Count after pop:", stack.Count())

	// Output:
	// Peeked: top
	// Count after peek: 3
	// Popped: top
	// Count after pop: 2
}
