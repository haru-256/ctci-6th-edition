package linked_list

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewLinkedList(t *testing.T) {
	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "creates empty list",
			test: func(t *testing.T) {
				list := NewLinkedList[int]()
				assert.NotNil(t, list)
				assert.Nil(t, list.Head)
				assert.Nil(t, list.Tail)
			},
		},
		{
			name: "creates empty string list",
			test: func(t *testing.T) {
				list := NewLinkedList[string]()
				assert.NotNil(t, list)
				assert.Nil(t, list.Head)
				assert.Nil(t, list.Tail)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.test)
	}
}

func TestNewNode(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected interface{}
	}{
		{"creates int node", 42, 42},
		{"creates string node", "hello", "hello"},
		{"creates float node", 3.14, 3.14},
		{"creates bool node", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case int:
				node := NewNode(v)
				assert.Equal(t, tt.expected, node.Value)
				assert.Nil(t, node.Prev)
				assert.Nil(t, node.Next)
			case string:
				node := NewNode(v)
				assert.Equal(t, tt.expected, node.Value)
				assert.Nil(t, node.Prev)
				assert.Nil(t, node.Next)
			case float64:
				node := NewNode(v)
				assert.Equal(t, tt.expected, node.Value)
				assert.Nil(t, node.Prev)
				assert.Nil(t, node.Next)
			case bool:
				node := NewNode(v)
				assert.Equal(t, tt.expected, node.Value)
				assert.Nil(t, node.Prev)
				assert.Nil(t, node.Next)
			}
		})
	}
}

func TestLinkedList_Prepend(t *testing.T) {
	tests := []struct {
		name           string
		initialValues  []int
		prependValue   int
		expectedValues []int
	}{
		{
			name:           "prepend to empty list",
			initialValues:  []int{},
			prependValue:   1,
			expectedValues: []int{1},
		},
		{
			name:           "prepend to single element list",
			initialValues:  []int{2},
			prependValue:   1,
			expectedValues: []int{1, 2},
		},
		{
			name:           "prepend to multiple element list",
			initialValues:  []int{3, 4, 5}, // These will be added as [5,4,3] due to prepend
			prependValue:   1,
			expectedValues: []int{1, 5, 4, 3}, // 1 prepended to [5,4,3]
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int]()

			// Setup initial list
			for _, v := range tt.initialValues {
				list.Prepend(v)
			}

			// Prepend the test value
			list.Prepend(tt.prependValue)

			// Verify the result
			values := collectValues(list)
			assert.Equal(t, tt.expectedValues, values)

			if len(tt.expectedValues) > 0 {
				assert.Equal(t, tt.expectedValues[0], list.Head.Value)
				assert.Equal(t, tt.expectedValues[len(tt.expectedValues)-1], list.Tail.Value)
			}
		})
	}
}

func TestLinkedList_Search(t *testing.T) {
	tests := []struct {
		name        string
		values      []int
		searchValue int
		expectFound bool
		expectedPos int // position in the list (0-indexed), -1 if not found
	}{
		{
			name:        "search in empty list",
			values:      []int{},
			searchValue: 1,
			expectFound: false,
			expectedPos: -1,
		},
		{
			name:        "search existing value at head",
			values:      []int{1, 2, 3},
			searchValue: 1,
			expectFound: true,
			expectedPos: 0,
		},
		{
			name:        "search existing value in middle",
			values:      []int{1, 2, 3},
			searchValue: 2,
			expectFound: true,
			expectedPos: 1,
		},
		{
			name:        "search existing value at tail",
			values:      []int{1, 2, 3},
			searchValue: 3,
			expectFound: true,
			expectedPos: 2,
		},
		{
			name:        "search non-existing value",
			values:      []int{1, 2, 3},
			searchValue: 4,
			expectFound: false,
			expectedPos: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int]()

			// Setup list (prepend in reverse order to get correct order)
			for i := len(tt.values) - 1; i >= 0; i-- {
				list.Prepend(tt.values[i])
			}

			result := list.Search(tt.searchValue)

			if tt.expectFound {
				require.NotNil(t, result)
				assert.Equal(t, tt.searchValue, result.Value)
			} else {
				assert.Nil(t, result)
			}
		})
	}
}

func TestLinkedList_Insert(t *testing.T) {
	tests := []struct {
		name           string
		initialValues  []int
		insertValue    int
		insertAfterIdx int // index of node to insert after, -1 for nil
		expectedValues []int
		expectError    bool
		expectedError  error
	}{
		{
			name:           "insert after nil node",
			initialValues:  []int{1, 2, 3},
			insertValue:    4,
			insertAfterIdx: -1,
			expectedValues: []int{1, 2, 3},
			expectError:    true,
			expectedError:  ErrorNodeIsNil,
		},
		{
			name:           "insert after head",
			initialValues:  []int{1, 3},
			insertValue:    2,
			insertAfterIdx: 0,
			expectedValues: []int{1, 2, 3},
			expectError:    false,
		},
		{
			name:           "insert after middle node",
			initialValues:  []int{1, 2, 4},
			insertValue:    3,
			insertAfterIdx: 1,
			expectedValues: []int{1, 2, 3, 4},
			expectError:    false,
		},
		{
			name:           "insert after tail",
			initialValues:  []int{1, 2, 3},
			insertValue:    4,
			insertAfterIdx: 2,
			expectedValues: []int{1, 2, 3, 4},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int]()

			// Setup initial list
			for i := len(tt.initialValues) - 1; i >= 0; i-- {
				list.Prepend(tt.initialValues[i])
			}

			// Find the node to insert after
			var afterNode *Node[int]
			if tt.insertAfterIdx >= 0 {
				afterNode = getNodeAtIndex(list, tt.insertAfterIdx)
			}

			err := list.Insert(tt.insertValue, afterNode)

			if tt.expectError {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				values := collectValues(list)
				assert.Equal(t, tt.expectedValues, values)

				// Verify tail is correctly updated
				if len(tt.expectedValues) > 0 {
					assert.Equal(t, tt.expectedValues[len(tt.expectedValues)-1], list.Tail.Value)
				}
			}
		})
	}
}

func TestLinkedList_Delete(t *testing.T) {
	tests := []struct {
		name           string
		initialValues  []int
		deleteValue    int
		expectedValues []int
		expectError    bool
		expectedError  error
	}{
		{
			name:           "delete from empty list",
			initialValues:  []int{},
			deleteValue:    1,
			expectedValues: []int{},
			expectError:    true,
			expectedError:  ErrorNodeNotFound,
		},
		{
			name:           "delete non-existing value",
			initialValues:  []int{1, 2, 3},
			deleteValue:    4,
			expectedValues: []int{1, 2, 3},
			expectError:    true,
			expectedError:  ErrorNodeNotFound,
		},
		{
			name:           "delete head from single element list",
			initialValues:  []int{1},
			deleteValue:    1,
			expectedValues: []int{},
			expectError:    false,
		},
		{
			name:           "delete head from multiple element list",
			initialValues:  []int{1, 2, 3},
			deleteValue:    1,
			expectedValues: []int{2, 3},
			expectError:    false,
		},
		{
			name:           "delete middle element",
			initialValues:  []int{1, 2, 3},
			deleteValue:    2,
			expectedValues: []int{1, 3},
			expectError:    false,
		},
		{
			name:           "delete tail",
			initialValues:  []int{1, 2, 3},
			deleteValue:    3,
			expectedValues: []int{1, 2},
			expectError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewLinkedList[int]()

			// Setup initial list
			for i := len(tt.initialValues) - 1; i >= 0; i-- {
				list.Prepend(tt.initialValues[i])
			}

			err := list.Delete(tt.deleteValue)

			if tt.expectError {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			values := collectValues(list)
			assert.Equal(t, tt.expectedValues, values)

			// Verify head and tail are correctly updated
			if len(tt.expectedValues) == 0 {
				assert.Nil(t, list.Head)
				assert.Nil(t, list.Tail)
			} else {
				assert.Equal(t, tt.expectedValues[0], list.Head.Value)
				assert.Equal(t, tt.expectedValues[len(tt.expectedValues)-1], list.Tail.Value)
			}
		})
	}
}

// Helper functions

// collectValues returns all values in the list from head to tail
func collectValues[T comparable](list *LinkedList[T]) []T {
	var values []T
	current := list.Head
	for current != nil {
		values = append(values, current.Value)
		current = current.Next
	}
	// Return empty slice instead of nil for consistent comparison
	if values == nil {
		return []T{}
	}
	return values
}

// getNodeAtIndex returns the node at the specified index (0-based)
func getNodeAtIndex[T comparable](list *LinkedList[T], index int) *Node[T] {
	current := list.Head
	for i := 0; i < index && current != nil; i++ {
		current = current.Next
	}
	return current
}
