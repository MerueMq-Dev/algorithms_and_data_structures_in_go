package main

import "testing"

// Задание 5. Очередь.

func makeQueue(values ...int) *Queue[int] {
	q := &Queue[int]{}

	for _, value := range values {
		q.Enqueue(value)
	}

	return q
}

// drain достаёт все элементы по порядку и опустошает очередь.
func drain(t *testing.T, q *Queue[int]) []int {
	t.Helper()

	var values []int

	for q.Size() > 0 {
		item, err := q.Dequeue()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		values = append(values, item)
	}

	return values
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

// 1. Size

func TestSizeEmptyQueue(t *testing.T) {
	q := makeQueue()

	if q.Size() != 0 {
		t.Errorf("expected 0, got %d", q.Size())
	}
}

func TestSizeAfterEnqueue(t *testing.T) {
	q := makeQueue(10, 20, 30)

	if q.Size() != 3 {
		t.Errorf("expected 3, got %d", q.Size())
	}
}

func TestSizeAfterDequeue(t *testing.T) {
	q := makeQueue(10, 20, 30)

	q.Dequeue()

	if q.Size() != 2 {
		t.Errorf("expected 2, got %d", q.Size())
	}
}

// 1. Enqueue

func TestEnqueueOneItem(t *testing.T) {
	q := makeQueue(10)

	assertInts(t, drain(t, q), []int{10})
}

func TestEnqueueKeepsOrder(t *testing.T) {
	q := makeQueue(10, 20, 30)

	assertInts(t, drain(t, q), []int{10, 20, 30})
}

func TestEnqueueAfterEmptying(t *testing.T) {
	q := makeQueue(10)

	q.Dequeue()
	q.Enqueue(20)
	q.Enqueue(30)

	assertInts(t, drain(t, q), []int{20, 30})
}

// 1. Dequeue

func TestDequeueReturnsHead(t *testing.T) {
	q := makeQueue(10, 20)

	got, err := q.Dequeue()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 10 {
		t.Errorf("expected 10, got %d", got)
	}
}

func TestDequeueEmptyQueue(t *testing.T) {
	q := makeQueue()

	if _, err := q.Dequeue(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestDequeueAfterEmptying(t *testing.T) {
	q := makeQueue(10)

	q.Dequeue()

	if _, err := q.Dequeue(); err == nil {
		t.Error("expected error, got nil")
	}

	if q.head != nil || q.tail != nil {
		t.Error("expected head and tail to be nil in empty queue")
	}
}

func TestDequeueMixedWithEnqueue(t *testing.T) {
	q := makeQueue(10, 20)

	q.Dequeue()
	q.Enqueue(30)

	assertInts(t, drain(t, q), []int{20, 30})
}

func TestLoopFromTask(t *testing.T) {
	q := &Queue[int]{}
	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)

	var order []int

	for q.Size() > 0 {
		item, _ := q.Dequeue()
		order = append(order, item)
	}

	assertInts(t, order, []int{1, 2, 3})
}

func TestQueueOfStrings(t *testing.T) {
	q := &Queue[string]{}

	q.Enqueue("a")
	q.Enqueue("b")

	got, err := q.Dequeue()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "a" {
		t.Errorf("expected %q, got %q", "a", got)
	}
}

// 3*. Вращение очереди

func TestRotateEmptyQueue(t *testing.T) {
	q := makeQueue()

	RotateQueue(q, 3)

	if q.Size() != 0 {
		t.Errorf("expected empty queue, got size %d", q.Size())
	}
}

func TestRotateVariants(t *testing.T) {
	cases := []struct {
		name string
		n    int
		want []int
	}{
		{"zero", 0, []int{1, 2, 3, 4}},
		{"by one", 1, []int{2, 3, 4, 1}},
		{"by two", 2, []int{3, 4, 1, 2}},
		{"full circle", 4, []int{1, 2, 3, 4}},
		{"more than size", 5, []int{2, 3, 4, 1}},
		{"negative", -1, []int{4, 1, 2, 3}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			q := makeQueue(1, 2, 3, 4)

			RotateQueue(q, c.n)

			assertInts(t, drain(t, q), c.want)
		})
	}
}

// 4*. Очередь на двух стеках

func TestStackQueueEmpty(t *testing.T) {
	q := &StackQueue[int]{}

	if q.Size() != 0 {
		t.Errorf("expected 0, got %d", q.Size())
	}

	if _, err := q.Dequeue(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestStackQueueKeepsOrder(t *testing.T) {
	q := &StackQueue[int]{}

	for _, value := range []int{10, 20, 30} {
		q.Enqueue(value)
	}

	for _, want := range []int{10, 20, 30} {
		got, err := q.Dequeue()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}
}

func TestStackQueueMixedOperations(t *testing.T) {
	q := &StackQueue[int]{}

	q.Enqueue(1)
	q.Enqueue(2)
	q.Dequeue()

	// Выходной стек ещё не пуст, новые элементы идут во входной.
	q.Enqueue(3)
	q.Enqueue(4)

	if q.Size() != 3 {
		t.Fatalf("expected size 3, got %d", q.Size())
	}

	for _, want := range []int{2, 3, 4} {
		got, _ := q.Dequeue()

		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}
}

// 5*. Обращение очереди

func TestReverseQueueVariants(t *testing.T) {
	cases := []struct {
		name   string
		values []int
		want   []int
	}{
		{"empty", nil, nil},
		{"single", []int{10}, []int{10}},
		{"many", []int{1, 2, 3, 4}, []int{4, 3, 2, 1}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			q := makeQueue(c.values...)

			ReverseQueue(q)

			assertInts(t, drain(t, q), c.want)
		})
	}
}

func TestReverseQueueTwiceRestoresOrder(t *testing.T) {
	q := makeQueue(1, 2, 3)

	ReverseQueue(q)
	ReverseQueue(q)

	assertInts(t, drain(t, q), []int{1, 2, 3})
}

// 6*. Круговая очередь

func TestCircularQueueEmpty(t *testing.T) {
	q := NewCircularQueue[int](3)

	if q.Size() != 0 || q.IsFull() {
		t.Error("new queue must be empty and not full")
	}

	if _, err := q.Dequeue(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestCircularQueueFull(t *testing.T) {
	q := NewCircularQueue[int](3)

	for _, value := range []int{1, 2, 3} {
		if err := q.Enqueue(value); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	if !q.IsFull() {
		t.Error("expected queue to be full")
	}

	if err := q.Enqueue(4); err == nil {
		t.Error("expected error on full queue, got nil")
	}

	if q.Size() != 3 {
		t.Errorf("expected size 3, got %d", q.Size())
	}
}

func TestCircularQueueWrapsAround(t *testing.T) {
	q := NewCircularQueue[int](3)

	q.Enqueue(1)
	q.Enqueue(2)
	q.Enqueue(3)
	q.Dequeue()
	q.Dequeue()

	// Хвост доходит до конца массива и начинается сначала.
	q.Enqueue(4)
	q.Enqueue(5)

	if !q.IsFull() {
		t.Error("expected queue to be full after wrap-around")
	}

	for _, want := range []int{3, 4, 5} {
		got, err := q.Dequeue()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}
}

func TestCircularQueueZeroCapacity(t *testing.T) {
	q := NewCircularQueue[int](0)

	if !q.IsFull() {
		t.Error("zero-capacity queue must be full")
	}

	if err := q.Enqueue(1); err == nil {
		t.Error("expected error, got nil")
	}

	if _, err := q.Dequeue(); err == nil {
		t.Error("expected error, got nil")
	}
}
