package hashtable

import (
	"testing"

	l "github.com/haru-256/ctci-6th-edition/pkg/linked_list"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHashChainTable(t *testing.T) {
	tests := []struct {
		name    string
		maxSize int64
	}{
		{"small table", 5},
		{"medium table", 100},
		{"large table", 1000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewHashChainTable[string](tt.maxSize)

			assert.NotNil(t, table)
			assert.Equal(t, tt.maxSize, table.MaxSize)
			assert.Equal(t, 0, table.Size())
			assert.Equal(t, int(tt.maxSize), len(table.Table))

			// Verify all buckets are initially nil
			for i := int64(0); i < tt.maxSize; i++ {
				assert.Nil(t, table.Table[i])
			}
		})
	}
}

func TestHashChainTable_Size(t *testing.T) {
	tests := []struct {
		name         string
		maxSize      int64
		insertValues []string
		expectedSize int
	}{
		{
			name:         "empty table",
			maxSize:      10,
			insertValues: []string{},
			expectedSize: 0,
		},
		{
			name:         "single element",
			maxSize:      10,
			insertValues: []string{"apple"},
			expectedSize: 1,
		},
		{
			name:         "multiple elements",
			maxSize:      10,
			insertValues: []string{"apple", "banana", "cherry"},
			expectedSize: 3,
		},
		{
			name:         "many elements",
			maxSize:      5,
			insertValues: []string{"a", "b", "c", "d", "e", "f", "g", "h"},
			expectedSize: 8,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewHashChainTable[string](tt.maxSize)

			for _, value := range tt.insertValues {
				err := table.Insert(value)
				require.NoError(t, err)
			}

			assert.Equal(t, tt.expectedSize, table.Size())
		})
	}
}

func TestHashChainTable_Insert_Strings(t *testing.T) {
	tests := []struct {
		name           string
		maxSize        int64
		insertValues   []string
		expectedErrors []error
		finalSize      int
	}{
		{
			name:           "insert strings",
			maxSize:        10,
			insertValues:   []string{"apple", "banana", "cherry"},
			expectedErrors: []error{nil, nil, nil},
			finalSize:      3,
		},
		{
			name:           "insert duplicate strings",
			maxSize:        5,
			insertValues:   []string{"apple", "banana", "apple"},
			expectedErrors: []error{nil, nil, ErrorAlreadyExists},
			finalSize:      2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewHashChainTable[string](tt.maxSize)
			for i, value := range tt.insertValues {
				err := table.Insert(value)
				if tt.expectedErrors[i] != nil {
					assert.ErrorIs(t, err, tt.expectedErrors[i])
				} else {
					assert.NoError(t, err)
				}
			}
			assert.Equal(t, tt.finalSize, table.Size())
		})
	}
}

func TestHashChainTable_Insert_Integers(t *testing.T) {
	tests := []struct {
		name           string
		maxSize        int64
		insertValues   []int
		expectedErrors []error
		finalSize      int
	}{
		{
			name:           "insert integers",
			maxSize:        5,
			insertValues:   []int{1, 2, 3, 4, 5},
			expectedErrors: []error{nil, nil, nil, nil, nil},
			finalSize:      5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewHashChainTable[int](tt.maxSize)
			for i, value := range tt.insertValues {
				err := table.Insert(value)
				if tt.expectedErrors[i] != nil {
					assert.ErrorIs(t, err, tt.expectedErrors[i])
				} else {
					assert.NoError(t, err)
				}
			}
			assert.Equal(t, tt.finalSize, table.Size())
		})
	}
}

func TestHashChainTable_Insert_Floats(t *testing.T) {
	tests := []struct {
		name           string
		maxSize        int64
		insertValues   []float64
		expectedErrors []error
		finalSize      int
	}{
		{
			name:           "insert floats",
			maxSize:        3,
			insertValues:   []float64{1.1, 2.2, 3.3},
			expectedErrors: []error{nil, nil, nil},
			finalSize:      3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewHashChainTable[float64](tt.maxSize)
			for i, value := range tt.insertValues {
				err := table.Insert(value)
				if tt.expectedErrors[i] != nil {
					assert.ErrorIs(t, err, tt.expectedErrors[i])
				} else {
					assert.NoError(t, err)
				}
			}
			assert.Equal(t, tt.finalSize, table.Size())
		})
	}
}

func TestHashChainTable_Insert_UnsupportedType(t *testing.T) {
	table := NewHashChainTable[bool](5)
	err := table.Insert(true)
	assert.ErrorIs(t, err, ErrorUnsupportedValueType)
	assert.Equal(t, 0, table.Size())
}

func TestHashChainTable_Search(t *testing.T) {
	tests := []struct {
		name          string
		maxSize       int64
		insertValues  []string
		searchValue   string
		expectFound   bool
		expectError   bool
		expectedError error
	}{
		{
			name:         "search in empty table",
			maxSize:      10,
			insertValues: []string{},
			searchValue:  "apple",
			expectFound:  false,
			expectError:  false,
		},
		{
			name:         "search existing value",
			maxSize:      10,
			insertValues: []string{"apple", "banana", "cherry"},
			searchValue:  "banana",
			expectFound:  true,
			expectError:  false,
		},
		{
			name:         "search non-existing value",
			maxSize:      10,
			insertValues: []string{"apple", "banana", "cherry"},
			searchValue:  "grape",
			expectFound:  false,
			expectError:  false,
		},
		{
			name:         "search with collision",
			maxSize:      2, // Small table to force collisions
			insertValues: []string{"a", "b", "c", "d"},
			searchValue:  "c",
			expectFound:  true,
			expectError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewHashChainTable[string](tt.maxSize)

			// Insert values
			for _, value := range tt.insertValues {
				err := table.Insert(value)
				require.NoError(t, err)
			}

			// Search for value
			node, err := table.Search(tt.searchValue)

			if tt.expectError {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				if tt.expectFound {
					require.NotNil(t, node)
					assert.Equal(t, tt.searchValue, node.Value)
				} else {
					assert.Nil(t, node)
				}
			}
		})
	}
}

func TestHashChainTable_Search_UnsupportedType(t *testing.T) {
	table := NewHashChainTable[bool](10)

	_, err := table.Search(true)
	assert.ErrorIs(t, err, ErrorUnsupportedValueType)
}

func TestHashChainTable_Delete(t *testing.T) {
	tests := []struct {
		name          string
		maxSize       int64
		insertValues  []string
		deleteValue   string
		expectError   bool
		expectedError error
		finalSize     int
	}{
		{
			name:          "delete from empty table",
			maxSize:       10,
			insertValues:  []string{},
			deleteValue:   "apple",
			expectError:   true, // Delete returns ErrorNodeNotFound for non-existing values
			expectedError: ErrorNodeNotFound,
			finalSize:     0,
		},
		{
			name:         "delete existing value",
			maxSize:      10,
			insertValues: []string{"apple", "banana", "cherry"},
			deleteValue:  "banana",
			expectError:  false,
			finalSize:    2,
		},
		{
			name:          "delete non-existing value",
			maxSize:       10,
			insertValues:  []string{"apple", "banana", "cherry"},
			deleteValue:   "grape",
			expectError:   true, // Delete returns ErrorNodeNotFound for non-existing values
			expectedError: ErrorNodeNotFound,
			finalSize:     3,
		},
		{
			name:         "delete with collision",
			maxSize:      2, // Small table to force collisions
			insertValues: []string{"a", "b", "c", "d"},
			deleteValue:  "c",
			expectError:  false,
			finalSize:    3,
		},
		{
			name:         "delete all from single bucket",
			maxSize:      1, // Force all values into one bucket
			insertValues: []string{"a", "b", "c"},
			deleteValue:  "b",
			expectError:  false,
			finalSize:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := NewHashChainTable[string](tt.maxSize)

			// Insert values
			for _, value := range tt.insertValues {
				err := table.Insert(value)
				require.NoError(t, err)
			}

			// Delete value
			err := table.Delete(tt.deleteValue)

			if tt.expectError {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.finalSize, table.Size())

			// Verify the value is actually deleted
			if !tt.expectError && tt.finalSize < len(tt.insertValues) {
				var node *l.Node[string]
				node, err = table.Search(tt.deleteValue)
				assert.NoError(t, err)
				assert.Nil(t, node)
			}
		})
	}
}

func TestHashChainTable_Delete_UnsupportedType(t *testing.T) {
	table := NewHashChainTable[bool](10)

	err := table.Delete(true)
	assert.ErrorIs(t, err, ErrorUnsupportedValueType)
}

func TestHashChainTable_Concurrency(t *testing.T) {
	table := NewHashChainTable[int](100)

	// Test concurrent insertions
	done := make(chan bool)

	// Insert values concurrently
	for i := 0; i < 10; i++ {
		go func(start int) {
			for j := start * 10; j < (start+1)*10; j++ {
				err := table.Insert(j)
				assert.NoError(t, err)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all values were inserted
	assert.Equal(t, 100, table.Size())

	// Verify all values can be found
	for i := 0; i < 100; i++ {
		node, err := table.Search(i)
		assert.NoError(t, err)
		assert.NotNil(t, node)
		assert.Equal(t, i, node.Value)
	}
}

func TestHashChainTable_GetHash(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		expectErr bool
	}{
		{"hash string", "hello", false},
		{"hash int", 42, false},
		{"hash float64", 3.14, false},
		{"hash unsupported bool", true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch v := tt.value.(type) {
			case string:
				table := NewHashChainTable[string](10)
				hash, err := table.getHash(v)
				if tt.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, hash)
				}
			case int:
				table := NewHashChainTable[int](10)
				hash, err := table.getHash(v)
				if tt.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, hash)
				}
			case float64:
				table := NewHashChainTable[float64](10)
				hash, err := table.getHash(v)
				if tt.expectErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, hash)
				}
			case bool:
				table := NewHashChainTable[bool](10)
				hash, err := table.getHash(v)
				if tt.expectErr {
					assert.ErrorIs(t, err, ErrorUnsupportedValueType)
					assert.Zero(t, hash)
				} else {
					assert.NoError(t, err)
					assert.NotZero(t, hash)
				}
			}
		})
	}
}

func TestHashChainTable_HashConsistency(t *testing.T) {
	table := NewHashChainTable[string](10)

	testValues := []string{"apple", "banana", "cherry", "date"}

	// Get hashes multiple times and verify consistency
	for _, value := range testValues {
		hash1, err1 := table.getHash(value)
		require.NoError(t, err1)

		hash2, err2 := table.getHash(value)
		require.NoError(t, err2)

		assert.Equal(t, hash1, hash2, "Hash should be consistent for value: %s", value)
	}
}

func TestHashChainTable_CompleteWorkflow(t *testing.T) {
	// Test a complete workflow: insert, search, delete
	table := NewHashChainTable[string](5)
	values := []string{"apple", "banana", "cherry", "date", "elderberry"}

	// Insert all values
	for _, value := range values {
		err := table.Insert(value)
		require.NoError(t, err)
	}
	assert.Equal(t, len(values), table.Size())

	// Search for all values
	for _, value := range values {
		node, err := table.Search(value)
		require.NoError(t, err)
		require.NotNil(t, node)
		assert.Equal(t, value, node.Value)
	}

	// Delete some values
	deleteValues := []string{"banana", "date"}
	for _, value := range deleteValues {
		err := table.Delete(value)
		require.NoError(t, err)
	}
	assert.Equal(t, len(values)-len(deleteValues), table.Size())

	// Verify deleted values are gone
	for _, value := range deleteValues {
		node, err := table.Search(value)
		require.NoError(t, err)
		assert.Nil(t, node)
	}

	// Verify remaining values are still there
	remainingValues := []string{"apple", "cherry", "elderberry"}
	for _, value := range remainingValues {
		node, err := table.Search(value)
		require.NoError(t, err)
		require.NotNil(t, node)
		assert.Equal(t, value, node.Value)
	}
}
