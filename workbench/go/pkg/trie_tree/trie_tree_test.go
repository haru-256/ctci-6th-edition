package trietree

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestTrieTree_NewTrieTree(t *testing.T) {
	trie := NewTrieTree[byte, string]()
	require.NotNil(t, trie, "NewTrieTree should return a valid trie")
	require.NotNil(t, trie.root, "NewTrieTree should create a valid root node")
	assert.False(t, trie.root.isEnd, "Root node should not be marked as end")
	assert.Equal(t, 0, len(trie.root.children), "Root node should have no children initially")
}

func TestTrieTree_Insert_and_Search(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Test basic insertion and search
	key1 := []byte("hello")
	value1 := "world"
	trie.Insert(key1, value1)

	result, found := trie.Search(key1)
	assert.True(t, found, "Should find inserted key")
	assert.Equal(t, value1, result, "Should return correct value")

	// Test searching non-existent key
	key2 := []byte("goodbye")
	_, found = trie.Search(key2)
	assert.False(t, found, "Should not find non-existent key")

	// Test inserting multiple keys
	key3 := []byte("help")
	value3 := "assistance"
	trie.Insert(key3, value3)

	result, found = trie.Search(key3)
	assert.True(t, found, "Should find second inserted key")
	assert.Equal(t, value3, result, "Should return correct value for second key")

	// Original key should still be there
	result, found = trie.Search(key1)
	assert.True(t, found, "Original key should still exist")
	assert.Equal(t, value1, result, "Original key should have correct value")
}

func TestTrieTree_Insert_OverwriteValue(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	key := []byte("test")
	value1 := "first"
	value2 := "second"

	trie.Insert(key, value1)
	trie.Insert(key, value2) // Overwrite

	result, found := trie.Search(key)
	assert.True(t, found, "Should find the key")
	assert.Equal(t, value2, result, "Should return overwritten value")
}

func TestTrieTree_Insert_EmptyKey(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	key := []byte{}
	value := "empty"
	trie.Insert(key, value)

	result, found := trie.Search(key)
	assert.True(t, found, "Should find empty key")
	assert.Equal(t, value, result, "Should return correct value for empty key")
}

func TestTrieTree_StartsWith(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Insert some keys
	trie.Insert([]byte("hello"), "world")
	trie.Insert([]byte("help"), "assistance")
	trie.Insert([]byte("helicopter"), "aircraft")

	tests := []struct {
		prefix   []byte
		expected bool
	}{
		{[]byte("hel"), true},
		{[]byte("hello"), true},
		{[]byte("help"), true},
		{[]byte("helicopter"), true},
		{[]byte("helicopters"), false},
		{[]byte("world"), false},
		{[]byte(""), true}, // Empty prefix should match
		{[]byte("xyz"), false},
	}

	for _, test := range tests {
		result := trie.StartsWith(test.prefix)
		assert.Equal(t, test.expected, result, "StartsWith(%s) should return %v", test.prefix, test.expected)
	}
}

func TestTrieTree_Delete(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Insert test data
	trie.Insert([]byte("hello"), "world")
	trie.Insert([]byte("help"), "assistance")
	trie.Insert([]byte("helicopter"), "aircraft")

	// Test deleting non-existent key
	err := trie.Delete([]byte("nonexistent"))
	assert.Error(t, err, "Should return error when deleting non-existent key")
	assert.Equal(t, ErrKeyNotFound, err, "Should return ErrKeyNotFound")

	// Test deleting existing key
	err = trie.Delete([]byte("help"))
	assert.NoError(t, err, "Should not return error when deleting existing key")

	// Verify key is deleted
	_, found := trie.Search([]byte("help"))
	assert.False(t, found, "Deleted key should not be found")

	// Verify other keys still exist
	_, found = trie.Search([]byte("hello"))
	assert.True(t, found, "Other keys should still exist after deletion")

	_, found = trie.Search([]byte("helicopter"))
	assert.True(t, found, "Other keys should still exist after deletion")

	// Verify prefix still works
	assert.True(t, trie.StartsWith([]byte("hel")), "Prefix should still work after deletion")
}

func TestTrieTree_Delete_EmptyKey(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	err := trie.Delete([]byte{})
	assert.Error(t, err, "Should return error when deleting empty key")
	assert.Equal(t, ErrKeyNotFound, err, "Should return ErrKeyNotFound for empty key")
}

func TestTrieTree_Delete_PrefixKey(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Insert keys where one is prefix of another
	trie.Insert([]byte("test"), "value1")
	trie.Insert([]byte("testing"), "value2")

	// Delete the shorter key
	err := trie.Delete([]byte("test"))
	assert.NoError(t, err, "Should be able to delete prefix key")

	// Verify shorter key is deleted
	_, found := trie.Search([]byte("test"))
	assert.False(t, found, "Deleted prefix key should not be found")

	// Verify longer key still exists
	result, found := trie.Search([]byte("testing"))
	assert.True(t, found, "Longer key should still exist after deleting prefix")
	assert.Equal(t, "value2", result, "Longer key should have correct value")
}

func TestTrieTree_Size(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Empty trie
	assert.Equal(t, 0, trie.Size(), "Empty trie size should be 0")

	// Add keys one by one
	keys := [][]byte{
		[]byte("hello"),
		[]byte("help"),
		[]byte("helicopter"),
		[]byte("world"),
	}

	for i, key := range keys {
		trie.Insert(key, "value")
		expectedSize := i + 1
		assert.Equal(t, expectedSize, trie.Size(), "Size should be %d after inserting %d keys", expectedSize, expectedSize)
	}

	// Delete a key
	err := trie.Delete([]byte("help"))
	if err != nil {
		t.Errorf("Failed to delete key: %v", err)
	}
	assert.Equal(t, len(keys)-1, trie.Size(), "Size should be %d after deleting 1 key", len(keys)-1)

	// Insert duplicate (should not increase size)
	trie.Insert([]byte("hello"), "new value")
	assert.Equal(t, len(keys)-1, trie.Size(), "Size should remain %d after overwriting existing key", len(keys)-1)
}

func TestTrieTree_IsEmpty(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Empty trie
	assert.True(t, trie.IsEmpty(), "New trie should be empty")

	// Add a key
	trie.Insert([]byte("test"), "value")
	assert.False(t, trie.IsEmpty(), "Trie with keys should not be empty")

	// Delete the key
	err := trie.Delete([]byte("test"))
	if err != nil {
		t.Errorf("Failed to delete key: %v", err)
	}
	assert.True(t, trie.IsEmpty(), "Trie should be empty after deleting all keys")
}

func TestTrieTree_Keys(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Empty trie
	keys := trie.Keys()
	assert.Empty(t, keys, "Empty trie should return empty keys slice")

	// Add some keys
	expectedKeys := [][]byte{
		[]byte("hello"),
		[]byte("help"),
		[]byte("helicopter"),
		[]byte("world"),
	}

	for _, key := range expectedKeys {
		trie.Insert(key, "value")
	}

	keys = trie.Keys()
	assert.Len(t, keys, len(expectedKeys), "Should return correct number of keys")

	// Convert to strings for easier comparison
	var resultStrings []string
	var expectedStrings []string

	for _, key := range keys {
		resultStrings = append(resultStrings, string(key))
	}
	for _, key := range expectedKeys {
		expectedStrings = append(expectedStrings, string(key))
	}

	sort.Strings(resultStrings)
	sort.Strings(expectedStrings)

	assert.Equal(t, expectedStrings, resultStrings, "Should return all inserted keys")
}

func TestTrieTree_KeysWithPrefix(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Insert test data
	testData := map[string]string{
		"hello":      "world",
		"help":       "assistance",
		"helicopter": "aircraft",
		"world":      "earth",
		"wonder":     "amazing",
	}

	for k, v := range testData {
		trie.Insert([]byte(k), v)
	}

	tests := []struct {
		prefix   []byte
		expected []string
	}{
		{[]byte("hel"), []string{"hello", "help", "helicopter"}},
		{[]byte("hello"), []string{"hello"}},
		{[]byte("w"), []string{"world", "wonder"}},
		{[]byte("wor"), []string{"world"}},
		{[]byte("xyz"), []string{}}, // No matches
	}

	for _, test := range tests {
		keys, err := trie.KeysWithPrefix(test.prefix)

		if len(test.expected) == 0 {
			// Expecting no matches
			assert.Error(t, err, "KeysWithPrefix(%s) should return error for non-existent prefix", test.prefix)
			continue
		}

		assert.NoError(t, err, "KeysWithPrefix(%s) should not return error", test.prefix)
		assert.Len(t, keys, len(test.expected), "KeysWithPrefix(%s) should return %d keys", test.prefix, len(test.expected))

		// Convert to strings for comparison
		var resultStrings []string
		for _, key := range keys {
			resultStrings = append(resultStrings, string(key))
		}

		sort.Strings(resultStrings)
		sort.Strings(test.expected)

		assert.Equal(t, test.expected, resultStrings, "KeysWithPrefix(%s) should return correct keys", test.prefix)
	}
}

func TestTrieTree_DifferentTypes(t *testing.T) {
	// Test with int keys
	intTrie := NewTrieTree[int, string]()
	intKey := []int{1, 2, 3}
	intTrie.Insert(intKey, "integer key")

	result, found := intTrie.Search(intKey)
	assert.True(t, found, "Should find integer key")
	assert.Equal(t, "integer key", result, "Should return correct value for integer key")

	// Test with string keys
	stringTrie := NewTrieTree[string, int]()
	stringKey := []string{"hello", "world"}
	stringTrie.Insert(stringKey, 42)

	intResult, found := stringTrie.Search(stringKey)
	assert.True(t, found, "Should find string key")
	assert.Equal(t, 42, intResult, "Should return correct integer value")
}

func TestTrieTree_KeysWithPrefix_EmptyPrefixReturnsAll(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	data := map[string]string{
		"alpha":    "a",
		"beta":     "b",
		"gamma":    "g",
		"alphabet": "ab",
	}

	for k, v := range data {
		trie.Insert([]byte(k), v)
	}

	keys, err := trie.KeysWithPrefix([]byte(""))
	require.NoError(t, err)
	require.Len(t, keys, len(data))

	got := make([]string, 0, len(keys))
	for _, k := range keys {
		got = append(got, string(k))
	}
	want := make([]string, 0, len(data))
	for k := range data {
		want = append(want, k)
	}
	sort.Strings(got)
	sort.Strings(want)
	assert.Equal(t, want, got)
}

func TestTrieTree_ConcurrentAccess(t *testing.T) {
	t.Parallel()

	trie := NewTrieTree[byte, string]()

	keys := []string{
		"a", "ab", "abc", "abcd", "abcde",
		"b", "bc", "bcd", "bcde",
		"c", "cd", "cde",
	}

	// Concurrent inserts
	var g errgroup.Group
	for _, k := range keys {
		k := k // capture
		g.Go(func() error {
			trie.Insert([]byte(k), k)
			return nil
		})
	}
	require.NoError(t, g.Wait(), "concurrent inserts should not fail")

	// Verify all keys were inserted
	for _, k := range keys {
		v, ok := trie.Search([]byte(k))
		require.True(t, ok, "key %q should exist after concurrent insert", k)
		require.Equal(t, k, v, "key %q should have correct value", k)
	}

	// Concurrent reads
	for _, k := range keys {
		k := k
		g.Go(func() error {
			v, ok := trie.Search([]byte(k))
			if !ok {
				t.Errorf("search failed for %q: key not found", k)
				return nil // don't fail errgroup, just log error
			}
			if v != k {
				t.Errorf("search failed for %q: expected %q, got %q", k, k, v)
			}
			return nil
		})
	}
	require.NoError(t, g.Wait(), "concurrent reads should not fail")

	// Mixed deletes and reads
	for i, k := range keys {
		k := k
		if i%2 == 0 {
			g.Go(func() error {
				_ = trie.Delete([]byte(k)) // ignore error; concurrent delete may happen twice
				return nil
			})
		} else {
			g.Go(func() error {
				trie.Search([]byte(k))
				return nil
			})
		}
	}
	require.NoError(t, g.Wait(), "mixed concurrent operations should not fail")

	// Finally, ensure at least half the keys remain (odd indices)
	remaining := 0
	for i, k := range keys {
		if i%2 == 1 {
			if _, ok := trie.Search([]byte(k)); ok {
				remaining++
			}
		}
	}
	assert.GreaterOrEqual(t, remaining, len(keys)/2, "at least half the keys should remain after concurrent deletes")
}

func TestTrieTree_ConcurrentStress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping stress test in short mode")
	}

	trie := NewTrieTree[byte, string]()
	var g errgroup.Group
	const numGoroutines = 10
	const keysPerGoroutine = 50

	// Multiple inserters with different key patterns
	for i := 0; i < numGoroutines; i++ {
		i := i
		g.Go(func() error {
			for j := 0; j < keysPerGoroutine; j++ {
				key := []byte(fmt.Sprintf("key-%d-%d", i, j))
				value := fmt.Sprintf("value-%d-%d", i, j)
				trie.Insert(key, value)
			}
			return nil
		})
	}

	// Multiple readers
	for i := 0; i < numGoroutines; i++ {
		i := i
		g.Go(func() error {
			for j := 0; j < keysPerGoroutine; j++ {
				key := []byte(fmt.Sprintf("key-%d-%d", i, j))
				trie.Search(key) // ignore result, just ensure no panic
			}
			return nil
		})
	}

	// Multiple deleters
	for i := 0; i < numGoroutines/2; i++ {
		i := i
		g.Go(func() error {
			for j := 0; j < keysPerGoroutine/2; j++ {
				key := []byte(fmt.Sprintf("key-%d-%d", i, j))
				_ = trie.Delete(key) // ignore error, might not exist
			}
			return nil
		})
	}

	require.NoError(t, g.Wait(), "stress test should not fail")

	// Verify trie is still functional
	testKey := []byte("test-after-stress")
	testValue := "test-value"
	trie.Insert(testKey, testValue)

	result, found := trie.Search(testKey)
	assert.True(t, found, "trie should be functional after stress test")
	assert.Equal(t, testValue, result, "trie should return correct value after stress test")
}

// Benchmark tests
func BenchmarkTrieTree_Insert(b *testing.B) {
	trie := NewTrieTree[byte, string]()
	key := []byte("benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Insert(key, "value")
	}
}

func BenchmarkTrieTree_Search(b *testing.B) {
	trie := NewTrieTree[byte, string]()
	key := []byte("benchmark")
	trie.Insert(key, "value")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Search(key)
	}
}

func BenchmarkTrieTree_Delete(b *testing.B) {
	trie := NewTrieTree[byte, string]()
	key := []byte("benchmark")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		trie.Insert(key, "value")
		b.StartTimer()
		err := trie.Delete(key)
		require.NoError(b, err, "Delete should not fail in benchmark")
	}
}

// Additional edge case tests
func TestTrieTree_EdgeCases(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Test with nil-like behavior (empty trie operations)
	keys := trie.Keys()
	assert.Empty(t, keys, "Empty trie should return empty keys")

	// Test deleting from empty trie
	err := trie.Delete([]byte("nonexistent"))
	assert.Equal(t, ErrKeyNotFound, err, "Deleting from empty trie should return ErrKeyNotFound")

	// Test StartsWith on empty trie
	assert.True(t, trie.StartsWith([]byte{}), "Empty prefix should always return true")
	assert.False(t, trie.StartsWith([]byte("test")), "Non-empty prefix on empty trie should return false")

	// Test KeysWithPrefix on empty trie
	_, err = trie.KeysWithPrefix([]byte("test"))
	assert.Equal(t, ErrKeyNotFound, err, "KeysWithPrefix on empty trie should return error")
}

func TestTrieTree_LargeKey(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Test with a very large key
	largeKey := make([]byte, 1000)
	for i := range largeKey {
		largeKey[i] = byte(i % 256)
	}

	trie.Insert(largeKey, "large")
	result, found := trie.Search(largeKey)
	assert.True(t, found, "Should find large key")
	assert.Equal(t, "large", result, "Should return correct value for large key")

	// Delete large key
	err := trie.Delete(largeKey)
	assert.NoError(t, err, "Should be able to delete large key")

	_, found = trie.Search(largeKey)
	assert.False(t, found, "Large key should be deleted")
}

func TestTrieTree_SpecialValues(t *testing.T) {
	trie := NewTrieTree[byte, *string]()

	// Test with nil pointer value
	key := []byte("test")
	var nilValue *string = nil
	trie.Insert(key, nilValue)

	result, found := trie.Search(key)
	assert.True(t, found, "Should find key with nil value")
	assert.Nil(t, result, "Should preserve nil value")

	// Test with zero-value string
	stringTrie := NewTrieTree[byte, string]()
	stringTrie.Insert(key, "")

	stringResult, found := stringTrie.Search(key)
	assert.True(t, found, "Should find key with empty string value")
	assert.Equal(t, "", stringResult, "Should preserve empty string value")
}

func TestTrieTree_ConcurrentPrefixes(t *testing.T) {
	trie := NewTrieTree[byte, string]()

	// Insert keys with shared prefixes
	trie.Insert([]byte("a"), "1")
	trie.Insert([]byte("ab"), "2")
	trie.Insert([]byte("abc"), "3")
	trie.Insert([]byte("abcd"), "4")
	trie.Insert([]byte("abcde"), "5")

	// Test that all keys exist
	expectedValues := map[string]string{
		"a":     "1",
		"ab":    "2",
		"abc":   "3",
		"abcd":  "4",
		"abcde": "5",
	}

	for keyStr, expectedValue := range expectedValues {
		result, found := trie.Search([]byte(keyStr))
		assert.True(t, found, "Key %s should be found", keyStr)
		assert.Equal(t, expectedValue, result, "Key %s should have correct value", keyStr)
	}

	// Test size
	assert.Equal(t, 5, trie.Size(), "Should have 5 keys")

	// Delete middle key and verify others remain
	err := trie.Delete([]byte("abc"))
	assert.NoError(t, err, "Should be able to delete middle key")

	// Verify deleted key is gone
	_, found := trie.Search([]byte("abc"))
	assert.False(t, found, "Deleted key should not be found")

	// Verify other keys still exist
	for keyStr, expectedValue := range expectedValues {
		if keyStr == "abc" {
			continue // Skip the deleted key
		}
		result, exists := trie.Search([]byte(keyStr))
		assert.True(t, exists, "Key %s should still exist after deleting abc", keyStr)
		assert.Equal(t, expectedValue, result, "Key %s should have correct value", keyStr)
	}

	assert.Equal(t, 4, trie.Size(), "Should have 4 keys after deletion")
}
