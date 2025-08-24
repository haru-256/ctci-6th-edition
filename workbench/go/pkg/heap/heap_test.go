package heap

import (
	"testing"

	"github.com/stretchr/testify/require"
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

	if heap == nil {
		t.Fatal("NewHeap returned nil")
	}

	if heap.Size() != 0 {
		t.Errorf("Expected size 0, got %d", heap.Size())
	}

	if len(heap.items) != 0 {
		t.Errorf("Expected empty items slice, got length %d", len(heap.items))
	}
}

func TestNewMaxHeap_Convenience(t *testing.T) {
	heap := NewMaxHeap[int]()

	if heap == nil {
		t.Fatal("NewMaxHeap returned nil")
	}

	if heap.Size() != 0 {
		t.Errorf("Expected size 0, got %d", heap.Size())
	}
}

func TestNewMinHeap_Convenience(t *testing.T) {
	heap := NewMinHeap[int]()

	if heap == nil {
		t.Fatal("NewMinHeap returned nil")
	}

	if heap.Size() != 0 {
		t.Errorf("Expected size 0, got %d", heap.Size())
	}
}

func TestHeap_Insert(t *testing.T) {
	heap := NewHeap[int](intCmp)

	// Test inserting single element
	heap.Insert(10)
	if heap.Size() != 1 {
		t.Errorf("Expected size 1, got %d", heap.Size())
	}

	top, err := heap.Peek()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if *top != 10 {
		t.Errorf("Expected 10, got %d", *top)
	}

	// Test inserting multiple elements
	heap.Insert(20)
	heap.Insert(5)
	heap.Insert(15)

	if heap.Size() != 4 {
		t.Errorf("Expected size 4, got %d", heap.Size())
	}

	// The root should be the maximum element for max heap
	top, err = heap.Peek()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if *top != 20 {
		t.Errorf("Expected top 20, got %d", *top)
	}
}

func TestMaxHeap_HeapProperty(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Insert elements in random order
	values := []int{10, 30, 20, 40, 50, 15, 25}

	for _, v := range values {
		heap.Insert(v)
	}

	// Verify max heap property: parent >= children
	for i := 0; i < heap.Size()/2; i++ {
		left := Left(i)
		right := Right(i)

		if left < heap.Size() && *heap.items[i] < *heap.items[left] {
			t.Errorf("Max heap property violated: parent[%d]=%d < left[%d]=%d",
				i, *heap.items[i], left, *heap.items[left])
		}

		if right < heap.Size() && *heap.items[i] < *heap.items[right] {
			t.Errorf("Max heap property violated: parent[%d]=%d < right[%d]=%d",
				i, *heap.items[i], right, *heap.items[right])
		}
	}
}

func TestMaxHeap_Pop(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Test pop from empty heap
	_, err := heap.Pop()
	if err != ErrorIsEmpty {
		t.Errorf("Expected ErrorIsEmpty, got %v", err)
	}

	// Insert elements and test pop
	heap.Insert(10)
	heap.Insert(20)
	heap.Insert(5)
	heap.Insert(15)

	// Pop should return elements in descending order
	expectedValues := []int{20, 15, 10, 5}
	for i, expectedValue := range expectedValues {
		var item *int
		item, err = heap.Pop()
		if err != nil {
			t.Fatalf("Unexpected error at pop %d: %v", i, err)
		}
		if *item != expectedValue {
			t.Errorf("Pop %d: expected %d, got %d", i, expectedValue, *item)
		}
	}

	// Heap should be empty now
	if heap.Size() != 0 {
		t.Errorf("Expected size 0 after all pops, got %d", heap.Size())
	}

	// Another pop should return error
	_, err = heap.Pop()
	if err != ErrorIsEmpty {
		t.Errorf("Expected ErrorIsEmpty after popping empty heap, got %v", err)
	}
}

func TestMaxHeap_Max(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Test max on empty heap
	_, err := heap.Peek()
	if err != ErrorIsEmpty {
		t.Errorf("Expected ErrorIsEmpty, got %v", err)
	}

	// Insert elements and test max
	heap.Insert(10)
	max, err := heap.Peek()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if *max != 10 {
		t.Errorf("Expected max 10, got %d", *max)
	}

	heap.Insert(20)
	max, err = heap.Peek()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if *max != 20 {
		t.Errorf("Expected max 20, got %d", *max)
	}

	// Size should remain the same after Peek()
	if heap.Size() != 2 {
		t.Errorf("Expected size 2 after Peek(), got %d", heap.Size())
	}
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
		result := Left(test.input)
		if result != test.expected {
			t.Errorf("Left(%d) = %d; expected %d", test.input, result, test.expected)
		}
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
		result := Right(test.input)
		if result != test.expected {
			t.Errorf("Right(%d) = %d; expected %d", test.input, result, test.expected)
		}
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
		result := Parent(test.input)
		if result != test.expected {
			t.Errorf("Parent(%d) = %d; expected %d", test.input, result, test.expected)
		}
	}
}

func TestBuildMaxHeap(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Manually add elements without maintaining heap property
	items := []int{10, 30, 20, 40, 50}
	for _, item := range items {
		heap.items = append(heap.items, &item)
	}

	// Build max heap
	BuildHeap(heap)

	// Verify max heap property
	for i := 0; i < heap.Size()/2; i++ {
		left := Left(i)
		right := Right(i)

		if left < heap.Size() && *heap.items[i] < *heap.items[left] {
			t.Errorf("Max heap property violated after BuildHeap: parent[%d]=%d < left[%d]=%d",
				i, *heap.items[i], left, *heap.items[left])
		}

		if right < heap.Size() && *heap.items[i] < *heap.items[right] {
			t.Errorf("Max heap property violated after BuildHeap: parent[%d]=%d < right[%d]=%d",
				i, *heap.items[i], right, *heap.items[right])
		}
	}

	// Root should be maximum
	if *heap.items[0] != 50 {
		t.Errorf("Expected root 50 after BuildHeap, got %d", *heap.items[0])
	}
}

func TestHeapSort(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Insert elements in random order
	values := []int{30, 10, 50, 20, 40}
	for _, v := range values {
		heap.Insert(v)
	}

	// Perform heap sort
	HeapSort(heap)

	// After heap sort, the array should be sorted in ascending order
	expected := []int{10, 20, 30, 40, 50}
	for i := 0; i < len(expected); i++ {
		if *heap.items[i] != expected[i] {
			t.Errorf("After HeapSort, items[%d] = %d; expected %d",
				i, *heap.items[i], expected[i])
		}
	}

	// Verify size remains the same
	if heap.Size() != len(values) {
		t.Errorf("Expected size %d after HeapSort, got %d", len(values), heap.Size())
	}
}

func TestHeapSort_LargerDataset(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Insert elements in random order (larger dataset)
	values := []int{64, 34, 25, 12, 22, 11, 90, 5, 77, 30, 55, 45, 60, 78, 1}
	for _, v := range values {
		heap.Insert(v)
	}

	// Perform heap sort
	HeapSort(heap)

	// After heap sort, the array should be sorted in ascending order
	expected := []int{1, 5, 11, 12, 22, 25, 30, 34, 45, 55, 60, 64, 77, 78, 90}
	for i := 0; i < len(expected); i++ {
		if *heap.items[i] != expected[i] {
			t.Errorf("After HeapSort, items[%d] = %d; expected %d",
				i, *heap.items[i], expected[i])
		}
	}
}

func TestHeapSort_EmptyHeap(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Test heap sort on empty heap
	HeapSort(heap)

	if heap.Size() != 0 {
		t.Errorf("Expected size 0 for empty heap after HeapSort, got %d", heap.Size())
	}
}

func TestHeapSort_SingleElement(t *testing.T) {
	heap := NewMaxHeap[int]()
	heap.Insert(42)

	HeapSort(heap)

	if heap.Size() != 1 {
		t.Errorf("Expected size 1 for single element heap after HeapSort, got %d", heap.Size())
	}

	if *heap.items[0] != 42 {
		t.Errorf("Expected 42, got %d", *heap.items[0])
	}
}

func TestMaxHeap_WithDifferentTypes(t *testing.T) {
	// Test with string values
	stringHeap := NewHeap[string](stringCmp)
	stringHeap.Insert("zebra")
	stringHeap.Insert("apple")
	stringHeap.Insert("banana")

	// Root should be "zebra" (lexicographically largest)
	max, err := stringHeap.Peek()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if *max != "zebra" {
		t.Errorf("Expected max 'zebra', got %s", *max)
	}

	// Test with float64 values
	floatHeap := NewHeap[float64](float64Cmp)
	floatHeap.Insert(3.14)
	floatHeap.Insert(2.71)
	floatHeap.Insert(1.41)

	// Root should be 3.14 (largest)
	max2, err := floatHeap.Peek()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if *max2 != 3.14 {
		t.Errorf("Expected max 3.14, got %f", *max2)
	}
}

func TestMaxHeap_IntegrationTest(t *testing.T) {
	heap := NewMaxHeap[int]()

	// Insert many elements
	elements := []int{50, 30, 70, 20, 40, 60, 80}

	for _, elem := range elements {
		heap.Insert(elem)
	}

	// Verify size
	if heap.Size() != len(elements) {
		t.Errorf("Expected size %d, got %d", len(elements), heap.Size())
	}

	// Pop all elements and verify they come out in descending order
	var poppedValues []int
	for heap.Size() > 0 {
		item, err := heap.Pop()
		if err != nil {
			t.Fatalf("Unexpected error during pop: %v", err)
		}
		poppedValues = append(poppedValues, *item)
	}

	// Verify order is descending
	for i := 1; i < len(poppedValues); i++ {
		if poppedValues[i-1] < poppedValues[i] {
			t.Errorf("Elements not popped in descending order: %v", poppedValues)
			break
		}
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
	heap.Insert(Person{Name: "Alice", Age: 30})
	heap.Insert(Person{Name: "Bob", Age: 25})
	heap.Insert(Person{Name: "Charlie", Age: 35})

	// Max should be Charlie (oldest)
	max, err := heap.Peek()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if max.Name != "Charlie" || max.Age != 35 {
		t.Errorf("Expected Charlie(35), got %s(%d)", max.Name, max.Age)
	}

	// Pop all and verify order
	expectedNames := []string{"Charlie", "Alice", "Bob"}
	for i, expectedName := range expectedNames {
		var person *Person
		person, err = heap.Pop()
		if err != nil {
			t.Fatalf("Unexpected error at pop %d: %v", i, err)
		}
		if person.Name != expectedName {
			t.Errorf("Pop %d: expected %s, got %s", i, expectedName, person.Name)
		}
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
	heap.Insert(Node[int, string]{Key: 10, Value: "ten"})
	heap.Insert(Node[int, string]{Key: 20, Value: "twenty"})
	heap.Insert(Node[int, string]{Key: 5, Value: "five"})

	// Max should be node with key 20
	max, err := heap.Peek()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if max.Key != 20 || max.Value != "twenty" {
		t.Errorf("Expected {20, twenty}, got {%d, %s}", max.Key, max.Value)
	}
}

// Benchmark tests
func BenchmarkMaxHeap_Insert(b *testing.B) {
	heap := NewMaxHeap[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i)
	}
}

func BenchmarkMaxHeap_Pop(b *testing.B) {
	// Pre-populate heap
	heap := NewMaxHeap[int]()
	for i := 0; i < b.N; i++ {
		heap.Insert(i)
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
		heap := NewMaxHeap[int]()

		// Insert 1000 elements
		for j := 1000; j > 0; j-- {
			heap.Insert(j)
		}

		b.StartTimer()
		HeapSort(heap)
	}
}
