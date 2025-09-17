/*
Package linked_list provides a generic doubly linked list implementation with comprehensive operations.

The doubly linked list maintains bidirectional links between nodes, allowing efficient insertion
and deletion at any position. It supports any comparable type and provides both head and tail
pointers for optimal performance at both ends of the list.

# Features

- Generic implementation supporting any comparable type
- Doubly linked structure with bidirectional navigation
- Efficient O(1) insertion and deletion at known positions
- Head and tail pointers for optimal end operations
- Linear search functionality with O(n) complexity
- Proper memory management with garbage collection support
- Thread-safe for concurrent use by multiple goroutines
- Simple and intuitive API design

# Thread Safety

This implementation is thread-safe and can be used concurrently by multiple goroutines.
All public methods use appropriate mutex locking:
- Read operations (Search) use RWMutex.RLock() for concurrent reads
- Write operations (Prepend, Insert, Delete) use RWMutex.Lock() for exclusive access
- The mutex prevents race conditions and ensures list consistency across goroutines

No external synchronization is required when using this linked list from multiple goroutines.

# Performance Characteristics

- Prepend: O(1)
- Insert after known node: O(1)
- Delete known node: O(1)
- Search: O(n)
- Delete by value: O(n) due to search phase
- Space: O(n)

The list excels at scenarios where frequent insertion/deletion is needed with known node references.

# Basic Usage

	// Create a new doubly linked list for integers
	list := linked_list.NewLinkedList[int]()

	// Add elements to the beginning
	list.Prepend(30)
	list.Prepend(20)
	list.Prepend(10)
	// List now contains: 10 <-> 20 <-> 30

	// Search for elements
	node := list.Search(20)
	if node != nil {
		fmt.Printf("Found: %d\n", node.Value)
	}

	// Insert after a specific node
	if node != nil {
		err := list.Insert(25, node)
		if err != nil {
			log.Fatal(err)
		}
	}
	// List now contains: 10 <-> 20 <-> 25 <-> 30

	// Delete by value
	err := list.Delete(20)
	if err != nil {
		log.Fatal(err)
	}
	// List now contains: 10 <-> 25 <-> 30

# Concurrent Usage

The linked list is thread-safe and can be used safely from multiple goroutines
without external synchronization:

	// Multiple goroutines can safely operate on the same list
	go func() {
		for i := 0; i < 100; i++ {
			list.Prepend(i)
		}
	}()

	go func() {
		for i := 0; i < 50; i++ {
			if node := list.Search(i); node != nil {
				fmt.Printf("Found: %d\n", node.Value)
			}
		}
	}()

	go func() {
		for i := 0; i < 25; i++ {
			list.Delete(i)
		}
	}()

# Advanced Usage with Custom Types

	// Define a custom type
	type Person struct {
		ID   int
		Name string
	}

	// Create list for custom type
	people := linked_list.NewLinkedList[Person]()

	// Add people to the list
	people.Prepend(Person{3, "Charlie"})
	people.Prepend(Person{2, "Bob"})
	people.Prepend(Person{1, "Alice"})

	// Search for specific person
	target := Person{2, "Bob"}
	bobNode := people.Search(target)
	if bobNode != nil {
		fmt.Printf("Found: %+v\n", bobNode.Value)

		// Insert a new person after Bob
		newPerson := Person{4, "David"}
		err := people.Insert(newPerson, bobNode)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Remove a person
	err := people.Delete(Person{1, "Alice"})
	if err != nil {
		log.Fatal(err)
	}

# Traversal Patterns

	list := linked_list.NewLinkedList[string]()

	// Populate list
	words := []string{"first", "second", "third", "fourth"}
	for _, word := range words {
		list.Prepend(word)
	}

	// Forward traversal from head
	fmt.Println("Forward traversal:")
	current := list.Head
	for current != nil {
		fmt.Printf("%s ", current.Value)
		current = current.Next
	}
	fmt.Println()

	// Backward traversal from tail
	fmt.Println("Backward traversal:")
	current = list.Tail
	for current != nil {
		fmt.Printf("%s ", current.Value)
		current = current.Prev
	}
	fmt.Println()

# Node Operations

	list := linked_list.NewLinkedList[int]()
	list.Prepend(1)
	list.Prepend(2)
	list.Prepend(3)

	// Get reference to specific nodes
	head := list.Head        // Node with value 3
	tail := list.Tail        // Node with value 1
	middle := head.Next      // Node with value 2

	// Insert using node references (O(1) operation)
	newNode := linked_list.NewNode(5)
	err := list.Insert(5, middle)  // Insert 5 after middle node
	if err != nil {
		log.Fatal(err)
	}

	// Verify the insertion
	if middle.Next != nil && middle.Next.Value == 5 {
		fmt.Println("Successfully inserted 5 after middle node")
	}

# Building Complex Data Structures

	// Use linked list as building block for other structures
	type Stack[T comparable] struct {
		list *linked_list.LinkedList[T]
	}

	func NewStack[T comparable]() *Stack[T] {
		return &Stack[T]{
			list: linked_list.NewLinkedList[T](),
		}
	}

	func (s *Stack[T]) Push(value T) {
		s.list.Prepend(value)  // Add to head for O(1) push
	}

	func (s *Stack[T]) Pop() (T, error) {
		var zero T
		if s.list.Head == nil {
			return zero, errors.New("stack is empty")
		}
		value := s.list.Head.Value
		err := s.list.Delete(value)
		return value, err
	}

	func (s *Stack[T]) IsEmpty() bool {
		return s.list.Head == nil
	}

	// Usage
	stack := NewStack[int]()
	stack.Push(1)
	stack.Push(2)
	stack.Push(3)

	for !stack.IsEmpty() {
		value, err := stack.Pop()
		if err != nil {
			break
		}
		fmt.Printf("Popped: %d\n", value)
	}

# Memory Management and Performance

	list := linked_list.NewLinkedList[int]()

	// Insert many elements
	for i := 0; i < 10000; i++ {
		list.Prepend(i)
	}

	// Delete elements (nodes are properly cleaned up)
	for i := 0; i < 5000; i++ {
		err := list.Delete(i)
		if err != nil {
			// Handle error - element might not exist
			continue
		}
	}

	// Remaining elements are still accessible
	remaining := 0
	current := list.Head
	for current != nil {
		remaining++
		current = current.Next
	}
	fmt.Printf("Remaining elements: %d\n", remaining)

Key memory management features:
- Deleted nodes have references cleared for garbage collection
- No memory leaks with proper reference management
- Efficient for scenarios with frequent insertions/deletions
- O(1) space overhead per element (two pointers + value)

# Error Handling

The linked list operations can return errors in specific conditions:

	list := linked_list.NewLinkedList[string]()

	// Attempting to delete from empty list
	err := list.Delete("nonexistent")
	if errors.Is(err, linked_list.ErrorNodeNotFound) {
		fmt.Println("Cannot delete from empty list")
	}

	// Attempting to insert after nil node
	err = list.Insert("value", nil)
	if errors.Is(err, linked_list.ErrorNodeIsNil) {
		fmt.Println("Cannot insert after nil node")
	}

	// Proper error handling
	node := list.Search("target")
	if node != nil {
		err := list.Insert("new_value", node)
		if err != nil {
			log.Printf("Insert failed: %v", err)
			return
		}
		fmt.Println("Insert successful")
	}

# Thread Safety

The linked list is not thread-safe. For concurrent access, use external synchronization:

	list := linked_list.NewLinkedList[int]()
	var mu sync.RWMutex

	// Safe concurrent reads
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			mu.RLock()
			node := list.Search(value)
			mu.RUnlock()
			if node != nil {
				fmt.Printf("Found: %d\n", node.Value)
			}
		}(i)
	}

	// Safe concurrent writes
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(value int) {
			defer wg.Done()
			mu.Lock()
			list.Prepend(value)
			mu.Unlock()
		}(i)
	}
	wg.Wait()

# Comparison with Other Data Structures

Choose linked lists when:
- Frequent insertion/deletion at arbitrary positions
- Unknown or highly variable size
- Memory usage is more important than cache locality
- Need for bidirectional traversal

Consider alternatives when:
- Random access by index is needed (use slices/arrays)
- Cache performance is critical (use contiguous arrays)
- Memory overhead is a concern (linked lists have pointer overhead)
- Sorting is frequently needed (use specialized data structures)

# Implementation Details

The doubly linked list implementation features:
- Each node contains Value, Next, and Prev pointers
- Head and Tail pointers maintained for O(1) end access
- Proper cleanup of references during deletion
- Generic type support with comparable constraint
- Simple node creation with NewNode function
- Bidirectional navigation for flexible traversal patterns

The structure is optimized for:
- Fast insertion/deletion with known node references
- Efficient traversal in both directions
- Minimal memory overhead per operation
- Clean integration with Go's garbage collector
*/
package linked_list
