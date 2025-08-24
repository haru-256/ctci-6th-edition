/*
Package hashtable provides a thread-safe hash table implementation using chaining for collision resolution.

The hash table uses separate chaining with linked lists to handle hash collisions, providing
reliable performance even under high load factors. It supports any comparable type and uses
FNV-1a hashing for consistent key distribution.

# Features

- Generic implementation supporting any comparable type
- Separate chaining collision resolution using doubly linked lists
- Thread-safe operations with read-write mutex protection
- FNV-1a hashing algorithm for consistent key distribution
- Configurable table size for optimal performance tuning
- Automatic memory management with garbage collection support

# Performance Characteristics

- Insert: O(1) average, O(n) worst case (with many collisions)
- Search: O(1) average, O(n) worst case (with many collisions)
- Delete: O(1) average, O(n) worst case (with many collisions)
- Space: O(n + m) where n is number of elements and m is table size

Load factor affects performance. For optimal results, keep load factor below 0.75.

# Basic Usage

	// Create a hash table with 10 buckets
	table := hashtable.NewHashChainTable[string](10)

	// Insert values
	err := table.Insert("apple")
	if err != nil {
		log.Fatal(err)
	}
	err = table.Insert("banana")
	if err != nil {
		log.Fatal(err)
	}
	err = table.Insert("cherry")
	if err != nil {
		log.Fatal(err)
	}

	// Search for values
	node, err := table.Search("banana")
	if err != nil {
		log.Fatal(err)
	}
	if node != nil {
		fmt.Printf("Found: %s\n", node.Value)
	}

	// Check table size
	fmt.Printf("Table contains %d elements\n", table.Size())

	// Delete a value
	err = table.Delete("banana")
	if err != nil {
		log.Fatal(err)
	}

# Advanced Usage

	// Create a hash table for custom types
	type User struct {
		ID   int
		Name string
	}

	userTable := hashtable.NewHashChainTable[User](50)

	// Insert users
	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}

	for _, user := range users {
		err := userTable.Insert(user)
		if err != nil {
			if errors.Is(err, hashtable.ErrorAlreadyExists) {
				fmt.Printf("User %+v already exists\n", user)
				continue
			}
			log.Fatal(err)
		}
	}

	// Search for specific user
	target := User{ID: 2, Name: "Bob"}
	node, err := userTable.Search(target)
	if err != nil {
		log.Fatal(err)
	}
	if node != nil {
		fmt.Printf("Found user: %+v\n", node.Value)
	}

# Collision Resolution

The hash table uses separate chaining with doubly linked lists:

	// Small table size to demonstrate collision handling
	smallTable := hashtable.NewHashChainTable[int](3)

	// Insert values that may collide
	values := []int{1, 4, 7, 10, 13} // These may hash to same buckets
	for _, value := range values {
		err := smallTable.Insert(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	// All values are stored and searchable despite collisions
	for _, value := range values {
		node, err := smallTable.Search(value)
		if err != nil {
			log.Fatal(err)
		}
		if node != nil {
			fmt.Printf("Found: %d\n", node.Value)
		}
	}

# Load Factor Management

For optimal performance, monitor and manage load factor:

	table := hashtable.NewHashChainTable[string](100)
	
	// Insert many elements
	for i := 0; i < 1000; i++ {
		value := fmt.Sprintf("item_%d", i)
		err := table.Insert(value)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Calculate load factor
	loadFactor := float64(table.Size()) / float64(table.MaxSize)
	fmt.Printf("Load factor: %.2f\n", loadFactor)

	if loadFactor > 0.75 {
		fmt.Println("Consider increasing table size for better performance")
	}

# Concurrency

The hash table is thread-safe for all operations:

	table := hashtable.NewHashChainTable[int](100)
	var wg sync.WaitGroup

	// Safe concurrent insertions
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			err := table.Insert(value)
			if err != nil && !errors.Is(err, hashtable.ErrorAlreadyExists) {
				log.Printf("Insert error: %v", err)
			}
		}(i)
	}
	wg.Wait()

	// Safe concurrent searches
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			node, err := table.Search(value)
			if err != nil {
				log.Printf("Search error: %v", err)
				return
			}
			if node != nil {
				fmt.Printf("Found: %d\n", node.Value)
			}
		}(i)
	}
	wg.Wait()

# Hash Function Details

The hash table uses FNV-1a hashing for key generation:

	// Supported types
	table.Insert(42)        // int
	table.Insert(3.14)      // float64  
	table.Insert("hello")   // string

	// Unsupported types will return ErrorUnsupportedValueType
	err := table.Insert(true) // bool not supported
	if errors.Is(err, hashtable.ErrorUnsupportedValueType) {
		fmt.Println("Boolean type not supported for hashing")
	}

The FNV-1a algorithm provides:
- Fast computation with good distribution
- Deterministic results for consistent behavior
- Low collision rate for typical data patterns
- Thread-safe hasher pool for concurrent access

# Error Handling

The package defines specific errors for different failure conditions:

- ErrorAlreadyExists: Returned when inserting duplicate values
- ErrorUnsupportedValueType: Returned for unsupported hash types
- ErrorNodeNotFound: Returned when deleting non-existent values

Handle errors appropriately in your application:

	err := table.Insert("duplicate")
	if err != nil {
		switch {
		case errors.Is(err, hashtable.ErrorAlreadyExists):
			fmt.Println("Value already exists in table")
		case errors.Is(err, hashtable.ErrorUnsupportedValueType):
			fmt.Println("Value type not supported for hashing")
		default:
			log.Fatal(err)
		}
	}

# Memory Management

The implementation efficiently manages memory:
- Empty buckets are set to nil to allow garbage collection
- Linked lists are deallocated when buckets become empty
- Hash pool reuses hasher instances to reduce allocations
- Deleted nodes have references cleared for garbage collection

# Performance Tips

1. Choose appropriate table size based on expected load:
   - Table size should be roughly equal to expected number of elements
   - Use prime numbers for table size to improve hash distribution

2. Monitor load factor and resize when necessary:
   - Keep load factor below 0.75 for optimal performance
   - Consider dynamic resizing for applications with varying load

3. Use appropriate types:
   - Prefer int, float64, and string for best hash performance
   - Ensure custom types implement proper equality comparison

4. Handle collisions gracefully:
   - Expect some collisions with any hash function
   - Performance degrades gracefully with separate chaining
   - Consider alternative data structures for very high collision rates
*/
package hashtable
