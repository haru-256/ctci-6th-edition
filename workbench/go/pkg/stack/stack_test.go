package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewStack(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "creates empty stack",
			test: func(t *testing.T) {
				s := NewStack(10)
				require.NotNil(t, s)
				assert.True(t, s.IsEmpty())
				assert.Equal(t, 0, s.top)
				assert.Equal(t, 10, s.size)
				assert.Equal(t, 10, len(s.items))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestStack_IsEmpty(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() *Stack
		expected bool
	}{
		{
			name: "empty stack returns true",
			setup: func() *Stack {
				return NewStack(10)
			},
			expected: true,
		},
		{
			name: "non-empty stack returns false",
			setup: func() *Stack {
				s := NewStack(10)
				_ = s.Push("item")
				return s
			},
			expected: false,
		},
		{
			name: "stack after push and pop returns true",
			setup: func() *Stack {
				s := NewStack(10)
				_ = s.Push("item")
				_, err := s.Pop()
				require.NoError(t, err)
				return s
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := tt.setup()
			result := stack.IsEmpty()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStack_Push(t *testing.T) {
	tests := []struct {
		name      string
		items     []any
		checkFunc func(t *testing.T, s *Stack)
	}{
		{
			name:  "push single item",
			items: []any{"hello"},
			checkFunc: func(t *testing.T, s *Stack) {
				assert.False(t, s.IsEmpty())
				assert.Equal(t, 1, s.top)
				assert.Equal(t, "hello", s.items[0])
			},
		},
		{
			name:  "push multiple items",
			items: []any{"first", "second", "third"},
			checkFunc: func(t *testing.T, s *Stack) {
				assert.False(t, s.IsEmpty())
				assert.Equal(t, 3, s.top)
				assert.Equal(t, "first", s.items[0])
				assert.Equal(t, "second", s.items[1])
				assert.Equal(t, "third", s.items[2])
			},
		},
		{
			name:  "push different types",
			items: []any{42, "string", 3.14, true},
			checkFunc: func(t *testing.T, s *Stack) {
				assert.False(t, s.IsEmpty())
				assert.Equal(t, 4, s.top)
				assert.Equal(t, 42, s.items[0])
				assert.Equal(t, "string", s.items[1])
				assert.Equal(t, 3.14, s.items[2])
				assert.Equal(t, true, s.items[3])
			},
		},
		{
			name:  "push nil item",
			items: []any{nil},
			checkFunc: func(t *testing.T, s *Stack) {
				assert.False(t, s.IsEmpty())
				assert.Equal(t, 1, s.top)
				assert.Nil(t, s.items[0])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := NewStack(10)

			for _, item := range tt.items {
				err := stack.Push(item)
				require.NoError(t, err)
			}

			tt.checkFunc(t, stack)
		})
	}
}

func TestStack_PushOverflow(t *testing.T) {
	tests := []struct {
		name      string
		stackSize int
		pushCount int
	}{
		{
			name:      "push to exactly full stack",
			stackSize: 3,
			pushCount: 3,
		},
		{
			name:      "push beyond stack capacity",
			stackSize: 2,
			pushCount: 5,
		},
		{
			name:      "push to size 1 stack",
			stackSize: 1,
			pushCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := NewStack(tt.stackSize)

			var err error
			for i := 0; i < tt.pushCount; i++ {
				err = stack.Push(i)
				if i < tt.stackSize {
					require.NoError(t, err, "Push %d should succeed", i)
					assert.Equal(t, i+1, stack.top)
				} else {
					require.Error(t, err, "Push %d should fail due to overflow", i)
					assert.ErrorIs(t, err, ErrorStackOverflow)
					assert.Equal(t, tt.stackSize, stack.top, "Stack top should remain at capacity")
				}
			}

			// Verify stack is at capacity but not beyond
			assert.Equal(t, tt.stackSize, stack.top)
			assert.Equal(t, tt.stackSize, len(stack.items))
		})
	}
}

func TestStack_Pop(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *Stack
		expectedItem  any
		expectedError error
		checkStack    func(t *testing.T, s *Stack)
	}{
		{
			name: "pop from empty stack returns error",
			setup: func() *Stack {
				return NewStack(10)
			},
			expectedItem:  nil,
			expectedError: ErrorStackUnderflow,
			checkStack: func(t *testing.T, s *Stack) {
				assert.True(t, s.IsEmpty())
			},
		},
		{
			name: "pop single item from stack",
			setup: func() *Stack {
				s := NewStack(10)
				_ = s.Push("item")
				return s
			},
			expectedItem:  "item",
			expectedError: nil,
			checkStack: func(t *testing.T, s *Stack) {
				assert.True(t, s.IsEmpty())
			},
		},
		{
			name: "pop from stack with multiple items returns last pushed",
			setup: func() *Stack {
				s := NewStack(10)
				_ = s.Push("first")
				_ = s.Push("second")
				_ = s.Push("third")
				return s
			},
			expectedItem:  "third",
			expectedError: nil,
			checkStack: func(t *testing.T, s *Stack) {
				assert.False(t, s.IsEmpty())
				assert.Equal(t, 2, s.top)
				assert.Equal(t, "first", s.items[0])
				assert.Equal(t, "second", s.items[1])
			},
		},
		{
			name: "pop nil item",
			setup: func() *Stack {
				s := NewStack(10)
				_ = s.Push(nil)
				return s
			},
			expectedItem:  nil,
			expectedError: nil,
			checkStack: func(t *testing.T, s *Stack) {
				assert.True(t, s.IsEmpty())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := tt.setup()

			item, err := stack.Pop()

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedItem, item)
			tt.checkStack(t, stack)
		})
	}
}

func TestStack_Peek(t *testing.T) {
	tests := []struct {
		name          string
		setup         func() *Stack
		expectedItem  any
		expectedError error
		checkStack    func(t *testing.T, s *Stack) // Peek should not modify the stack
	}{
		{
			name: "peek empty stack returns error",
			setup: func() *Stack {
				return NewStack(10)
			},
			expectedItem:  nil,
			expectedError: ErrorStackUnderflow,
			checkStack: func(t *testing.T, s *Stack) {
				assert.True(t, s.IsEmpty())
			},
		},
		{
			name: "peek single item stack",
			setup: func() *Stack {
				s := NewStack(10)
				_ = s.Push("item")
				return s
			},
			expectedItem:  "item",
			expectedError: nil,
			checkStack: func(t *testing.T, s *Stack) {
				assert.False(t, s.IsEmpty())
				assert.Equal(t, 1, s.top)
				assert.Equal(t, "item", s.items[0])
			},
		},
		{
			name: "peek stack with multiple items returns last pushed",
			setup: func() *Stack {
				s := NewStack(10)
				_ = s.Push("first")
				_ = s.Push("second")
				_ = s.Push("third")
				return s
			},
			expectedItem:  "third",
			expectedError: nil,
			checkStack: func(t *testing.T, s *Stack) {
				assert.False(t, s.IsEmpty())
				assert.Equal(t, 3, s.top)
				assert.Equal(t, "first", s.items[0])
				assert.Equal(t, "second", s.items[1])
				assert.Equal(t, "third", s.items[2])
			},
		},
		{
			name: "peek nil item",
			setup: func() *Stack {
				s := NewStack(10)
				_ = s.Push(nil)
				return s
			},
			expectedItem:  nil,
			expectedError: nil,
			checkStack: func(t *testing.T, s *Stack) {
				assert.False(t, s.IsEmpty())
				assert.Equal(t, 1, s.top)
				assert.Nil(t, s.items[0])
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := tt.setup()
			originalTop := stack.top

			item, err := stack.Peek()

			if tt.expectedError != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedItem, item)

			// Ensure peek doesn't modify the stack
			assert.Equal(t, originalTop, stack.top)
			tt.checkStack(t, stack)
		})
	}
}

func TestStack_SequentialOperations(t *testing.T) {
	tests := []struct {
		name       string
		operations func(t *testing.T, s *Stack)
	}{
		{
			name: "push-pop-push sequence",
			operations: func(t *testing.T, s *Stack) {
				// Push first item
				err := s.Push("first")
				require.NoError(t, err)
				assert.False(t, s.IsEmpty())

				// Peek should return first item
				item, err := s.Peek()
				require.NoError(t, err)
				assert.Equal(t, "first", item)

				// Pop should return first item
				item, err = s.Pop()
				require.NoError(t, err)
				assert.Equal(t, "first", item)
				assert.True(t, s.IsEmpty())

				// Push second item
				err = s.Push("second")
				require.NoError(t, err)
				assert.False(t, s.IsEmpty())

				// Peek should return second item
				item, err = s.Peek()
				require.NoError(t, err)
				assert.Equal(t, "second", item)
			},
		},
		{
			name: "multiple push and pop operations",
			operations: func(t *testing.T, s *Stack) {
				items := []string{"a", "b", "c", "d", "e"}

				// Push all items
				for _, item := range items {
					err := s.Push(item)
					require.NoError(t, err)
				}

				// Stack should contain all items
				assert.Equal(t, len(items), s.top)

				// Pop all items in LIFO order
				for i := len(items) - 1; i >= 0; i-- {
					item, err := s.Pop()
					require.NoError(t, err)
					assert.Equal(t, items[i], item)
				}

				// Stack should be empty
				assert.True(t, s.IsEmpty())
			},
		},
		{
			name: "alternating push and peek operations",
			operations: func(t *testing.T, s *Stack) {
				// Push and peek multiple times
				err := s.Push(1)
				require.NoError(t, err)
				item, err := s.Peek()
				require.NoError(t, err)
				assert.Equal(t, 1, item)

				err = s.Push(2)
				require.NoError(t, err)
				item, err = s.Peek()
				require.NoError(t, err)
				assert.Equal(t, 2, item)

				err = s.Push(3)
				require.NoError(t, err)
				item, err = s.Peek()
				require.NoError(t, err)
				assert.Equal(t, 3, item)

				// Stack should still have all 3 items
				assert.Equal(t, 3, s.top)
				assert.False(t, s.IsEmpty())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := NewStack(10)
			tt.operations(t, stack)
		})
	}
}

func TestStack_ErrorConditions(t *testing.T) {
	tests := []struct {
		name      string
		operation func(t *testing.T, s *Stack)
	}{
		{
			name: "multiple pops from empty stack",
			operation: func(t *testing.T, s *Stack) {
				// First pop should return error
				item, err := s.Pop()
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorStackUnderflow)
				assert.Nil(t, item)

				// Second pop should also return error
				item, err = s.Pop()
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorStackUnderflow)
				assert.Nil(t, item)

				assert.True(t, s.IsEmpty())
			},
		},
		{
			name: "multiple peeks from empty stack",
			operation: func(t *testing.T, s *Stack) {
				// First peek should return error
				item, err := s.Peek()
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorStackUnderflow)
				assert.Nil(t, item)

				// Second peek should also return error
				item, err = s.Peek()
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorStackUnderflow)
				assert.Nil(t, item)

				assert.True(t, s.IsEmpty())
			},
		},
		{
			name: "pop until empty then continue popping",
			operation: func(t *testing.T, s *Stack) {
				// Add some items
				err := s.Push("a")
				require.NoError(t, err)
				err = s.Push("b")
				require.NoError(t, err)

				// Pop all items
				_, err = s.Pop()
				require.NoError(t, err)
				_, err = s.Pop()
				require.NoError(t, err)

				assert.True(t, s.IsEmpty())

				// Try to pop from empty stack
				item, err := s.Pop()
				require.Error(t, err)
				assert.ErrorIs(t, err, ErrorStackUnderflow)
				assert.Nil(t, item)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stack := NewStack(10)
			tt.operation(t, stack)
		})
	}
}
