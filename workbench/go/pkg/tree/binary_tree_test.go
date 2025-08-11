// Package tree provides a generic binary search tree implementation.
package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBinaryTree(t *testing.T) {
	t.Run("without root", func(t *testing.T) {
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)
		require.NotNil(t, tree)
		assert.Nil(t, tree.root)
		assert.Equal(t, 0, tree.Size())
	})

	t.Run("with root", func(t *testing.T) {
		rootNode := NewNode[uint64](123, "root")
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)
		tree.root = rootNode
		tree.size = 1

		require.NotNil(t, tree)
		assert.Equal(t, rootNode, tree.root)
		assert.Equal(t, 1, tree.Size())
	})
}

func TestGetHash(t *testing.T) {
	tree, err := NewBinaryTree[any]()
	require.NoError(t, err)

	testCases := []struct {
		name          string
		value1        any
		value2        any
		expectEqual   bool
		expectErr     bool
		expectedError error
	}{
		{"same strings", "hello", "hello", true, false, nil},
		{"different strings", "hello", "world", false, false, nil},
		{"same ints", 123, 123, true, false, nil},
		{"different ints", 123, 456, false, false, nil},
		{"same float64s", 123.456, 123.456, true, false, nil},
		{"different float64s", 123.456, 987.654, false, false, nil},
		{"unsupported type", true, false, false, true, ErrorUnsupportedValueType},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			hash1, err1 := tree.getHash(tc.value1)
			if tc.expectErr {
				assert.ErrorIs(t, err1, tc.expectedError)
				return
			}
			require.NoError(t, err1)

			hash2, err2 := tree.getHash(tc.value2)
			require.NoError(t, err2)

			if tc.expectEqual {
				assert.Equal(t, hash1, hash2)
			} else {
				assert.NotEqual(t, hash1, hash2)
			}
		})
	}
}

func TestBinaryTree_InsertAndFind(t *testing.T) {
	t.Run("insert first element", func(t *testing.T) {
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)
		err = tree.InsertInOrder("apple")
		require.NoError(t, err)
		assert.Equal(t, 1, tree.Size())
		require.NotNil(t, tree.root)
		assert.Equal(t, "apple", tree.root.value)
	})

	t.Run("insert multiple elements", func(t *testing.T) {
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)
		values := []string{"apple", "orange", "banana", "grape", "pineapple"}
		for i, v := range values {
			err = tree.InsertInOrder(v)
			require.NoError(t, err)
			assert.Equal(t, i+1, tree.Size())
		}

		// Find all inserted elements
		for _, v := range values {
			t.Run("find_"+v, func(t *testing.T) {
				var foundNode *Node[uint64, string]
				foundNode, err = tree.Find(v)
				assert.NoError(t, err)
				assert.NotNil(t, foundNode)
				if foundNode != nil {
					assert.Equal(t, v, foundNode.value)
				}
			})
		}

		// Find non-existent element
		var foundNode *Node[uint64, string]
		foundNode, err = tree.Find("watermelon")
		assert.NoError(t, err)
		assert.Nil(t, foundNode)
	})

	t.Run("find in empty tree", func(t *testing.T) {
		emptyTree, err := NewBinaryTree[string]()
		require.NoError(t, err)
		foundNode, err := emptyTree.Find("anything") // Ensure no panic on empty tree
		assert.Error(t, err)
		assert.Equal(t, ErrorNodeIsNil, err)
		assert.Nil(t, foundNode)
	})

	t.Run("insert unsupported type", func(t *testing.T) {
		errorTree, err := NewBinaryTree[any]()
		require.NoError(t, err)
		err = errorTree.InsertInOrder(true) // bool is unsupported
		assert.ErrorIs(t, err, ErrorUnsupportedValueType)
		assert.Equal(t, 0, errorTree.Size())
	})
}

func TestBinaryTree_Delete(t *testing.T) {
	t.Run("delete from empty tree", func(t *testing.T) {
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)
		_, err = tree.Delete("nonexistent")
		assert.ErrorIs(t, err, ErrorNodeIsNil)
		assert.Equal(t, 0, tree.Size())
	})

	t.Run("delete leaf node", func(t *testing.T) {
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)

		// Insert nodes
		values := []string{"apple", "banana", "cherry"}
		for _, v := range values {
			err = tree.InsertInOrder(v)
			require.NoError(t, err)
		}
		initialSize := tree.Size()

		// Delete a leaf node
		deletedKey, err := tree.Delete("cherry")
		require.NoError(t, err)
		assert.NotZero(t, deletedKey)
		assert.Equal(t, initialSize-1, tree.Size())

		// Verify it's deleted
		foundNode, err := tree.Find("cherry")
		assert.NoError(t, err)
		assert.Nil(t, foundNode)

		// Verify other nodes still exist
		foundNode, err = tree.Find("apple")
		assert.NoError(t, err)
		assert.NotNil(t, foundNode)
	})

	t.Run("delete node with one child", func(t *testing.T) {
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)

		// Insert nodes to create a node with one child
		values := []string{"banana", "apple", "cherry", "date"}
		for _, v := range values {
			err = tree.InsertInOrder(v)
			require.NoError(t, err)
		}
		initialSize := tree.Size()

		// Delete node with one child
		deletedKey, err := tree.Delete("cherry")
		require.NoError(t, err)
		assert.NotZero(t, deletedKey)
		assert.Equal(t, initialSize-1, tree.Size())

		// Verify it's deleted
		foundNode, err := tree.Find("cherry")
		assert.NoError(t, err)
		assert.Nil(t, foundNode)

		// Verify child node still exists
		foundNode, err = tree.Find("date")
		assert.NoError(t, err)
		assert.NotNil(t, foundNode)
	})

	t.Run("delete node with two children", func(t *testing.T) {
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)

		// Insert nodes to create a node with two children
		values := []string{"banana", "apple", "cherry", "blueberry", "date"}
		for _, v := range values {
			err = tree.InsertInOrder(v)
			require.NoError(t, err)
		}
		initialSize := tree.Size()

		// Delete node with two children
		deletedKey, err := tree.Delete("cherry")
		require.NoError(t, err)
		assert.NotZero(t, deletedKey)
		assert.Equal(t, initialSize-1, tree.Size())

		// Verify it's deleted
		foundNode, err := tree.Find("cherry")
		assert.NoError(t, err)
		assert.Nil(t, foundNode)

		// Verify children still exist
		foundNode, err = tree.Find("blueberry")
		assert.NoError(t, err)
		assert.NotNil(t, foundNode)

		foundNode, err = tree.Find("date")
		assert.NoError(t, err)
		assert.NotNil(t, foundNode)
	})

	t.Run("delete root node", func(t *testing.T) {
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)

		// Insert only root
		err = tree.InsertInOrder("root")
		require.NoError(t, err)
		assert.Equal(t, 1, tree.Size())

		// Delete root
		deletedKey, err := tree.Delete("root")
		require.NoError(t, err)
		assert.NotZero(t, deletedKey)
		assert.Equal(t, 0, tree.Size())
		assert.Nil(t, tree.root)
	})

	t.Run("delete nonexistent node", func(t *testing.T) {
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)

		// Insert some nodes
		values := []string{"apple", "banana", "cherry"}
		for _, v := range values {
			err = tree.InsertInOrder(v)
			require.NoError(t, err)
		}
		initialSize := tree.Size()

		// Try to delete nonexistent node
		deletedKey, err := tree.Delete("nonexistent")
		require.Error(t, err) // Should return ErrorNodeNotFound
		require.Equal(t, ErrorNodeNotFound, err)
		assert.Zero(t, deletedKey)                // Hash will be computed even if not found
		assert.Equal(t, initialSize, tree.Size()) // Size should remain unchanged

		// Verify existing nodes are still there
		foundNode, err := tree.Find("apple")
		assert.NoError(t, err)
		assert.NotNil(t, foundNode)
	})

	t.Run("delete with unsupported type", func(t *testing.T) {
		tree, err := NewBinaryTree[any]()
		require.NoError(t, err)

		// Insert a node first to avoid ErrorNodeIsNil
		err = tree.InsertInOrder("test")
		require.NoError(t, err)

		_, err = tree.Delete(true) // bool is unsupported
		assert.ErrorIs(t, err, ErrorUnsupportedValueType)
	})

	t.Run("delete all nodes", func(t *testing.T) {
		tree, err := NewBinaryTree[string]()
		require.NoError(t, err)

		// Insert multiple nodes
		values := []string{"banana", "apple", "cherry", "date", "elderberry"}
		for _, v := range values {
			err = tree.InsertInOrder(v)
			require.NoError(t, err)
		}

		// Delete all nodes one by one
		for i, v := range values {
			var deletedKey uint64
			deletedKey, err = tree.Delete(v)
			require.NoError(t, err)
			assert.NotZero(t, deletedKey)
			expectedSize := len(values) - i - 1
			assert.Equal(t, expectedSize, tree.Size())
		}

		// Tree should be empty
		assert.Equal(t, 0, tree.Size())
		assert.Nil(t, tree.root)
	})
}
