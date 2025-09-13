// Package trietree provides a generic trie (prefix tree) data structure implementation.
//
// A trie is a tree-like data structure that is used to store a dynamic set of strings
// or sequences, where each node represents a common prefix of some strings. This implementation
// is generic and can work with any comparable key type and any value type.
//
// Example usage with byte slices (strings):
//
//	trie := trietree.NewTrieTree[byte, string]()
//	trie.Insert([]byte("hello"), "world")
//	trie.Insert([]byte("help"), "assistance")
//
//	if value, found := trie.Search([]byte("hello")); found {
//		fmt.Println(value) // Output: world
//	}
//
//	if trie.StartsWith([]byte("hel")) {
//		fmt.Println("Found prefix") // This will be printed
//	}
//
// Example usage with integer sequences:
//
//	trie := trietree.NewTrieTree[int, string]()
//	trie.Insert([]int{1, 2, 3}, "sequence")
//	value, found := trie.Search([]int{1, 2, 3})
//
// Time Complexities:
//   - Insert: O(m) where m is the length of the key
//   - Search: O(m) where m is the length of the key
//   - Delete: O(m) where m is the length of the key
//   - StartsWith: O(m) where m is the length of the prefix
//   - Size: O(n) where n is the total number of nodes in the trie
//   - Keys: O(n*m) where n is the number of keys and m is the average key length
//   - KeysWithPrefix: O(k*m) where k is the number of matching keys and m is the average key length
//
// Space Complexity: O(ALPHABET_SIZE * N * M) where ALPHABET_SIZE is the number of possible
// key elements, N is the number of keys, and M is the average length of the keys.
package trietree
