package heap

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

// intCmp compares two integers for max heap (larger values have higher priority)
func intCmp(a, b *int) int {
	if *a > *b {
		return 1
	} else if *a < *b {
		return -1
	}
	return 0
}

// stringCmp compares two strings lexicographically for max heap
func stringCmp(a, b *string) int {
	if *a > *b {
		return 1
	} else if *a < *b {
		return -1
	}
	return 0
}

// float64Cmp compares two float64 values for max heap
func float64Cmp(a, b *float64) int {
	if *a > *b {
		return 1
	} else if *a < *b {
		return -1
	}
	return 0
}

func TestNewHeap(t *testing.T) {
	heap := NewHeap[int](intCmp)

	require.NotNil(t, heap, "NewHeap should not return nil")
	assert.Equal(t, 0, heap.Size(), "Expected size 0")

	items := heap.GetItems()
	assert.Empty(t, items, "Expected empty items slice")
}

func TestNewMaxHeap_Convenience(t *testing.T) {
	heap := NewMaxHeap[int]()

	require.NotNil(t, heap, "NewMaxHeap should not return nil")
	assert.Equal(t, 0, heap.Size(), "Expected size 0")
}

func TestNewMinHeap_Convenience(t *testing.T) {
	heap := NewMinHeap[int]()

	require.NotNil(t, heap, "NewMinHeap should not return nil")
	assert.Equal(t, 0, heap.Size(), "Expected size 0")
}

func TestMinHeap_Operations(t *testing.T) {
	heap := NewMinHeap[int]()

	// Insert elements
	require.NoError(t, heap.Insert(30))
	require.NoError(t, heap.Insert(10))
	require.NoError(t, heap.Insert(20))
	require.NoError(t, heap.Insert(5))

	// Root should be minimum (5)
	min, err := heap.Peek()
	require.NoError(t, err, "Peek should not return error")
	assert.Equal(t, 5, *min, "Expected min 5")

	// Pop should return elements in ascending order
	expectedValues := []int{5, 10, 20, 30}
	for i, expectedValue := range expectedValues {
		var item *int
		item, err = heap.Pop()
		require.NoError(t, err, "Pop %d should not return error", i)
		assert.Equal(t, expectedValue, *item, "Pop %d: expected %d", i, expectedValue)
	}
}

func TestHeap_Insert(t *testing.T) {
	heap := NewHeap[int](intCmp)

	// Test inserting single element
	require.NoError(t, heap.Insert(10))
	assert.Equal(t, 1, heap.Size(), "Expected size 1")

	top, err := heap.Peek()
	require.NoError(t, err, "Peek should not return error")
	assert.Equal(t, 10, *top, "Expected 10")

	// Test inserting multiple elements
	require.NoError(t, heap.Insert(20))
	require.NoError(t, heap.Insert(5))
	require.NoError(t, heap.Insert(15))

	assert.Equal(t, 4, heap.Size(), "Expected size 4")

	// The root should be the maximum element for max heap
	top, err = heap.Peek()
	require.NoError(t, err, "Peek should not return error")
	assert.Equal(t, 20, *top, "Expected top 20")
}

func TestMaxHeap_HeapProperty(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Insert elements in random order
	values := []int{10, 30, 20, 40, 50, 15, 25}

	for _, v := range values {
		require.NoError(t, heap.Insert(v))
	}

	// Get items using GetItems() to verify heap property
	items := heap.GetItems()

	// Verify max heap property: parent >= children
	for i := 0; i < heap.Size()/2; i++ {
		left := Left(i)
		right := Right(i)

		if left < heap.Size() {
			assert.GreaterOrEqual(t, *items[i], *items[left],
				"Max heap property violated: parent[%d]=%d < left[%d]=%d",
				i, *items[i], left, *items[left])
		}

		if right < heap.Size() {
			assert.GreaterOrEqual(t, *items[i], *items[right],
				"Max heap property violated: parent[%d]=%d < right[%d]=%d",
				i, *items[i], right, *items[right])
		}
	}
}

func TestMaxHeap_Pop(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Test pop from empty heap
	_, err := heap.Pop()
	assert.Equal(t, ErrorIsEmpty, err, "Expected ErrorIsEmpty")

	// Insert elements and test pop
	require.NoError(t, heap.Insert(10))
	require.NoError(t, heap.Insert(20))
	require.NoError(t, heap.Insert(5))
	require.NoError(t, heap.Insert(15))

	// Pop should return elements in descending order
	expectedValues := []int{20, 15, 10, 5}
	for i, expectedValue := range expectedValues {
		var item *int
		item, err = heap.Pop()
		require.NoError(t, err, "Unexpected error at pop %d", i)
		assert.Equal(t, expectedValue, *item, "Pop %d: expected %d", i, expectedValue)
	}

	// Heap should be empty now
	assert.Equal(t, 0, heap.Size(), "Expected size 0 after all pops")

	// Another pop should return error
	_, err = heap.Pop()
	assert.Equal(t, ErrorIsEmpty, err, "Expected ErrorIsEmpty after popping empty heap")
}

func TestMaxHeap_Max(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Test max on empty heap
	_, err := heap.Peek()
	assert.Equal(t, ErrorIsEmpty, err, "Expected ErrorIsEmpty")

	// Insert elements and test max
	require.NoError(t, heap.Insert(10))
	max, err := heap.Peek()
	require.NoError(t, err, "Unexpected error")
	assert.Equal(t, 10, *max, "Expected max 10")

	require.NoError(t, heap.Insert(20))
	max, err = heap.Peek()
	require.NoError(t, err, "Unexpected error")
	assert.Equal(t, 20, *max, "Expected max 20")

	// Size should remain the same after Peek()
	assert.Equal(t, 2, heap.Size(), "Expected size 2 after Peek()")
}

func TestLeft(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 1},
		{1, 3},
		{2, 5},
		{3, 7},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, Left(test.input), "Left(%d) should be %d", test.input, test.expected)
	}
}

func TestRight(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{0, 2},
		{1, 4},
		{2, 6},
		{3, 8},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, Right(test.input), "Right(%d) should be %d", test.input, test.expected)
	}
}

func TestParent(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{1, 0},
		{2, 0},
		{3, 1},
		{4, 1},
		{5, 2},
		{6, 2},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, Parent(test.input), "Parent(%d) should be %d", test.input, test.expected)
	}
}

func TestBuildMaxHeap(t *testing.T) {
	// Create array with unordered elements
	items := []int{10, 30, 20, 40, 50}
	itemPtrs := make([]*int, len(items))
	for i, item := range items {
		itemPtrs[i] = &item
	}

	// Build max heap from array
	heap, err := BuildMaxHeap(itemPtrs)
	require.NoError(t, err)

	// Verify max heap property
	for i := 0; i < heap.Size()/2; i++ {
		left := Left(i)
		right := Right(i)

		heapItems := heap.GetItems()
		if left < heap.Size() {
			assert.GreaterOrEqual(t, *heapItems[i], *heapItems[left],
				"Max heap property violated after BuildMaxHeap: parent[%d]=%d < left[%d]=%d",
				i, *heapItems[i], left, *heapItems[left])
		}

		if right < heap.Size() {
			assert.GreaterOrEqual(t, *heapItems[i], *heapItems[right],
				"Max heap property violated after BuildMaxHeap: parent[%d]=%d < right[%d]=%d",
				i, *heapItems[i], right, *heapItems[right])
		}
	}

	// Root should be maximum
	max, err := heap.Peek()
	require.NoError(t, err, "Unexpected error")
	assert.Equal(t, 50, *max, "Expected root 50 after BuildMaxHeap")
}

func TestBuildMinHeap(t *testing.T) {
	// Create array with unordered elements
	items := []int{50, 20, 40, 10, 30}
	itemPtrs := make([]*int, len(items))
	for i, item := range items {
		itemPtrs[i] = &item
	}

	// Build min heap from array
	heap, err := BuildMinHeap(itemPtrs)
	require.NoError(t, err)

	// Verify min heap property
	for i := 0; i < heap.Size()/2; i++ {
		left := Left(i)
		right := Right(i)

		heapItems := heap.GetItems()
		if left < heap.Size() {
			assert.LessOrEqual(t, *heapItems[i], *heapItems[left],
				"Min heap property violated after BuildMinHeap: parent[%d]=%d > left[%d]=%d",
				i, *heapItems[i], left, *heapItems[left])
		}

		if right < heap.Size() {
			assert.LessOrEqual(t, *heapItems[i], *heapItems[right],
				"Min heap property violated after BuildMinHeap: parent[%d]=%d > right[%d]=%d",
				i, *heapItems[i], right, *heapItems[right])
		}
	}

	// Root should be minimum
	min, err := heap.Peek()
	require.NoError(t, err, "Unexpected error")
	assert.Equal(t, 10, *min, "Expected root 10 after BuildMinHeap")
}

func TestHeapSort(t *testing.T) {
	// Create array with elements in random order
	values := []int{30, 10, 50, 20, 40}
	itemPtrs := make([]*int, len(values))
	for i, v := range values {
		itemPtrs[i] = &v
	}

	// Perform heap sort
	sorted, err := HeapSort(itemPtrs)
	require.NoError(t, err)

	// After heap sort, the array should be sorted in ascending order
	expected := []int{10, 20, 30, 40, 50}
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], *sorted[i], "After HeapSort, items[%d] should be %d", i, expected[i])
	}

	// Verify size remains the same
	assert.Equal(t, len(values), len(sorted), "Expected size %d after HeapSort", len(values))
}

func TestHeapSort_LargerDataset(t *testing.T) {
	// Create array with elements in random order (larger dataset)
	values := []int{64, 34, 25, 12, 22, 11, 90, 5, 77, 30, 55, 45, 60, 78, 1}
	itemPtrs := make([]*int, len(values))
	for i, v := range values {
		itemPtrs[i] = &v
	}

	// Perform heap sort
	sorted, err := HeapSort(itemPtrs)
	require.NoError(t, err)

	// After heap sort, the array should be sorted in ascending order
	expected := []int{1, 5, 11, 12, 22, 25, 30, 34, 45, 55, 60, 64, 77, 78, 90}
	for i := 0; i < len(expected); i++ {
		assert.Equal(t, expected[i], *sorted[i], "After HeapSort, items[%d] should be %d", i, expected[i])
	}
}

func TestHeapSort_EmptyHeap(t *testing.T) {
	// Test heap sort on empty array
	var itemPtrs []*int

	sorted, err := HeapSort(itemPtrs)
	require.NoError(t, err)

	assert.Empty(t, sorted, "Expected empty array after HeapSort")
}

func TestHeapSort_SingleElement(t *testing.T) {
	// Test heap sort on single element array
	value := 42
	itemPtrs := []*int{&value}

	sorted, err := HeapSort(itemPtrs)
	require.NoError(t, err)

	assert.Len(t, sorted, 1, "Expected size 1 for single element array after HeapSort")
	assert.Equal(t, 42, *sorted[0], "Expected 42")
}

func TestMaxHeap_WithDifferentTypes(t *testing.T) {
	// Test with string values
	stringHeap := NewHeap[string](stringCmp)
	require.NoError(t, stringHeap.Insert("zebra"))
	require.NoError(t, stringHeap.Insert("apple"))
	require.NoError(t, stringHeap.Insert("banana"))

	// Root should be "zebra" (lexicographically largest)
	max, err := stringHeap.Peek()
	require.NoError(t, err, "Unexpected error")
	assert.Equal(t, "zebra", *max, "Expected max string 'zebra'")

	// Test with float64 values
	floatHeap := NewHeap[float64](float64Cmp)
	require.NoError(t, floatHeap.Insert(3.14))
	require.NoError(t, floatHeap.Insert(2.71))
	require.NoError(t, floatHeap.Insert(1.41))

	// Root should be 3.14 (largest)
	max2, err := floatHeap.Peek()
	require.NoError(t, err, "Unexpected error")
	assert.Equal(t, 3.14, *max2, "Expected max 3.14")
}

func TestMaxHeap_IntegrationTest(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Insert many elements
	elements := []int{50, 30, 70, 20, 40, 60, 80}

	for _, elem := range elements {
		require.NoError(t, heap.Insert(elem))
	}

	// Verify size
	assert.Equal(t, len(elements), heap.Size(), "Expected size %d", len(elements))

	// Pop all elements and verify they come out in descending order
	var poppedValues []int
	for heap.Size() > 0 {
		item, err := heap.Pop()
		require.NoError(t, err, "Unexpected error during pop")
		poppedValues = append(poppedValues, *item)
	}

	// Verify order is descending
	for i := 1; i < len(poppedValues); i++ {
		assert.GreaterOrEqual(t, poppedValues[i-1], poppedValues[i],
			"Elements not popped in descending order: %v", poppedValues)
	}
}

// Test with custom type
type Person struct {
	Name string
	Age  int
}

func personCmpByAge(a, b *Person) int {
	if a.Age > b.Age {
		return 1
	} else if a.Age < b.Age {
		return -1
	}
	return 0
}

func TestMaxHeap_CustomType(t *testing.T) {
	heap := NewHeap[Person](personCmpByAge)

	// Insert people
	require.NoError(t, heap.Insert(Person{Name: "Alice", Age: 30}))
	require.NoError(t, heap.Insert(Person{Name: "Bob", Age: 25}))
	require.NoError(t, heap.Insert(Person{Name: "Charlie", Age: 35}))

	// Max should be Charlie (oldest)
	max, err := heap.Peek()
	require.NoError(t, err, "Unexpected error")
	assert.Equal(t, "Charlie", max.Name, "Expected Charlie")
	assert.Equal(t, 35, max.Age, "Expected age 35")

	// Pop all and verify order
	expectedNames := []string{"Charlie", "Alice", "Bob"}
	var person *Person
	for i, expectedName := range expectedNames {
		person, err = heap.Pop()
		require.NoError(t, err, "Unexpected error at pop %d", i)
		assert.Equal(t, expectedName, person.Name, "Pop %d: expected %s", i, expectedName)
	}
}

// Test with Node type for backwards compatibility
func TestMaxHeap_WithNodeType(t *testing.T) {
	// Define comparison function for Node
	nodeCmp := func(a, b *Node[int, string]) int {
		if a.Key > b.Key {
			return 1
		} else if a.Key < b.Key {
			return -1
		}
		return 0
	}

	heap := NewHeap[Node[int, string]](nodeCmp)

	// Insert nodes
	require.NoError(t, heap.Insert(Node[int, string]{Key: 10, Value: "ten"}))
	require.NoError(t, heap.Insert(Node[int, string]{Key: 20, Value: "twenty"}))
	require.NoError(t, heap.Insert(Node[int, string]{Key: 5, Value: "five"}))

	// Max should be node with key 20
	max, err := heap.Peek()
	require.NoError(t, err, "Unexpected error")
	assert.Equal(t, 20, max.Key, "Expected key 20")
	assert.Equal(t, "twenty", max.Value, "Expected value 'twenty'")
}

// Benchmark tests
func BenchmarkMaxHeap_Insert(b *testing.B) {
	heap := NewMaxHeap[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = heap.Insert(i)
	}
}

func BenchmarkMaxHeap_Pop(b *testing.B) {
	// Pre-populate heap
	heap := NewMaxHeap[int]()
	for i := 0; i < b.N; i++ {
		_ = heap.Insert(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := heap.Pop()
		require.NoError(b, err, "Pop should not return an error")
	}
}

func BenchmarkHeapSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		// Create array with 1000 elements
		itemPtrs := make([]*int, 1000)
		for j := 1000; j > 0; j-- {
			value := j
			itemPtrs[1000-j] = &value
		}

		b.StartTimer()
		_, _ = HeapSort(itemPtrs)
	}
}

// TestHeap_ConcurrentInsertPop tests thread safety of Insert and Pop operations
// by running multiple goroutines concurrently using errgroup
func TestHeap_ConcurrentInsertPop(t *testing.T) {
	heap := NewMaxHeap[int]()

	const numGoroutines = 10
	const numOpsPerGoroutine = 100

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	// Launch goroutines that insert values
	for i := 0; i < numGoroutines; i++ {
		goroutineID := i
		g.Go(func() error {
			for j := 0; j < numOpsPerGoroutine; j++ {
				value := goroutineID*numOpsPerGoroutine + j
				_ = heap.Insert(value)
			}
			return nil
		})
	}

	// Launch goroutines that pop values
	for i := 0; i < numGoroutines; i++ {
		g.Go(func() error {
			for j := 0; j < numOpsPerGoroutine; j++ {
				// Try to pop, but handle empty heap gracefully
				_, _ = heap.Pop() // Ignore errors as heap might be empty
			}
			return nil
		})
	}

	// Wait for all operations to complete
	err := g.Wait()
	require.NoError(t, err, "All goroutines should complete without error")

	// Verify no deadlocks occurred by checking heap is still functional
	require.NoError(t, heap.Insert(999))
	top, err := heap.Peek()
	require.NoError(t, err, "Heap should still be functional after concurrent operations")
	require.NotNil(t, top, "Should be able to peek after concurrent operations")
	assert.Equal(t, 999, *top, "Inserted value should be at the top of max heap")
}

// TestHeap_ConcurrentReads tests that multiple read operations can run concurrently
// without blocking each other (RWMutex behavior) using errgroup
func TestHeap_ConcurrentReads(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Populate heap with test data
	testData := []int{50, 30, 70, 20, 40, 60, 80}
	for _, val := range testData {
		require.NoError(t, heap.Insert(val))
	}

	const numReaders = 20
	const numReadsPerReader = 50

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	// Launch concurrent Peek operations
	for i := 0; i < numReaders; i++ {
		g.Go(func() error {
			for j := 0; j < numReadsPerReader; j++ {
				top, err := heap.Peek()
				require.NoError(t, err, "Peek should not fail")
				require.NotNil(t, top, "Peek should not return nil")
				assert.Equal(t, 80, *top, "Max element should be 80")
			}
			return nil
		})
	}

	// Launch concurrent Size operations
	for i := 0; i < numReaders; i++ {
		g.Go(func() error {
			for j := 0; j < numReadsPerReader; j++ {
				size := heap.Size()
				assert.Equal(t, len(testData), size, "Size should be consistent")
				assert.Greater(t, size, 0, "Size should be positive")
			}
			return nil
		})
	}

	// Launch concurrent GetItems operations
	for i := 0; i < numReaders; i++ {
		g.Go(func() error {
			for j := 0; j < numReadsPerReader; j++ {
				items := heap.GetItems()
				require.NotNil(t, items, "GetItems should not return nil")
				assert.Equal(t, len(testData), len(items), "Items length should match")
			}
			return nil
		})
	}

	// Wait for all readers to complete
	err := g.Wait()
	require.NoError(t, err, "All read operations should complete without error")

	// Verify heap is still functional and correct
	top, err := heap.Peek()
	require.NoError(t, err, "Final peek should work")
	assert.Equal(t, 80, *top, "Max element should still be 80")
	assert.Equal(t, len(testData), heap.Size(), "Size should be unchanged")
}

// TestHeap_ConcurrentUpDownHeap tests thread safety of UpHeap and DownHeap public methods
// These methods should be safe to call externally (e.g., from priority queue package)
func TestHeap_ConcurrentUpDownHeap(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Populate heap with some initial data to have valid indices
	initialData := []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100}
	for _, val := range initialData {
		require.NoError(t, heap.Insert(val))
	}

	const numGoroutines = 10
	const numOpsPerGoroutine = 20

	ctx := context.Background()
	g, _ := errgroup.WithContext(ctx)

	// Start concurrent operations
	startConcurrentUpHeap(g, heap, numGoroutines, numOpsPerGoroutine)
	startConcurrentDownHeap(g, heap, numGoroutines, numOpsPerGoroutine)
	startMixedOperations(g, heap, numGoroutines, numOpsPerGoroutine)

	// Wait for all operations to complete
	err := g.Wait()
	require.NoError(t, err, "All concurrent operations should complete without error")

	// Verify heap functionality after concurrent operations
	verifyHeapIntegrityAfterConcurrency(t, heap)
}

// startConcurrentUpHeap launches goroutines that perform UpHeap operations
func startConcurrentUpHeap(g *errgroup.Group, heap *Heap[int], numGoroutines, numOpsPerGoroutine int) {
	for i := 0; i < numGoroutines; i++ {
		g.Go(func() error {
			return performUpHeapOperations(heap, numOpsPerGoroutine)
		})
	}
}

// performUpHeapOperations executes UpHeap operations for a single goroutine
func performUpHeapOperations(heap *Heap[int], numOps int) error {
	for j := 0; j < numOps; j++ {
		// Use valid indices (heap size may change during concurrent operations)
		size := heap.Size()
		if size > 0 {
			index := j % size
			if err := heap.UpHeap(index); err != nil {
				return err
			}
		}
	}
	return nil
}

// startConcurrentDownHeap launches goroutines that perform DownHeap operations
func startConcurrentDownHeap(g *errgroup.Group, heap *Heap[int], numGoroutines, numOpsPerGoroutine int) {
	for i := 0; i < numGoroutines; i++ {
		g.Go(func() error {
			return performDownHeapOperations(heap, numOpsPerGoroutine)
		})
	}
}

// performDownHeapOperations executes DownHeap operations for a single goroutine
func performDownHeapOperations(heap *Heap[int], numOps int) error {
	for j := 0; j < numOps; j++ {
		// Use valid indices (heap size may change during concurrent operations)
		size := heap.Size()
		if size > 0 {
			index := j % size
			if err := heap.DownHeap(index); err != nil {
				return err
			}
		}
	}
	return nil
}

// startMixedOperations launches goroutines that perform mixed Insert/Peek operations
func startMixedOperations(g *errgroup.Group, heap *Heap[int], numGoroutines, numOpsPerGoroutine int) {
	for i := 0; i < numGoroutines; i++ {
		goroutineID := i
		g.Go(func() error {
			return performMixedOperations(heap, goroutineID, numOpsPerGoroutine)
		})
	}
}

// performMixedOperations executes mixed Insert/Peek operations for a single goroutine
func performMixedOperations(heap *Heap[int], goroutineID, numOps int) error {
	for j := 0; j < numOps; j++ {
		if j%2 == 0 {
			// Insert operation
			value := goroutineID*1000 + j
			if err := heap.Insert(value); err != nil {
				return err
			}
		} else {
			// Peek operation (ignore errors for empty heap)
if _, err := heap.Peek(); err != nil && err != ErrorIsEmpty {
	return err
}
		}
	}
	return nil
}

// verifyHeapIntegrityAfterConcurrency checks that the heap maintains its properties after concurrent operations
func verifyHeapIntegrityAfterConcurrency(t *testing.T, heap *Heap[int]) {
	// Verify heap is still functional and maintains heap property
	// Insert a known maximum value
	require.NoError(t, heap.Insert(9999))
	top, err := heap.Peek()
	require.NoError(t, err, "Heap should still be functional after concurrent UpHeap/DownHeap")
	require.NotNil(t, top, "Peek should return a value")
	assert.Equal(t, 9999, *top, "Heap property should be maintained after concurrent operations")

	// Verify we can still perform basic operations
	assert.Greater(t, heap.Size(), 0, "Heap should not be empty")

	// Pop the maximum and verify heap still works
	poppedMax, err := heap.Pop()
	require.NoError(t, err, "Should be able to pop after concurrent operations")
	require.NotNil(t, poppedMax, "Popped value should not be nil")
	assert.Equal(t, 9999, *poppedMax, "Should pop the maximum value")
}
