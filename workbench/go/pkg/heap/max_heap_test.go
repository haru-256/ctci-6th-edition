package heap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewMaxHeap(t *testing.T) {
	heap := NewMaxHeap[int, string]()

	if heap == nil {
		t.Fatal("NewMaxHeap returned nil")
	}

	if heap.Size() != 0 {
		t.Errorf("Expected size 0, got %d", heap.Size())
	}

	if len(heap.values) != 0 {
		t.Errorf("Expected empty values slice, got length %d", len(heap.values))
	}
}

func TestMaxHeap_Insert(t *testing.T) {
	heap := NewMaxHeap[int, string]()

	// Test inserting single element
	heap.Insert(10, "ten")
	if heap.Size() != 1 {
		t.Errorf("Expected size 1, got %d", heap.Size())
	}

	max, err := heap.Max()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if max.Key != 10 || max.Value != "ten" {
		t.Errorf("Expected {10, ten}, got {%v, %v}", max.Key, max.Value)
	}

	// Test inserting multiple elements
	heap.Insert(20, "twenty")
	heap.Insert(5, "five")
	heap.Insert(15, "fifteen")

	if heap.Size() != 4 {
		t.Errorf("Expected size 4, got %d", heap.Size())
	}

	// The root should be the maximum element
	max, err = heap.Max()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if max.Key != 20 {
		t.Errorf("Expected max key 20, got %d", max.Key)
	}
}

func TestMaxHeap_HeapProperty(t *testing.T) {
	heap := NewMaxHeap[int, string]()

	// Insert elements in random order
	values := []struct {
		key   int
		value string
	}{
		{10, "ten"},
		{30, "thirty"},
		{20, "twenty"},
		{40, "forty"},
		{50, "fifty"},
		{15, "fifteen"},
		{25, "twenty-five"},
	}

	for _, v := range values {
		heap.Insert(v.key, v.value)
	}

	// Verify max heap property: parent >= children
	for i := 0; i < heap.Size()/2; i++ {
		left := Left(i)
		right := Right(i)

		if left < heap.Size() && heap.values[i].Key < heap.values[left].Key {
			t.Errorf("Max heap property violated: parent[%d]=%d < left[%d]=%d",
				i, heap.values[i].Key, left, heap.values[left].Key)
		}

		if right < heap.Size() && heap.values[i].Key < heap.values[right].Key {
			t.Errorf("Max heap property violated: parent[%d]=%d < right[%d]=%d",
				i, heap.values[i].Key, right, heap.values[right].Key)
		}
	}
}

func TestMaxHeap_Pop(t *testing.T) {
	heap := NewMaxHeap[int, string]()

	// Test pop from empty heap
	_, err := heap.Pop()
	if err != ErrorIsEmpty {
		t.Errorf("Expected ErrorIsEmpty, got %v", err)
	}

	// Insert elements and test pop
	heap.Insert(10, "ten")
	heap.Insert(20, "twenty")
	heap.Insert(5, "five")
	heap.Insert(15, "fifteen")

	// Pop should return elements in descending order
	expectedKeys := []int{20, 15, 10, 5}
	for i, expectedKey := range expectedKeys {
		var node *Node[int, string]
		node, err = heap.Pop()
		if err != nil {
			t.Fatalf("Unexpected error at pop %d: %v", i, err)
		}
		if node.Key != expectedKey {
			t.Errorf("Pop %d: expected key %d, got %d", i, expectedKey, node.Key)
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
	heap := NewMaxHeap[int, string]()

	// Test max on empty heap
	_, err := heap.Max()
	if err != ErrorIsEmpty {
		t.Errorf("Expected ErrorIsEmpty, got %v", err)
	}

	// Insert elements and test max
	heap.Insert(10, "ten")
	max, err := heap.Max()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if max.Key != 10 {
		t.Errorf("Expected max key 10, got %d", max.Key)
	}

	heap.Insert(20, "twenty")
	max, err = heap.Max()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if max.Key != 20 {
		t.Errorf("Expected max key 20, got %d", max.Key)
	}

	// Size should remain the same after Max()
	if heap.Size() != 2 {
		t.Errorf("Expected size 2 after Max(), got %d", heap.Size())
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
	heap := NewMaxHeap[int, string]()

	// Manually add elements without maintaining heap property
	heap.values = []*Node[int, string]{
		{Key: 10, Value: "ten"},
		{Key: 30, Value: "thirty"},
		{Key: 20, Value: "twenty"},
		{Key: 40, Value: "forty"},
		{Key: 50, Value: "fifty"},
	}

	// Build max heap
	BuildMaxHeap(heap)

	// Verify max heap property
	for i := 0; i < heap.Size()/2; i++ {
		left := Left(i)
		right := Right(i)

		if left < heap.Size() && heap.values[i].Key < heap.values[left].Key {
			t.Errorf("Max heap property violated after BuildMaxHeap: parent[%d]=%d < left[%d]=%d",
				i, heap.values[i].Key, left, heap.values[left].Key)
		}

		if right < heap.Size() && heap.values[i].Key < heap.values[right].Key {
			t.Errorf("Max heap property violated after BuildMaxHeap: parent[%d]=%d < right[%d]=%d",
				i, heap.values[i].Key, right, heap.values[right].Key)
		}
	}

	// Root should be maximum
	if heap.values[0].Key != 50 {
		t.Errorf("Expected root key 50 after BuildMaxHeap, got %d", heap.values[0].Key)
	}
}

func TestHeapSort(t *testing.T) {
	heap := NewMaxHeap[int, string]()

	// Insert elements in random order
	values := []int{30, 10, 50, 20, 40}
	for _, v := range values {
		heap.Insert(v, "value")
	}

	// Perform heap sort
	HeapSort(heap)

	// After heap sort, the array should be sorted in ascending order
	expected := []int{10, 20, 30, 40, 50}
	for i := 0; i < len(expected); i++ {
		if heap.values[i].Key != expected[i] {
			t.Errorf("After HeapSort, values[%d] = %d; expected %d",
				i, heap.values[i].Key, expected[i])
		}
	}

	// Verify size remains the same
	if heap.Size() != len(values) {
		t.Errorf("Expected size %d after HeapSort, got %d", len(values), heap.Size())
	}
}

func TestHeapSort_LargerDataset(t *testing.T) {
	heap := NewMaxHeap[int, string]()

	// Insert elements in random order (larger dataset)
	values := []int{64, 34, 25, 12, 22, 11, 90, 5, 77, 30, 55, 45, 60, 78, 1}
	for _, v := range values {
		heap.Insert(v, "value")
	}

	// Perform heap sort
	HeapSort(heap)

	// After heap sort, the array should be sorted in ascending order
	expected := []int{1, 5, 11, 12, 22, 25, 30, 34, 45, 55, 60, 64, 77, 78, 90}
	for i := 0; i < len(expected); i++ {
		if heap.values[i].Key != expected[i] {
			t.Errorf("After HeapSort, values[%d] = %d; expected %d",
				i, heap.values[i].Key, expected[i])
		}
	}
}

func TestHeapSort_EmptyHeap(t *testing.T) {
	heap := NewMaxHeap[int, string]()

	// Test heap sort on empty heap
	HeapSort(heap)

	if heap.Size() != 0 {
		t.Errorf("Expected size 0 for empty heap after HeapSort, got %d", heap.Size())
	}
}

func TestHeapSort_SingleElement(t *testing.T) {
	heap := NewMaxHeap[int, string]()
	heap.Insert(42, "forty-two")

	HeapSort(heap)

	if heap.Size() != 1 {
		t.Errorf("Expected size 1 for single element heap after HeapSort, got %d", heap.Size())
	}

	if heap.values[0].Key != 42 {
		t.Errorf("Expected key 42, got %d", heap.values[0].Key)
	}
}

func TestMaxHeap_WithDifferentTypes(t *testing.T) {
	// Test with string keys and int values
	stringHeap := NewMaxHeap[string, int]()
	stringHeap.Insert("zebra", 1)
	stringHeap.Insert("apple", 2)
	stringHeap.Insert("banana", 3)

	// Root should be "zebra" (lexicographically largest)
	max, err := stringHeap.Max()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if max.Key != "zebra" {
		t.Errorf("Expected max key 'zebra', got %s", max.Key)
	}

	// Test with float keys
	floatHeap := NewMaxHeap[float64, string]()
	floatHeap.Insert(3.14, "pi")
	floatHeap.Insert(2.71, "e")
	floatHeap.Insert(1.41, "sqrt2")

	// Root should be 3.14 (largest)
	max2, err := floatHeap.Max()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if max2.Key != 3.14 {
		t.Errorf("Expected max key 3.14, got %f", max2.Key)
	}
}

func TestNode(t *testing.T) {
	node := Node[int, string]{Key: 42, Value: "answer"}

	if node.Key != 42 {
		t.Errorf("Expected key 42, got %d", node.Key)
	}

	if node.Value != "answer" {
		t.Errorf("Expected value 'answer', got %s", node.Value)
	}
}

func TestMaxHeap_IntegrationTest(t *testing.T) {
	heap := NewMaxHeap[int, string]()

	// Insert many elements
	elements := []struct {
		key   int
		value string
	}{
		{50, "fifty"},
		{30, "thirty"},
		{70, "seventy"},
		{20, "twenty"},
		{40, "forty"},
		{60, "sixty"},
		{80, "eighty"},
	}

	for _, elem := range elements {
		heap.Insert(elem.key, elem.value)
	}

	// Verify size
	if heap.Size() != len(elements) {
		t.Errorf("Expected size %d, got %d", len(elements), heap.Size())
	}

	// Pop all elements and verify they come out in descending order
	var poppedKeys []int
	for heap.Size() > 0 {
		node, err := heap.Pop()
		if err != nil {
			t.Fatalf("Unexpected error during pop: %v", err)
		}
		poppedKeys = append(poppedKeys, node.Key)
	}

	// Verify order is descending
	for i := 1; i < len(poppedKeys); i++ {
		if poppedKeys[i-1] < poppedKeys[i] {
			t.Errorf("Elements not popped in descending order: %v", poppedKeys)
			break
		}
	}
}

// Benchmark tests
func BenchmarkMaxHeap_Insert(b *testing.B) {
	heap := NewMaxHeap[int, string]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, "value")
	}
}

func BenchmarkMaxHeap_Pop(b *testing.B) {
	// Pre-populate heap
	heap := NewMaxHeap[int, string]()
	for i := 0; i < b.N; i++ {
		heap.Insert(i, "value")
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
		heap := NewMaxHeap[int, string]()

		// Insert 1000 elements
		for j := 1000; j > 0; j-- {
			heap.Insert(j, "value")
		}

		b.StartTimer()
		HeapSort(heap)
	}
}
