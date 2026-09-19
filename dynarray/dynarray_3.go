package main

import "testing"

// Задание 3. Динамический массив.

func makeArray(values ...int) *DynArray[int] {
	array := &DynArray[int]{}
	array.Init()

	for _, value := range values {
		array.Append(value)
	}

	return array
}

// fill добавляет n элементов: 0, 1, ... n-1.
func fill(n int) *DynArray[int] {
	array := &DynArray[int]{}
	array.Init()

	for i := 0; i < n; i++ {
		array.Append(i)
	}

	return array
}

func assertArray(t *testing.T, array *DynArray[int], expected []int) {
	t.Helper()

	if array.count != len(expected) {
		t.Fatalf("expected count %d, got %d", len(expected), array.count)
	}

	for i, want := range expected {
		got, err := array.GetItem(i)

		if err != nil {
			t.Fatalf("unexpected error at %d: %v", i, err)
		}

		if got != want {
			t.Errorf("index %d: expected %d, got %d", i, want, got)
		}
	}

	if array.capacity < array.count {
		t.Errorf("capacity %d is less than count %d", array.capacity, array.count)
	}

	if len(array.array) != array.capacity {
		t.Errorf("buffer length %d does not match capacity %d", len(array.array), array.capacity)
	}
}

func assertCapacity(t *testing.T, array *DynArray[int], expected int) {
	t.Helper()

	if array.capacity != expected {
		t.Errorf("expected capacity %d, got %d", expected, array.capacity)
	}
}

// Init и MakeArray

func TestInitCreatesMinimalBuffer(t *testing.T) {
	array := &DynArray[int]{}
	array.Init()

	assertArray(t, array, []int{})
	assertCapacity(t, array, 16)
}

func TestMakeArrayKeepsItems(t *testing.T) {
	array := makeArray(1, 2, 3)

	array.MakeArray(64)

	assertArray(t, array, []int{1, 2, 3})
	assertCapacity(t, array, 64)
}

func TestMakeArrayKeepsMinimalCapacity(t *testing.T) {
	array := makeArray(1, 2)

	array.MakeArray(4)

	assertArray(t, array, []int{1, 2})
	assertCapacity(t, array, 16)
}

// Append и GetItem

func TestAppendGrowsBufferTwice(t *testing.T) {
	array := fill(16)

	assertCapacity(t, array, 16)

	array.Append(100)

	if array.count != 17 {
		t.Errorf("expected count 17, got %d", array.count)
	}

	assertCapacity(t, array, 32)
}

func TestGetItemReturnsValues(t *testing.T) {
	array := makeArray(10, 20, 30)

	for i, want := range []int{10, 20, 30} {
		got, err := array.GetItem(i)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != want {
			t.Errorf("index %d: expected %d, got %d", i, want, got)
		}
	}
}

func TestGetItemBadIndex(t *testing.T) {
	array := makeArray(10, 20)

	for _, index := range []int{-1, 2, 100} {
		if _, err := array.GetItem(index); err == nil {
			t.Errorf("expected error for index %d, got nil", index)
		}
	}
}

func TestGetItemOnEmptyArray(t *testing.T) {
	array := makeArray()

	if _, err := array.GetItem(0); err == nil {
		t.Error("expected error for empty array, got nil")
	}
}

// Задача 5. Insert: буфер не превышен

func TestInsertWithoutGrowth(t *testing.T) {
	array := makeArray(10, 20, 30)

	if err := array.Insert(15, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertArray(t, array, []int{10, 15, 20, 30})
	assertCapacity(t, array, 16)
}

func TestInsertAtBeginning(t *testing.T) {
	array := makeArray(10, 20)

	if err := array.Insert(5, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertArray(t, array, []int{5, 10, 20})
	assertCapacity(t, array, 16)
}

func TestInsertAtCountAppendsToTail(t *testing.T) {
	array := makeArray(10, 20)

	if err := array.Insert(30, array.count); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertArray(t, array, []int{10, 20, 30})
}

func TestInsertIntoEmptyArray(t *testing.T) {
	array := makeArray()

	if err := array.Insert(10, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertArray(t, array, []int{10})
	assertCapacity(t, array, 16)
}

// Задача 5. Insert: буфер превышен

func TestInsertGrowsBuffer(t *testing.T) {
	array := fill(16)

	assertCapacity(t, array, 16)

	if err := array.Insert(100, 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertCapacity(t, array, 32)

	if array.count != 17 {
		t.Errorf("expected count 17, got %d", array.count)
	}

	first, _ := array.GetItem(0)
	last, _ := array.GetItem(16)

	if first != 100 {
		t.Errorf("expected first item 100, got %d", first)
	}

	if last != 15 {
		t.Errorf("expected last item 15, got %d", last)
	}
}

func TestInsertGrowsBufferTwice(t *testing.T) {
	array := fill(32)

	assertCapacity(t, array, 32)

	if err := array.Insert(100, 32); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertCapacity(t, array, 64)
}

// Задача 5. Insert: недопустимая позиция

func TestInsertBadIndex(t *testing.T) {
	array := makeArray(10, 20, 30)

	for _, index := range []int{-1, 4, 100} {
		if err := array.Insert(99, index); err == nil {
			t.Errorf("expected error for index %d, got nil", index)
		}
	}

	assertArray(t, array, []int{10, 20, 30})
}

// Задача 5. Remove: буфер не меняется

func TestRemoveKeepsBuffer(t *testing.T) {
	array := fill(17)

	assertCapacity(t, array, 32)

	if err := array.Remove(0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if array.count != 16 {
		t.Errorf("expected count 16, got %d", array.count)
	}

	assertCapacity(t, array, 32)
}

func TestRemoveKeepsMinimalBuffer(t *testing.T) {
	array := makeArray(10, 20, 30)

	if err := array.Remove(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertArray(t, array, []int{10, 30})
	assertCapacity(t, array, 16)
}

func TestRemoveFromBeginningAndEnd(t *testing.T) {
	array := makeArray(10, 20, 30)

	if err := array.Remove(0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertArray(t, array, []int{20, 30})

	if err := array.Remove(1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertArray(t, array, []int{20})
}

func TestRemoveOnlyItem(t *testing.T) {
	array := makeArray(10)

	if err := array.Remove(0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertArray(t, array, []int{})
	assertCapacity(t, array, 16)
}

// Задача 5. Remove: буфер уменьшается

func TestRemoveShrinksBuffer(t *testing.T) {
	array := fill(17)

	assertCapacity(t, array, 32)

	// 17 -> 16: заполнено ровно наполовину, буфер прежний.
	if err := array.Remove(0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertCapacity(t, array, 32)

	// 16 -> 15: заполнено меньше половины, буфер делится на 1.5.
	if err := array.Remove(0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if array.count != 15 {
		t.Errorf("expected count 15, got %d", array.count)
	}

	assertCapacity(t, array, 21)
}

func TestRemoveShrinksBufferRepeatedly(t *testing.T) {
	array := fill(33)

	assertCapacity(t, array, 64)

	for array.count > 31 {
		if err := array.Remove(0); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// 31/64 < 0.5 -> 64 / 1.5 = 42
	assertCapacity(t, array, 42)

	for array.count > 20 {
		if err := array.Remove(0); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	// 20/42 < 0.5 -> 42 / 1.5 = 28
	assertCapacity(t, array, 28)
}

func TestRemoveNeverShrinksBelowMinimum(t *testing.T) {
	array := fill(20)

	for array.count > 0 {
		if err := array.Remove(0); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	assertArray(t, array, []int{})
	assertCapacity(t, array, 16)
}

func TestRemoveKeepsRemainingItems(t *testing.T) {
	array := fill(17)

	if err := array.Remove(8); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := make([]int, 0, 16)
	for i := 0; i < 17; i++ {
		if i != 8 {
			expected = append(expected, i)
		}
	}

	assertArray(t, array, expected)
}

// Задача 5. Remove: недопустимая позиция

func TestRemoveBadIndex(t *testing.T) {
	array := makeArray(10, 20, 30)

	for _, index := range []int{-1, 3, 100} {
		if err := array.Remove(index); err == nil {
			t.Errorf("expected error for index %d, got nil", index)
		}
	}

	assertArray(t, array, []int{10, 20, 30})
}

func TestRemoveFromEmptyArray(t *testing.T) {
	array := makeArray()

	if err := array.Remove(0); err == nil {
		t.Error("expected error for empty array, got nil")
	}
}

// Совместная работа Insert и Remove

func TestInsertAndRemoveMixed(t *testing.T) {
	array := makeArray(10, 30)

	if err := array.Insert(20, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := array.Insert(40, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertArray(t, array, []int{10, 20, 30, 40})

	if err := array.Remove(0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertArray(t, array, []int{20, 30, 40})
}

func TestArrayOfStrings(t *testing.T) {
	array := &DynArray[string]{}
	array.Init()

	array.Append("a")

	if err := array.Insert("b", 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := array.Insert("c", 0); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for i, want := range []string{"c", "a", "b"} {
		got, err := array.GetItem(i)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != want {
			t.Errorf("index %d: expected %q, got %q", i, want, got)
		}
	}
}

// 6*. Банковский метод

func TestBankArrayKeepsItems(t *testing.T) {
	array := NewBankDynArray[int]()

	for i := 0; i < 20; i++ {
		array.Append(i)
	}

	if array.count != 20 {
		t.Fatalf("expected count 20, got %d", array.count)
	}

	for i := 0; i < 20; i++ {
		got, err := array.GetItem(i)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != i {
			t.Errorf("index %d: expected %d, got %d", i, i, got)
		}
	}
}

func TestBankBalanceNeverNegative(t *testing.T) {
	array := NewBankDynArray[int]()

	for i := 0; i < 1000; i++ {
		array.Append(i)

		if array.Bank() < 0 {
			t.Fatalf("bank went negative at append %d: %d", i, array.Bank())
		}
	}
}

func TestBankPaysForAllCopies(t *testing.T) {
	array := NewBankDynArray[int]()

	for i := 0; i < 500; i++ {
		array.Append(i)
	}

	if array.TotalPaid() < array.TotalCopied()+array.count {
		t.Errorf("paid %d is not enough for %d copies and %d writes",
			array.TotalPaid(), array.TotalCopied(), array.count)
	}
}

func TestBankAmortizedCostIsConstant(t *testing.T) {
	array := NewBankDynArray[int]()

	const n = 1000
	for i := 0; i < n; i++ {
		array.Append(i)
	}

	work := n + array.TotalCopied()

	if work > appendCost*n {
		t.Errorf("total work %d exceeds amortized bound %d", work, appendCost*n)
	}
}

// 7*. Многомерный массив

func TestMultiArrayRequiresDimensions(t *testing.T) {
	if _, err := NewMultiArray[int](); err == nil {
		t.Error("expected error for zero dimensions, got nil")
	}

	if _, err := NewMultiArray[int](2, -1); err == nil {
		t.Error("expected error for negative dimension, got nil")
	}
}

func TestMultiArraySetAndGet(t *testing.T) {
	array, err := NewMultiArray[int](2, 3, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := array.Set(42, 1, 2, 3); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := array.Get(1, 2, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 42 {
		t.Errorf("expected 42, got %d", got)
	}
}

func TestMultiArrayEmptyCellIsZero(t *testing.T) {
	array, _ := NewMultiArray[int](2, 2)

	got, err := array.Get(1, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 0 {
		t.Errorf("expected zero value, got %d", got)
	}
}

func TestMultiArrayKeepsCellsIndependent(t *testing.T) {
	array, _ := NewMultiArray[int](2, 2)

	cells := [][2]int{{0, 0}, {0, 1}, {1, 0}, {1, 1}}

	for i, cell := range cells {
		if err := array.Set(i+1, cell[0], cell[1]); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	for i, cell := range cells {
		got, err := array.Get(cell[0], cell[1])

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != i+1 {
			t.Errorf("cell %v: expected %d, got %d", cell, i+1, got)
		}
	}
}

func TestMultiArrayGrowsDimension(t *testing.T) {
	array, _ := NewMultiArray[int](2, 2)

	if err := array.Set(7, 5, 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dims := array.Dims()

	if dims[0] != 6 {
		t.Errorf("expected first dimension 6, got %d", dims[0])
	}

	got, err := array.Get(5, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 7 {
		t.Errorf("expected 7, got %d", got)
	}
}

func TestMultiArrayBadIndices(t *testing.T) {
	array, _ := NewMultiArray[int](2, 2)

	if _, err := array.Get(0); err == nil {
		t.Error("expected error for wrong indices count, got nil")
	}

	if _, err := array.Get(2, 0); err == nil {
		t.Error("expected error for out of range index, got nil")
	}

	if _, err := array.Get(-1, 0); err == nil {
		t.Error("expected error for negative index, got nil")
	}

	if err := array.Set(1, 0); err == nil {
		t.Error("expected error for wrong indices count, got nil")
	}

	if err := array.Set(1, -1, 0); err == nil {
		t.Error("expected error for negative index, got nil")
	}
}

func TestMultiArraySingleDimension(t *testing.T) {
	array, err := NewMultiArray[string](3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := array.Set("hello", 2); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, err := array.Get(2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "hello" {
		t.Errorf("expected %q, got %q", "hello", got)
	}
}
