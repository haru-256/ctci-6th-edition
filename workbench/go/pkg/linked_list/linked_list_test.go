package linked_list

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
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
				assert.Nil(t, list.Head())
				assert.Nil(t, list.Tail())
			},
		},
		{
			name: "creates empty string list",
			test: func(t *testing.T) {
				list := NewLinkedList[string]()
				assert.NotNil(t, list)
				assert.Nil(t, list.Head())
				assert.Nil(t, list.Tail())
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
				assert.Equal(t, tt.expectedValues[0], list.Head().Value)
				assert.Equal(t, tt.expectedValues[len(tt.expectedValues)-1], list.Tail().Value)
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
					assert.Equal(t, tt.expectedValues[len(tt.expectedValues)-1], list.Tail().Value)
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
				assert.Nil(t, list.Head())
				assert.Nil(t, list.Tail())
			} else {
				assert.Equal(t, tt.expectedValues[0], list.Head().Value)
				assert.Equal(t, tt.expectedValues[len(tt.expectedValues)-1], list.Tail().Value)
			}
		})
	}
}

// Helper functions

// collectValues returns all values in the list from head to tail
func collectValues[T comparable](list *LinkedList[T]) []T {
	values := []T{}
	current := list.Head()
	for current != nil {
		values = append(values, current.Value)
		current = current.Next
	}
	return values
}

// getNodeAtIndex returns the node at the specified index (0-based)
func getNodeAtIndex[T comparable](list *LinkedList[T], index int) *Node[T] {
	current := list.Head()
	for i := 0; i < index && current != nil; i++ {
		current = current.Next
	}
	return current
}

// TestLinkedList_ConcurrentSearch tests that multiple goroutines can safely
// read from the linked list simultaneously without data races.
func TestLinkedList_ConcurrentSearch(t *testing.T) {
	list := NewLinkedList[int]()

	// Prepare initial data
	for i := 1; i <= 10; i++ {
		list.Prepend(i)
	}

	const numReaders = 10
	const readsPerReader = 100

	var wg sync.WaitGroup

	// Start multiple readers
	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()

			for j := 0; j < readsPerReader; j++ {
				// Search for values that should exist
				node := list.Search(readerID%10 + 1)
				if readerID%10+1 <= 10 {
					assert.NotNil(t, node, "Reader %d should find value %d", readerID, readerID%10+1)
				}

				// Search for values that don't exist
				node = list.Search(100 + readerID)
				assert.Nil(t, node, "Reader %d should not find value %d", readerID, 100+readerID)

				// Access head and tail
				head := list.Head()
				tail := list.Tail()
				assert.NotNil(t, head, "Head should not be nil")
				assert.NotNil(t, tail, "Tail should not be nil")
			}
		}(i)
	}

	wg.Wait()
}

// TestLinkedList_ConcurrentWrites tests that multiple goroutines can safely
// write to the linked list without causing data races or corruption.
func TestLinkedList_ConcurrentWrites(t *testing.T) {
	list := NewLinkedList[int]()

	const numWriters = 10
	const writesPerWriter = 50

	var g errgroup.Group

	// Start multiple writers performing prepends
	for i := 0; i < numWriters; i++ {
		writerID := i
		g.Go(func() error {
			for j := 0; j < writesPerWriter; j++ {
				list.Prepend(writerID*writesPerWriter + j)
			}
			return nil
		})
	}

	// Start multiple writers performing deletes
	for i := 0; i < numWriters/2; i++ {
		writerID := i
		g.Go(func() error {
			for j := 0; j < writesPerWriter/2; j++ {
				// Try to delete values (some may not exist, which is OK)
				_ = list.Delete(writerID*writesPerWriter/2 + j)
			}
			return nil
		})
	}

	require.NoError(t, g.Wait())

	// Verify list integrity - should be able to traverse without panics
	values := collectValues(list)
	assert.GreaterOrEqual(t, len(values), 0, "List should be accessible after concurrent operations")

	// Verify head and tail consistency
	head := list.Head()
	tail := list.Tail()
	if len(values) == 0 {
		assert.Nil(t, head, "Head should be nil for empty list")
		assert.Nil(t, tail, "Tail should be nil for empty list")
	} else {
		assert.NotNil(t, head, "Head should not be nil for non-empty list")
		assert.NotNil(t, tail, "Tail should not be nil for non-empty list")
		assert.Equal(t, values[0], head.Value, "Head value should match first collected value")
		assert.Equal(t, values[len(values)-1], tail.Value, "Tail value should match last collected value")
	}
}

// TestLinkedList_MixedConcurrentOperations tests a realistic scenario with
// mixed concurrent reads and writes.
func TestLinkedList_MixedConcurrentOperations(t *testing.T) {
	list := NewLinkedList[int]()

	// Initialize with some data
	for i := 1; i <= 20; i++ {
		list.Prepend(i)
	}

	const numOperations = 100
	var g errgroup.Group

	// Concurrent readers
	g.Go(func() error {
		for i := 0; i < numOperations; i++ {
			// Search for existing values
			node := list.Search(i%20 + 1)
			if node != nil {
				assert.GreaterOrEqual(t, node.Value, 1)
				assert.LessOrEqual(t, node.Value, 20)
			}

			// Access head/tail
			head := list.Head()
			tail := list.Tail()
			if head != nil && tail != nil {
				assert.NotNil(t, head)
				assert.NotNil(t, tail)
			}
		}
		return nil
	})

	// Concurrent prependers
	g.Go(func() error {
		for i := 0; i < numOperations/2; i++ {
			list.Prepend(100 + i)
		}
		return nil
	})

	// Concurrent inserters (find nodes and insert after them)
	g.Go(func() error {
		for i := 0; i < numOperations/4; i++ {
			// Find a node to insert after
			target := i%20 + 1
			node := list.Search(target)
			if node != nil {
				err := list.Insert(200+i, node)
				assert.NoError(t, err)
			}
		}
		return nil
	})

	// Concurrent deleters
	g.Go(func() error {
		for i := 0; i < numOperations/4; i++ {
			// Try to delete various values
			_ = list.Delete(i%20 + 1) // Some may not exist anymore
			if err := list.Delete(i%20 + 1); err != nil && err != ErrorNodeNotFound {
				return err
			}
		}
		return nil
	})

	require.NoError(t, g.Wait())

	// Final integrity check
	values := collectValues(list)
	assert.GreaterOrEqual(t, len(values), 0, "List should be accessible after mixed operations")

	// Verify list structure integrity
	current := list.Head()
	count := 0
	for current != nil {
		count++
		if current.Next != nil {
			assert.Equal(t, current, current.Next.Prev, "Forward/backward links should be consistent")
		}
		current = current.Next

		// Prevent infinite loops in case of corruption
		assert.Less(t, count, 1000, "List traversal should not exceed reasonable bounds")
	}
	assert.Equal(t, len(values), count, "Collected values count should match traversal count")
}

// TestLinkedList_ConcurrentInsertDelete tests the scenario where insertions
// and deletions happen concurrently, which exercises the most complex locking behavior.
func TestLinkedList_ConcurrentInsertDelete(t *testing.T) {
	list := NewLinkedList[int]()

	// Initialize with sequential values
	for i := 10; i >= 1; i-- {
		list.Prepend(i)
	}

	const numOperations = 50
	var g errgroup.Group

	// Concurrent insertions after specific nodes
	g.Go(func() error {
		for i := 0; i < numOperations; i++ {
			target := (i % 10) + 1
			node := list.Search(target)
			if node != nil {
				err := list.Insert(1000+i, node)
				assert.NoError(t, err)
			}
		}
		return nil
	})

	// Concurrent deletions
	g.Go(func() error {
		for i := 0; i < numOperations; i++ {
			target := (i % 10) + 1
			_ = list.Delete(target) // May fail if already deleted
		}
		return nil
	})

	// Concurrent prepends to keep adding data
	g.Go(func() error {
		for i := 0; i < numOperations; i++ {
			list.Prepend(2000 + i)
		}
		return nil
	})

	require.NoError(t, g.Wait())

	// Verify final state
	values := collectValues(list)
	assert.GreaterOrEqual(t, len(values), 0, "List should be accessible")

	// Check that all prepended values from the third goroutine are present
	foundPrependedValues := 0
	for _, value := range values {
		if value >= 2000 && value < 2000+numOperations {
			foundPrependedValues++
		}
	}
	assert.Equal(t, numOperations, foundPrependedValues, "All prepended values should be present")
}
