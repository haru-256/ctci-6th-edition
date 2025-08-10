package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInsertInOrder(t *testing.T) {
	t.Run("insert into single node tree", func(t *testing.T) {
		root := NewTreeNode(10)

		// Insert smaller
		err := root.InsertInOrder(5)
		require.NoError(t, err)
		assert.NotNil(t, root.left)
		assert.Equal(t, 5, root.left.data)
		assert.Equal(t, root, root.left.parent)
		assert.Nil(t, root.right)
		assert.Equal(t, 2, root.Size())

		// Insert larger
		err = root.InsertInOrder(15)
		require.NoError(t, err)
		assert.NotNil(t, root.right)
		assert.Equal(t, 15, root.right.data)
		assert.Equal(t, root, root.right.parent)
		assert.Equal(t, 3, root.Size())
	})

	t.Run("insert equal value", func(t *testing.T) {
		root := NewTreeNode(10)
		err := root.InsertInOrder(10)
		require.NoError(t, err)
		assert.NotNil(t, root.left)
		assert.Equal(t, 10, root.left.data)
		assert.Equal(t, 2, root.Size())
	})

	t.Run("insert into multi-level tree and check sizes", func(t *testing.T) {
		root := NewTreeNode(20)
		require.Equal(t, 1, root.Size())

		err := root.InsertInOrder(10)
		require.NoError(t, err)
		require.Equal(t, 2, root.Size())
		require.Equal(t, 1, root.left.Size())

		err = root.InsertInOrder(30)
		require.NoError(t, err)
		require.Equal(t, 3, root.Size())
		require.Equal(t, 1, root.right.Size())

		// Insert into left subtree
		err = root.InsertInOrder(5) // smaller than 10
		require.NoError(t, err)
		assert.NotNil(t, root.left.left)
		assert.Equal(t, 5, root.left.left.data)
		assert.Equal(t, root.left, root.left.left.parent)
		assert.Equal(t, 4, root.Size())
		assert.Equal(t, 2, root.left.Size())

		err = root.InsertInOrder(15) // larger than 10
		require.NoError(t, err)
		assert.NotNil(t, root.left.right)
		assert.Equal(t, 15, root.left.right.data)
		assert.Equal(t, root.left, root.left.right.parent)
		assert.Equal(t, 5, root.Size())
		assert.Equal(t, 3, root.left.Size())

		// Insert into right subtree
		err = root.InsertInOrder(25) // smaller than 30
		require.NoError(t, err)
		assert.NotNil(t, root.right.left)
		assert.Equal(t, 25, root.right.left.data)
		assert.Equal(t, root.right, root.right.left.parent)
		assert.Equal(t, 6, root.Size())
		assert.Equal(t, 2, root.right.Size())

		err = root.InsertInOrder(35) // larger than 30
		require.NoError(t, err)
		assert.NotNil(t, root.right.right)
		assert.Equal(t, 35, root.right.right.data)
		assert.Equal(t, root.right, root.right.right.parent)
		assert.Equal(t, 7, root.Size())
		assert.Equal(t, 3, root.right.Size())
	})
}
