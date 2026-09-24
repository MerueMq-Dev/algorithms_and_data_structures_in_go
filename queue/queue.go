package main

import "os"

// Задание 5. Очередь.
type Queue[T any] struct {
	head  *queueNode[T]
	tail  *queueNode[T]
	count int
}

type queueNode[T any] struct {
	next  *queueNode[T]
	value T
}

// 1. Количество элементов в очереди.
// Время: O(1)
// Память: O(1)
func (q *Queue[T]) Size() int {
	return q.count
}

// 1. Снятие элемента с головы очереди.
// 2. Время: O(1) — голова под рукой, ничего сдвигать не нужно.
// Память: O(1)
func (q *Queue[T]) Dequeue() (T, error) {
	var result T

	if q.head == nil {
		return result, os.ErrNotExist
	}

	node := q.head
	result = node.value

	q.head = node.next
	node.next = nil
	q.count--

	if q.head == nil {
		q.tail = nil
	}

	return result, nil
}

// 1. Добавление в хвост очереди.
// 2. Время: O(1) — хвост хранится в отдельном поле.
// Память: O(1)
func (q *Queue[T]) Enqueue(itm T) {
	node := &queueNode[T]{value: itm}

	if q.tail == nil {
		q.head = node
	} else {
		q.tail.next = node
	}

	q.tail = node
	q.count++
}
