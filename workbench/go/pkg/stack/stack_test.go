package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStack(t *testing.T) {
	t.Run("valid size", func(t *testing.T) {
		s := NewStack[int](10)
		assert.NotNil(t, s)
		assert.True(t, s.IsEmpty())
		assert.False(t, s.IsFull())
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

	item := 42
	require.NoError(t, s.Push(&item))
	assert.False(t, s.IsEmpty())

	_, err := s.Pop()
	require.NoError(t, err)
	assert.True(t, s.IsEmpty())
}

func TestIsFull(t *testing.T) {
	s := NewStack[int](2)
	assert.False(t, s.IsFull())

	item1 := 1
	require.NoError(t, s.Push(&item1))
	assert.False(t, s.IsFull())

	item2 := 2
	require.NoError(t, s.Push(&item2))
	assert.True(t, s.IsFull())
}

func TestPush(t *testing.T) {
	t.Run("successful push", func(t *testing.T) {
		s := NewStack[string](3)

		str1 := "hello"
		err := s.Push(&str1)
		assert.NoError(t, err)
		assert.False(t, s.IsEmpty())

		str2 := "world"
		err = s.Push(&str2)
		assert.NoError(t, err)
		assert.False(t, s.IsFull())

		str3 := "test"
		err = s.Push(&str3)
		assert.NoError(t, err)
		assert.True(t, s.IsFull())
	})

	t.Run("push to full stack", func(t *testing.T) {
		s := NewStack[int](2)

		item1 := 1
		require.NoError(t, s.Push(&item1))
		item2 := 2
		require.NoError(t, s.Push(&item2))
		assert.True(t, s.IsFull())

		item3 := 3
		err := s.Push(&item3)
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

		require.NoError(t, s.Push(&item1))
		require.NoError(t, s.Push(&item2))
		require.NoError(t, s.Push(&item3))

		// Pop in LIFO order
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 30, *popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 20, *popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 10, *popped)

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

		require.NoError(t, s.Push(&item1))
		require.NoError(t, s.Push(&item2))

		// Peek should return the top item without removing it
		peeked, err := s.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 20, *peeked)

		// Stack should remain unchanged
		assert.False(t, s.IsEmpty())

		// Peek again should return the same item
		peeked, err = s.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 20, *peeked)

		// Pop should return the same item that was peeked
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 20, *popped)
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
		temp := item // Create a copy to get address
		require.NoError(t, s.Push(&temp))
	}

	assert.True(t, s.IsFull())

	// Pop all items in reverse order (LIFO)
	for i := len(items) - 1; i >= 0; i-- {
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, items[i], *popped)
	}

	assert.True(t, s.IsEmpty())
}

func TestGenericTypes(t *testing.T) {
	t.Run("string stack", func(t *testing.T) {
		s := NewStack[string](2)

		str1 := "hello"
		str2 := "world"

		require.NoError(t, s.Push(&str1))
		require.NoError(t, s.Push(&str2))

		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, "world", *popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, "hello", *popped)
	})

	t.Run("struct stack", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}

		s := NewStack[Person](2)

		person1 := Person{Name: "Alice", Age: 30}
		person2 := Person{Name: "Bob", Age: 25}

		require.NoError(t, s.Push(&person1))
		require.NoError(t, s.Push(&person2))

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
			temp := item
			require.NoError(t, s.Push(&temp))
		}

		// Pop in reverse order
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 3.14, *popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, "string", *popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 42, *popped)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("single element stack", func(t *testing.T) {
		s := NewStack[int](1)

		item := 42
		require.NoError(t, s.Push(&item))
		assert.True(t, s.IsFull())

		// Try to push another item
		item2 := 43
		err := s.Push(&item2)
		assert.Equal(t, ErrorStackOverflow, err)

		// Pop the item
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 42, *popped)
		assert.True(t, s.IsEmpty())
	})

	t.Run("multiple push/pop cycles", func(t *testing.T) {
		s := NewStack[int](3)

		// First cycle
		for i := 1; i <= 3; i++ {
			item := i
			require.NoError(t, s.Push(&item))
		}

		for i := 3; i >= 1; i-- {
			popped, err := s.Pop()
			assert.NoError(t, err)
			assert.Equal(t, i, *popped)
		}
		assert.True(t, s.IsEmpty())

		// Second cycle
		for i := 4; i <= 6; i++ {
			item := i
			require.NoError(t, s.Push(&item))
		}

		for i := 6; i >= 4; i-- {
			popped, err := s.Pop()
			assert.NoError(t, err)
			assert.Equal(t, i, *popped)
		}
		assert.True(t, s.IsEmpty())
	})

	t.Run("alternating operations", func(t *testing.T) {
		s := NewStack[int](5)

		// Push, peek, push, pop pattern
		item1 := 1
		require.NoError(t, s.Push(&item1))

		peeked, err := s.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 1, *peeked)

		item2 := 2
		require.NoError(t, s.Push(&item2))

		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 2, *popped)

		// Stack should still have the first item
		assert.False(t, s.IsEmpty())

		peeked, err = s.Peek()
		assert.NoError(t, err)
		assert.Equal(t, 1, *peeked)
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
		require.NoError(t, s.Push(&item1))
		require.NoError(t, s.Push(&item2))
		assert.True(t, s.IsFull())

		// Multiple pushes should all fail
		for i := 3; i <= 5; i++ {
			item := i
			err := s.Push(&item)
			assert.Equal(t, ErrorStackOverflow, err)
		}

		// Stack should still be full with original items
		assert.True(t, s.IsFull())

		// Verify original items are still there
		popped, err := s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 2, *popped)

		popped, err = s.Pop()
		assert.NoError(t, err)
		assert.Equal(t, 1, *popped)
	})
}
