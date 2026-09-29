package main

import "os"

// Задание 6. Двусторонняя очередь.
// Двусвязный список с фиктивными узлами по краям.
type Deque[T any] struct {
	head  *dequeNode[T]
	tail  *dequeNode[T]
	count int
}

type dequeNode[T any] struct {
	prev  *dequeNode[T]
	next  *dequeNode[T]
	value T
}

// 1. Количество элементов в деке.
// Время: O(1)
// Память: O(1)
func (d *Deque[T]) Size() int {
	return d.count
}

// 1. Добавление в голову.
// Время: O(1)
// Память: O(1)
func (d *Deque[T]) AddFront(itm T) {
	d.ensureInit()
	d.insertAfter(d.head, itm)
}

// 1. Добавление в хвост.
// Время: O(1)
// Память: O(1)
func (d *Deque[T]) AddTail(itm T) {
	d.ensureInit()
	d.insertAfter(d.tail.prev, itm)
}

// 1. Снятие элемента с головы.
// Время: O(1)
// Память: O(1)
func (d *Deque[T]) RemoveFront() (T, error) {
	var result T

	if d.count == 0 {
		return result, os.ErrNotExist
	}

	return d.unlink(d.head.next), nil
}

// 1. Снятие элемента с хвоста.
// Время: O(1)
// Память: O(1)
func (d *Deque[T]) RemoveTail() (T, error) {
	var result T

	if d.count == 0 {
		return result, os.ErrNotExist
	}

	return d.unlink(d.tail.prev), nil
}

func (d *Deque[T]) insertAfter(prev *dequeNode[T], itm T) {
	node := &dequeNode[T]{prev: prev, next: prev.next, value: itm}

	prev.next.prev = node
	prev.next = node
	d.count++
}

func (d *Deque[T]) unlink(node *dequeNode[T]) T {
	node.prev.next = node.next
	node.next.prev = node.prev

	node.prev = nil
	node.next = nil
	d.count--

	return node.value
}

func (d *Deque[T]) ensureInit() {
	if d.head != nil {
		return
	}

	d.head = &dequeNode[T]{}
	d.tail = &dequeNode[T]{}
	d.head.next = d.tail
	d.tail.prev = d.head
}

// 2. Как выровнять сложность addHead/removeHead и addTail/removeTail.
//
// В массиве один из концов стоит O(n): при вставке и удалении сдвигаются все
// элементы. Двусвязный список даёт O(1) с обоих концов — так сделано выше.
