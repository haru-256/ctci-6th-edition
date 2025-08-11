// Package tree provides a generic binary search tree implementation.
package tree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNode(t *testing.T) {
	t.Run("NewNode", func(t *testing.T) {
		node := NewNode(10, "apple")
		require.NotNil(t, node)
		assert.Equal(t, 10, node.key)
		assert.Equal(t, "apple", node.value)
		assert.Nil(t, node.parent)
		assert.Nil(t, node.left)
		assert.Nil(t, node.right)
	})

	t.Run("SetLeftChild", func(t *testing.T) {
		parent := NewNode(10, "parent")
		child := NewNode(5, "child")

		err := parent.setLeftChild(child)
		require.NoError(t, err)

		assert.Equal(t, child, parent.left)
		assert.Equal(t, parent, child.parent)
	})

	t.Run("SetRightChild", func(t *testing.T) {
		parent := NewNode(10, "parent")
		child := NewNode(15, "child")

		err := parent.setRightChild(child)
		require.NoError(t, err)

		assert.Equal(t, child, parent.right)
		assert.Equal(t, parent, child.parent)
	})

	t.Run("SetChild_ErrorCases", func(t *testing.T) {
		parent := NewNode(10, "parent")
		child := NewNode(5, "child")

		err := (*Node[int, string])(nil).setLeftChild(child)
		assert.ErrorIs(t, err, ErrorNodeIsNil, "setLeftChild on nil parent should fail")

		err = parent.setLeftChild(nil)
		assert.ErrorIs(t, err, ErrorNodeIsNil, "setLeftChild with nil child should fail")

		err = (*Node[int, string])(nil).setRightChild(child)
		assert.ErrorIs(t, err, ErrorNodeIsNil, "setRightChild on nil parent should fail")

		err = parent.setRightChild(nil)
		assert.ErrorIs(t, err, ErrorNodeIsNil, "setRightChild with nil child should fail")
	})
}

func TestInsertInOrder_Node(t *testing.T) {
	type insertAction struct {
		key   int
		value string
	}
	type validation func(t *testing.T, root *Node[int, string])

	testCases := []struct {
		name       string
		rootNode   *Node[int, string]
		inserts    []insertAction
		validation validation
	}{
		{
			name:     "insert smaller and larger",
			rootNode: NewNode(20, "root"),
			inserts: []insertAction{
				{10, "ten"},
				{30, "thirty"},
			},
			validation: func(t *testing.T, root *Node[int, string]) {
				if assert.NotNil(t, root.left) {
					assert.Equal(t, 10, root.left.key)
					assert.Equal(t, "ten", root.left.value)
					assert.Equal(t, root, root.left.parent)
				}
				if assert.NotNil(t, root.right) {
					assert.Equal(t, 30, root.right.key)
					assert.Equal(t, "thirty", root.right.value)
					assert.Equal(t, root, root.right.parent)
				}
			},
		},
		{
			name:     "insert duplicate key",
			rootNode: NewNode(20, "root"),
			inserts: []insertAction{
				{10, "ten"},
				{20, "root_duplicate"},
			},
			validation: func(t *testing.T, root *Node[int, string]) {
				// 20 <= 20 -> left
				// 20 > 10 -> right
				if assert.NotNil(t, root.left) && assert.NotNil(t, root.left.right) {
					assert.Equal(t, 20, root.left.right.key)
					assert.Equal(t, "root_duplicate", root.left.right.value)
					assert.Equal(t, root.left, root.left.right.parent)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for _, insert := range tc.inserts {
				err := tc.rootNode.insertInOrder(insert.key, insert.value)
				require.NoError(t, err)
			}
			tc.validation(t, tc.rootNode)
		})
	}

	t.Run("error on nil node", func(t *testing.T) {
		err := (*Node[int, string])(nil).insertInOrder(1, "one")
		assert.ErrorIs(t, err, ErrorNodeIsNil)
	})
}

func TestFind_Node(t *testing.T) {
	// Setup tree
	root := NewNode(20, "twenty")
	require.NoError(t, root.insertInOrder(10, "ten"))
	require.NoError(t, root.insertInOrder(30, "thirty"))
	require.NoError(t, root.insertInOrder(5, "five"))
	require.NoError(t, root.insertInOrder(15, "fifteen"))

	testCases := []struct {
		name        string
		keyToFind   int
		valueToFind string
		expectFound bool
		expectedKey int
		expectedVal string
	}{
		{"find root", 20, "twenty", true, 20, "twenty"},
		{"find middle node", 10, "ten", true, 10, "ten"},
		{"find leaf node", 5, "five", true, 5, "five"},
		{"find another leaf", 15, "fifteen", true, 15, "fifteen"},
		{"find non-existent node", 99, "ninety-nine", false, 0, ""},
		{"find existing key, wrong value", 10, "not-ten", false, 0, ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			foundNode := root.find(tc.keyToFind, tc.valueToFind)
			if tc.expectFound {
				if assert.NotNil(t, foundNode) {
					assert.Equal(t, tc.expectedKey, foundNode.key)
					assert.Equal(t, tc.expectedVal, foundNode.value)
				}
			} else {
				assert.Nil(t, foundNode)
			}
		})
	}

	t.Run("find on nil node", func(t *testing.T) {
		foundNode := (*Node[int, string])(nil).find(1, "one")
		assert.Nil(t, foundNode)
	})
}
