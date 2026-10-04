package main

import (
	"constraints"
	"slices"
	"testing"
)

// Задание 7. Упорядоченный список.

func makeList(asc bool, values ...int) *OrderedList[int] {
	l := &OrderedList[int]{}
	l.Clear(asc)

	for _, value := range values {
		l.Add(value)
	}

	return l
}

func assertInts(t *testing.T, actual []int, expected []int) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("index %d: expected %d, got %d", i, expected[i], actual[i])
		}
	}
}

func listValues[T constraints.Ordered](l *OrderedList[T]) []T {
	var values []T

	for node := l.head; node != nil; node = node.next {
		values = append(values, node.value)
	}

	return values
}

// Сверяет значения в обе стороны и Count.
func assertList(t *testing.T, l *OrderedList[int], expected []int) {
	t.Helper()

	var backward []int
	for node := l.tail; node != nil; node = node.prev {
		backward = append([]int{node.value}, backward...)
	}

	assertInts(t, listValues(l), expected)
	assertInts(t, backward, expected)

	if l.Count() != len(expected) {
		t.Errorf("expected count %d, got %d", len(expected), l.Count())
	}
}

// 2, 5. Compare

func TestCompareNumbers(t *testing.T) {
	l := makeList(true)

	cases := []struct {
		v1, v2 int
		want   int
	}{
		{1, 2, -1},
		{2, 1, 1},
		{2, 2, 0},
		{-5, 3, -1},
	}

	for _, c := range cases {
		if got := l.Compare(c.v1, c.v2); got != c.want {
			t.Errorf("Compare(%d, %d): expected %d, got %d", c.v1, c.v2, c.want, got)
		}
	}
}

func TestCompareFloats(t *testing.T) {
	l := &OrderedList[float64]{}

	if l.Compare(1.5, 2.5) != -1 || l.Compare(2.5, 1.5) != 1 || l.Compare(1.5, 1.5) != 0 {
		t.Error("unexpected float comparison")
	}
}

func TestCompareStringsTrimsSpaces(t *testing.T) {
	l := &OrderedList[string]{}

	cases := []struct {
		v1, v2 string
		want   int
	}{
		{"abc", "abd", -1},
		{"b", "a", 1},
		{"  abc  ", "abc", 0},
		{" b", "a ", 1},
		{"\tz\n", "z", 0},
	}

	for _, c := range cases {
		if got := l.Compare(c.v1, c.v2); got != c.want {
			t.Errorf("Compare(%q, %q): expected %d, got %d", c.v1, c.v2, c.want, got)
		}
	}
}

// 1. Clear

func TestClearSetsOrder(t *testing.T) {
	l := makeList(true, 3, 1, 2)

	l.Clear(false)
	assertList(t, l, nil)

	for _, value := range []int{1, 3, 2} {
		l.Add(value)
	}

	assertList(t, l, []int{3, 2, 1})
}

// 3. Add

func TestAddVariants(t *testing.T) {
	cases := []struct {
		name   string
		asc    bool
		values []int
		want   []int
	}{
		{"asc empty", true, nil, nil},
		{"asc single", true, []int{5}, []int{5}},
		{"asc to head", true, []int{5, 3}, []int{3, 5}},
		{"asc to tail", true, []int{3, 5}, []int{3, 5}},
		{"asc to middle", true, []int{1, 9, 5}, []int{1, 5, 9}},
		{"asc duplicates", true, []int{2, 1, 2, 1}, []int{1, 1, 2, 2}},
		{"desc single", false, []int{5}, []int{5}},
		{"desc to head", false, []int{3, 5}, []int{5, 3}},
		{"desc to tail", false, []int{5, 3}, []int{5, 3}},
		{"desc to middle", false, []int{1, 9, 5}, []int{9, 5, 1}},
		{"desc duplicates", false, []int{2, 1, 2, 1}, []int{2, 2, 1, 1}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assertList(t, makeList(c.asc, c.values...), c.want)
		})
	}
}

func TestAddStrings(t *testing.T) {
	l := &OrderedList[string]{}
	l.Clear(true)

	for _, value := range []string{"pear", " apple", "fig "} {
		l.Add(value)
	}

	got := listValues(l)
	want := []string{" apple", "fig ", "pear"}

	if !slices.Equal(got, want) {
		t.Errorf("expected %q, got %q", want, got)
	}
}

// 6. Find

func TestFindVariants(t *testing.T) {
	cases := []struct {
		name  string
		asc   bool
		value int
		found bool
	}{
		{"asc head", true, 1, true},
		{"asc middle", true, 5, true},
		{"asc tail", true, 9, true},
		{"asc before all", true, 0, false},
		{"asc between", true, 4, false},
		{"asc after all", true, 10, false},
		{"desc head", false, 9, true},
		{"desc middle", false, 5, true},
		{"desc tail", false, 1, true},
		{"desc before all", false, 10, false},
		{"desc between", false, 4, false},
		{"desc after all", false, 0, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := makeList(c.asc, 1, 5, 9)

			node, err := l.Find(c.value)

			if !c.found {
				if err == nil {
					t.Errorf("expected error, got node %d", node.value)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if node.value != c.value {
				t.Errorf("expected %d, got %d", c.value, node.value)
			}
		})
	}
}

func TestFindEmpty(t *testing.T) {
	if _, err := makeList(true).Find(1); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestFindReturnsFirstOfDuplicates(t *testing.T) {
	l := makeList(true, 1, 2, 2, 3)

	node, _ := l.Find(2)

	if node.prev == nil || node.prev.value != 1 {
		t.Error("expected first of duplicates")
	}
}

func TestFindStopsEarly(t *testing.T) {
	cases := []struct {
		name  string
		asc   bool
		value int
	}{
		{"asc", true, 4},
		{"desc", false, 6},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := makeList(c.asc, 1, 3, 5, 7, 9)

			// Хвост портим: при честной остановке поиск до него не дойдёт.
			l.tail.value = c.value

			if _, err := l.Find(c.value); err == nil {
				t.Error("expected search to stop before tail")
			}
		})
	}
}

// 4. Delete

func TestDeleteVariants(t *testing.T) {
	cases := []struct {
		name  string
		asc   bool
		value int
		want  []int
	}{
		{"asc head", true, 1, []int{5, 9}},
		{"asc middle", true, 5, []int{1, 9}},
		{"asc tail", true, 9, []int{1, 5}},
		{"asc missing", true, 4, []int{1, 5, 9}},
		{"desc head", false, 9, []int{5, 1}},
		{"desc middle", false, 5, []int{9, 1}},
		{"desc tail", false, 1, []int{9, 5}},
		{"desc missing", false, 4, []int{9, 5, 1}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := makeList(c.asc, 1, 5, 9)

			l.Delete(c.value)

			assertList(t, l, c.want)

			if _, err := l.Find(c.value); err == nil {
				t.Errorf("expected %d to be absent", c.value)
			}
		})
	}
}

func TestDeleteOnlyOneDuplicate(t *testing.T) {
	for _, asc := range []bool{true, false} {
		l := makeList(asc, 2, 2, 2)

		l.Delete(2)

		assertList(t, l, []int{2, 2})
	}
}

func TestDeleteLastItem(t *testing.T) {
	l := makeList(true, 7)

	l.Delete(7)

	assertList(t, l, nil)

	if l.head != nil || l.tail != nil {
		t.Error("expected head and tail to be nil")
	}
}

func TestDeleteFromEmpty(t *testing.T) {
	l := makeList(false)

	l.Delete(1)

	assertList(t, l, nil)
}

// 8*. Удаление дубликатов

func TestRemoveDuplicates(t *testing.T) {
	cases := []struct {
		name   string
		asc    bool
		values []int
		want   []int
	}{
		{"empty", true, nil, nil},
		{"no duplicates", true, []int{1, 2, 3}, []int{1, 2, 3}},
		{"all equal", true, []int{4, 4, 4}, []int{4}},
		{"asc mixed", true, []int{3, 1, 2, 1, 3, 3}, []int{1, 2, 3}},
		{"desc mixed", false, []int{3, 1, 2, 1, 3, 3}, []int{3, 2, 1}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := makeList(c.asc, c.values...)

			l.RemoveDuplicates()

			assertList(t, l, c.want)
		})
	}
}

// 9*. Слияние

func TestMerge(t *testing.T) {
	cases := []struct {
		name      string
		asc1      bool
		asc2      bool
		values1   []int
		values2   []int
		want      []int
		wantOrder bool
	}{
		{"both empty", true, true, nil, nil, nil, true},
		{"first empty", true, true, nil, []int{2, 1}, []int{1, 2}, true},
		{"second empty", false, true, []int{1, 2}, nil, []int{2, 1}, false},
		{"asc asc", true, true, []int{1, 4, 6}, []int{2, 4, 7}, []int{1, 2, 4, 4, 6, 7}, true},
		{"desc desc", false, false, []int{1, 4, 6}, []int{2, 4, 7}, []int{7, 6, 4, 4, 2, 1}, false},
		{"asc desc", true, false, []int{1, 5}, []int{3, 2, 6}, []int{1, 2, 3, 5, 6}, true},
		{"desc asc", false, true, []int{1, 5}, []int{3, 2, 6}, []int{6, 5, 3, 2, 1}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l1 := makeList(c.asc1, c.values1...)
			l2 := makeList(c.asc2, c.values2...)
			before1 := listValues(l1)
			before2 := listValues(l2)

			merged := Merge(l1, l2)

			assertList(t, merged, c.want)

			if merged._ascending != c.wantOrder {
				t.Errorf("expected ascending %v", c.wantOrder)
			}

			assertList(t, l1, before1)
			assertList(t, l2, before2)
		})
	}
}

// 10*. Под-список

func TestHasSublist(t *testing.T) {
	cases := []struct {
		name string
		asc  bool
		sub  []int
		subA bool
		want bool
	}{
		{"empty sub", true, nil, true, true},
		{"whole list", true, []int{1, 2, 2, 3, 5}, true, true},
		{"at head", true, []int{1, 2}, true, true},
		{"at tail", true, []int{3, 5}, true, true},
		{"in middle", true, []int{2, 3}, true, true},
		{"duplicates", true, []int{2, 2, 3}, true, true},
		{"gap", true, []int{1, 3}, true, false},
		{"missing value", true, []int{4}, true, false},
		{"too long", true, []int{1, 2, 2, 3, 5, 6}, true, false},
		{"desc list", false, []int{3, 2}, false, true},
		{"desc gap", false, []int{5, 2}, false, false},
		{"sub in other order", true, []int{2, 3}, false, true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			l := makeList(c.asc, 1, 2, 2, 3, 5)
			sub := makeList(c.subA, c.sub...)

			if got := l.HasSublist(sub); got != c.want {
				t.Errorf("expected %v, got %v", c.want, got)
			}
		})
	}
}

// 11*. Самое частое значение

func TestMostFrequent(t *testing.T) {
	cases := []struct {
		name   string
		asc    bool
		values []int
		want   int
	}{
		{"single", true, []int{7}, 7},
		{"clear winner", true, []int{1, 2, 2, 3, 2}, 2},
		{"winner at tail", true, []int{1, 3, 3, 3, 2}, 3},
		{"asc tie", true, []int{1, 1, 2, 2}, 1},
		{"desc tie", false, []int{1, 1, 2, 2}, 2},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := makeList(c.asc, c.values...).MostFrequent()

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != c.want {
				t.Errorf("expected %d, got %d", c.want, got)
			}
		})
	}
}

func TestMostFrequentEmpty(t *testing.T) {
	if _, err := makeList(true).MostFrequent(); err == nil {
		t.Error("expected error, got nil")
	}
}

// 12*. Индекс за O(log n)

func makeArray(asc bool, values ...int) *OrderedArray[int] {
	a := NewOrderedArray[int](asc)

	for _, value := range values {
		a.Add(value)
	}

	return a
}

func TestOrderedArrayAdd(t *testing.T) {
	assertInts(t, makeArray(true, 5, 1, 3, 1).items, []int{1, 1, 3, 5})
	assertInts(t, makeArray(false, 5, 1, 3, 1).items, []int{5, 3, 1, 1})
}

func TestOrderedArrayIndex(t *testing.T) {
	cases := []struct {
		name  string
		asc   bool
		value int
		want  int
	}{
		{"asc first", true, 1, 0},
		{"asc duplicates", true, 3, 1},
		{"asc last", true, 9, 4},
		{"asc missing", true, 4, -1},
		{"asc after all", true, 10, -1},
		{"desc first", false, 9, 0},
		{"desc duplicates", false, 3, 2},
		{"desc last", false, 1, 4},
		{"desc missing", false, 4, -1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := makeArray(c.asc, 1, 3, 3, 5, 9)

			got, err := a.Index(c.value)

			if c.want == -1 {
				if err == nil {
					t.Errorf("expected error, got index %d", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if got != c.want {
				t.Errorf("expected %d, got %d", c.want, got)
			}
		})
	}
}

func TestOrderedArrayDelete(t *testing.T) {
	a := makeArray(true, 1, 3, 3, 5)

	a.Delete(3)
	assertInts(t, a.items, []int{1, 3, 5})

	a.Delete(4)
	assertInts(t, a.items, []int{1, 3, 5})

	a.Delete(1)
	a.Delete(5)
	a.Delete(3)

	if a.Count() != 0 {
		t.Errorf("expected empty array, got %v", a.items)
	}

	if _, err := a.Index(3); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestOrderedArrayStrings(t *testing.T) {
	a := NewOrderedArray[string](true)

	for _, value := range []string{"pear", "apple", "fig"} {
		a.Add(value)
	}

	if i, err := a.Index("  fig "); err != nil || i != 1 {
		t.Errorf("expected index 1, got %d, %v", i, err)
	}
}
