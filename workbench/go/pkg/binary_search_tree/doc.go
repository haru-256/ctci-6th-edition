/*
Package binary_search_tree provides a generic binary search tree implementation with hash-based keys.

The binary search tree maintains elements in sorted order using hash values as keys, allowing
for efficient search, insertion, and deletion operations. The implementation uses FNV-1a hashing
for consistent key generation and supports any comparable value type.

# Features

- Generic implementation supporting any comparable value type
- Hash-based key generation using FNV-1a algorithm
- Thread-safe operations with read-write mutex protection
- Efficient O(log n) average-case performance for core operations
- Self-balancing through proper insertion order management
- Automatic memory management with garbage collection support

# Performance Characteristics

- Insert: O(log n) average, O(n) worst case
- Delete: O(log n) average, O(n) worst case
- Search: O(log n) average, O(n) worst case
- Space: O(n)

Note: Performance depends on hash distribution. Good hash functions provide balanced trees.

# Basic Usage

	// Create a new binary search tree for strings
	tree, err := binary_search_tree.NewBinaryTree[string]()
	if err != nil {
		log.Fatal(err)
	}

	// Insert values
	err = tree.InsertInOrder("apple")
	if err != nil {
		log.Fatal(err)
	}
	err = tree.InsertInOrder("banana")
	if err != nil {
		log.Fatal(err)
	}
	err = tree.InsertInOrder("cherry")
	if err != nil {
		log.Fatal(err)
	}

	// Search for values
	node, err := tree.Find("banana")
	if err != nil {
		log.Fatal(err)
	}
	if node != nil {
		fmt.Printf("Found: %s\n", node.Value())
	}

	// Delete a value
	deletedKey, err := tree.Delete("banana")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted key: %d\n", deletedKey)

# Advanced Usage

	// Create a tree for custom types
	type Product struct {
		ID   int
		Name string
	}

	productTree, err := binary_search_tree.NewBinaryTree[Product]()
	if err != nil {
		log.Fatal(err)
	}

	// Insert products
	products := []Product{
		{ID: 1, Name: "Laptop"},
		{ID: 2, Name: "Mouse"},
		{ID: 3, Name: "Keyboard"},
	}

	for _, product := range products {
		err := productTree.InsertInOrder(product)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Check tree size
	fmt.Printf("Tree contains %d products\n", productTree.Size())

	// Find a specific product
	target := Product{ID: 2, Name: "Mouse"}
	node, err := productTree.Find(target)
	if err != nil {
		log.Fatal(err)
	}
	if node != nil {
		fmt.Printf("Found product: %+v\n", node.Value())
	}

# Hash Function Details

The tree uses FNV-1a hashing for key generation, which provides:
- Fast computation with good distribution properties
- Deterministic results for consistent tree structure
- Support for int, float64, and string types
- Thread-safe hasher pool for concurrent access

Supported types for hashing:
- int: Converted to 8-byte representation
- float64: Uses IEEE 754 bit representation
- string: Direct byte conversion
- Other comparable types: May return ErrorUnsupportedValueType

# Concurrency

The binary search tree is thread-safe for all operations:

	var wg sync.WaitGroup
	tree, _ := binary_search_tree.NewBinaryTree[int]()

	// Safe concurrent insertions
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			tree.InsertInOrder(value)
		}(i)
	}
	wg.Wait()

	// Safe concurrent reads
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			node, err := tree.Find(value)
			if err == nil && node != nil {
				fmt.Printf("Found: %d\n", node.Value())
			}
		}(i)
	}
	wg.Wait()

# Error Handling

The package defines specific errors for different failure conditions:

- ErrorNodeIsNil: Returned when operating on nil nodes or empty trees
- ErrorNodeNotFound: Returned when deletion target doesn't exist
- ErrorUnsupportedValueType: Returned for unsupported hash types

Always check for errors when performing tree operations:

	node, err := tree.Find("nonexistent")
	if err != nil {
		if errors.Is(err, binary_search_tree.ErrorNodeIsNil) {
			fmt.Println("Tree is empty")
		}
		return
	}

# Memory Management

The implementation is designed for efficient memory usage:
- Nodes are allocated individually for optimal memory layout
- Deleted nodes are properly dereferenced for garbage collection
- Hash pool reuses hasher instances to reduce allocations
- Parent-child relationships are maintained bidirectionally

# Implementation Notes

The tree uses a hash-based approach where:
1. Values are hashed to generate 64-bit unsigned integer keys
2. Keys determine node placement following BST properties
3. Duplicate keys are handled by placing new nodes in left subtree
4. Deletion implements standard BST deletion with successor replacement

This approach provides consistent ordering based on hash values rather than
natural value ordering, which can be beneficial for maintaining balance
with certain data distributions.
*/
package binary_search_tree
