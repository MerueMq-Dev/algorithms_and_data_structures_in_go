package main

import (
	"math/rand/v2"
	"testing"
)

// Задание 6. Двусторонняя очередь.

func makeDeque(values ...int) *Deque[int] {
	d := &Deque[int]{}

	for _, value := range values {
		d.AddTail(value)
	}

	return d
}

func items(d *Deque[int]) []int {
	var values []int

	if d.head == nil {
		return values
	}

	for node := d.head.next; node != d.tail; node = node.next {
		values = append(values, node.value)
	}

	return values
}

func contains(values []int, itm int) bool {
	for _, value := range values {
		if value == itm {
			return true
		}
	}

	return false
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

func assertSize(t *testing.T, d interface{ Size() int }, expected int) {
	t.Helper()

	if d.Size() != expected {
		t.Fatalf("expected size %d, got %d", expected, d.Size())
	}
}

// 1. Size

func TestSizeEmptyDeque(t *testing.T) {
	assertSize(t, makeDeque(), 0)
}

func TestSizeAfterAdds(t *testing.T) {
	d := makeDeque(10, 20)
	d.AddFront(5)

	assertSize(t, d, 3)
}

func TestSizeAfterRemoves(t *testing.T) {
	d := makeDeque(10, 20, 30)

	d.RemoveFront()
	d.RemoveTail()

	assertSize(t, d, 1)
}

// 3. AddFront

func TestAddFrontToEmpty(t *testing.T) {
	d := makeDeque()

	d.AddFront(10)

	assertSize(t, d, 1)
	assertInts(t, items(d), []int{10})
}

func TestAddFrontToNonEmpty(t *testing.T) {
	d := makeDeque(10, 20)

	d.AddFront(5)

	assertSize(t, d, 3)

	if !contains(items(d), 5) {
		t.Fatal("expected 5 to be in deque")
	}

	assertInts(t, items(d), []int{5, 10, 20})
}

func TestAddFrontReversesOrder(t *testing.T) {
	d := makeDeque()

	for _, value := range []int{1, 2, 3} {
		d.AddFront(value)
	}

	assertInts(t, items(d), []int{3, 2, 1})
}

// 3. AddTail

func TestAddTailToEmpty(t *testing.T) {
	d := makeDeque()

	d.AddTail(10)

	assertSize(t, d, 1)
	assertInts(t, items(d), []int{10})
}

func TestAddTailToNonEmpty(t *testing.T) {
	d := makeDeque(10, 20)

	d.AddTail(30)

	assertSize(t, d, 3)

	if !contains(items(d), 30) {
		t.Fatal("expected 30 to be in deque")
	}

	assertInts(t, items(d), []int{10, 20, 30})
}

func TestAddTailAfterEmptying(t *testing.T) {
	d := makeDeque(10)

	d.RemoveFront()
	d.AddTail(20)

	assertSize(t, d, 1)
	assertInts(t, items(d), []int{20})
}

// 3. RemoveFront

func TestRemoveFrontReturnsHead(t *testing.T) {
	d := makeDeque(10, 20, 30)

	got, err := d.RemoveFront()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 10 {
		t.Errorf("expected 10, got %d", got)
	}

	assertSize(t, d, 2)

	if contains(items(d), 10) {
		t.Error("expected 10 to be removed")
	}

	assertInts(t, items(d), []int{20, 30})
}

func TestRemoveFrontEmpty(t *testing.T) {
	d := makeDeque()

	if _, err := d.RemoveFront(); err == nil {
		t.Error("expected error, got nil")
	}

	assertSize(t, d, 0)
}

func TestRemoveFrontLastItem(t *testing.T) {
	d := makeDeque(10)

	d.RemoveFront()

	assertSize(t, d, 0)

	if d.head.next != d.tail || d.tail.prev != d.head {
		t.Error("expected dummy nodes to point at each other")
	}

	if _, err := d.RemoveFront(); err == nil {
		t.Error("expected error, got nil")
	}
}

// 3. RemoveTail

func TestRemoveTailReturnsTail(t *testing.T) {
	d := makeDeque(10, 20, 30)

	got, err := d.RemoveTail()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 30 {
		t.Errorf("expected 30, got %d", got)
	}

	assertSize(t, d, 2)

	if contains(items(d), 30) {
		t.Error("expected 30 to be removed")
	}

	assertInts(t, items(d), []int{10, 20})
}

func TestRemoveTailEmpty(t *testing.T) {
	d := makeDeque()

	if _, err := d.RemoveTail(); err == nil {
		t.Error("expected error, got nil")
	}

	assertSize(t, d, 0)
}

func TestRemoveTailLastItem(t *testing.T) {
	d := makeDeque(10)

	d.RemoveTail()

	assertSize(t, d, 0)

	if d.head.next != d.tail || d.tail.prev != d.head {
		t.Error("expected dummy nodes to point at each other")
	}
}

func TestMixedOperations(t *testing.T) {
	d := makeDeque()

	d.AddTail(2)
	d.AddFront(1)
	d.AddTail(3)
	d.RemoveFront()
	d.AddFront(0)
	d.RemoveTail()

	assertSize(t, d, 2)
	assertInts(t, items(d), []int{0, 2})
}

func TestDequeOfStrings(t *testing.T) {
	d := &Deque[string]{}

	d.AddFront("b")
	d.AddFront("a")
	d.AddTail("c")

	first, _ := d.RemoveFront()
	last, _ := d.RemoveTail()

	if first != "a" || last != "c" {
		t.Errorf("expected a and c, got %q and %q", first, last)
	}
}

// 4*. Палиндром

func TestIsPalindromeVariants(t *testing.T) {
	cases := []struct {
		input string
		want  bool
	}{
		{"", true},
		{"a", true},
		{"aa", true},
		{"ab", false},
		{"aba", true},
		{"abba", true},
		{"abca", false},
		{"шалаш", true},
		{"Шалаш", false},
		{"топот", true},
		{"a b a", true},
	}

	for _, c := range cases {
		t.Run(c.input, func(t *testing.T) {
			if got := IsPalindrome(c.input); got != c.want {
				t.Errorf("IsPalindrome(%q): expected %v, got %v", c.input, c.want, got)
			}
		})
	}
}

// 5*. Минимум деки

func assertMin(t *testing.T, d *MinDeque[int], expected int) {
	t.Helper()

	got, err := d.Min()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != expected {
		t.Errorf("expected min %d, got %d", expected, got)
	}
}

func TestMinDequeEmpty(t *testing.T) {
	d := &MinDeque[int]{}

	if _, err := d.Min(); err == nil {
		t.Error("expected error, got nil")
	}

	if _, err := d.RemoveFront(); err == nil {
		t.Error("expected error, got nil")
	}

	if _, err := d.RemoveTail(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestMinDequeAddBothSides(t *testing.T) {
	d := &MinDeque[int]{}

	d.AddTail(5)
	assertMin(t, d, 5)

	d.AddFront(3)
	assertMin(t, d, 3)

	d.AddTail(7)
	assertMin(t, d, 3)

	d.AddTail(1)
	assertMin(t, d, 1)
}

func TestMinDequeRemoveRestoresMin(t *testing.T) {
	d := &MinDeque[int]{}

	for _, value := range []int{1, 4, 2, 6, 3} {
		d.AddTail(value)
	}

	cases := []struct {
		remove func() (int, error)
		item   int
		min    int
	}{
		{d.RemoveFront, 1, 2},
		{d.RemoveTail, 3, 2},
		{d.RemoveFront, 4, 2},
		{d.RemoveFront, 2, 6},
	}

	for _, c := range cases {
		got, err := c.remove()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != c.item {
			t.Errorf("expected %d, got %d", c.item, got)
		}

		assertMin(t, d, c.min)
	}

	assertSize(t, d, 1)
}

func TestMinDequeDuplicates(t *testing.T) {
	d := &MinDeque[int]{}

	d.AddTail(2)
	d.AddTail(1)
	d.AddFront(1)

	d.RemoveFront()
	assertMin(t, d, 1)

	d.RemoveTail()
	assertMin(t, d, 2)
}

func TestMinDequeRandomAgainstModel(t *testing.T) {
	r := rand.New(rand.NewPCG(1, 2))
	d := &MinDeque[int]{}
	var model []int

	for step := 0; step < 5000; step++ {
		value := r.IntN(100)

		switch r.IntN(4) {
		case 0:
			d.AddFront(value)
			model = append([]int{value}, model...)
		case 1:
			d.AddTail(value)
			model = append(model, value)
		case 2:
			got, err := d.RemoveFront()

			if len(model) == 0 {
				if err == nil {
					t.Fatalf("step %d: expected error, got nil", step)
				}
				continue
			}

			if got != model[0] {
				t.Fatalf("step %d: expected %d, got %d", step, model[0], got)
			}

			model = model[1:]
		case 3:
			got, err := d.RemoveTail()

			if len(model) == 0 {
				if err == nil {
					t.Fatalf("step %d: expected error, got nil", step)
				}
				continue
			}

			last := len(model) - 1

			if got != model[last] {
				t.Fatalf("step %d: expected %d, got %d", step, model[last], got)
			}

			model = model[:last]
		}

		assertSize(t, d, len(model))

		if len(model) > 0 {
			assertMin(t, d, sliceMin(model))
		}
	}
}

func sliceMin(values []int) int {
	result := values[0]

	for _, value := range values[1:] {
		result = min(result, value)
	}

	return result
}

// 6*. Дека на динамическом массиве

func TestArrayDequeEmpty(t *testing.T) {
	d := &ArrayDeque[int]{}

	assertSize(t, d, 0)

	if _, err := d.RemoveFront(); err == nil {
		t.Error("expected error, got nil")
	}

	if _, err := d.RemoveTail(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestArrayDequeBothEnds(t *testing.T) {
	d := &ArrayDeque[int]{}

	d.AddFront(2)
	d.AddFront(1)
	d.AddTail(3)

	assertSize(t, d, 3)

	if d.head != len(d.buffer)-2 {
		t.Errorf("expected head to wrap to %d, got %d", len(d.buffer)-2, d.head)
	}

	first, _ := d.RemoveFront()
	last, _ := d.RemoveTail()

	if first != 1 || last != 3 {
		t.Errorf("expected 1 and 3, got %d and %d", first, last)
	}

	assertSize(t, d, 1)
}

func TestArrayDequeGrowKeepsOrder(t *testing.T) {
	d := &ArrayDeque[int]{}

	for i := 1; i <= minCapacity; i++ {
		d.AddFront(-i)
		d.AddTail(i)
	}

	if len(d.buffer) != minCapacity*growFactor {
		t.Errorf("expected capacity %d, got %d", minCapacity*growFactor, len(d.buffer))
	}

	for i := -minCapacity; i <= minCapacity; i++ {
		if i == 0 {
			continue
		}

		got, err := d.RemoveFront()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != i {
			t.Fatalf("expected %d, got %d", i, got)
		}
	}

	assertSize(t, d, 0)
}

func TestArrayDequeShrinks(t *testing.T) {
	d := &ArrayDeque[int]{}

	for i := 0; i < 100; i++ {
		d.AddTail(i)
	}

	grown := len(d.buffer)

	for i := 0; i < 95; i++ {
		d.RemoveFront()
	}

	if len(d.buffer) >= grown {
		t.Errorf("expected buffer to shrink from %d, got %d", grown, len(d.buffer))
	}

	if len(d.buffer) < minCapacity {
		t.Errorf("expected capacity at least %d, got %d", minCapacity, len(d.buffer))
	}

	for want := 95; want < 100; want++ {
		got, _ := d.RemoveFront()

		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}
}

func TestArrayDequeRandomAgainstModel(t *testing.T) {
	r := rand.New(rand.NewPCG(3, 4))
	d := &ArrayDeque[int]{}
	var model []int

	maxCapacity := 0

	for step := 0; step < 5000; step++ {
		value := r.IntN(1000)

		// Сначала буфер растёт, потом сжимается.
		add := r.IntN(4) < 3
		if step >= 2500 {
			add = !add
		}

		switch {
		case add && r.IntN(2) == 0:
			d.AddFront(value)
			model = append([]int{value}, model...)
		case add:
			d.AddTail(value)
			model = append(model, value)
		case len(model) == 0:
			continue
		case r.IntN(2) == 0:
			got, _ := d.RemoveFront()

			if got != model[0] {
				t.Fatalf("step %d: expected %d, got %d", step, model[0], got)
			}

			model = model[1:]
		default:
			got, _ := d.RemoveTail()
			last := len(model) - 1

			if got != model[last] {
				t.Fatalf("step %d: expected %d, got %d", step, model[last], got)
			}

			model = model[:last]
		}

		assertSize(t, d, len(model))
		maxCapacity = max(maxCapacity, len(d.buffer))
	}

	if maxCapacity <= minCapacity {
		t.Errorf("expected buffer to grow beyond %d", minCapacity)
	}

	if len(d.buffer) >= maxCapacity {
		t.Errorf("expected buffer to shrink from %d, got %d", maxCapacity, len(d.buffer))
	}
}
