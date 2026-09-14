package main

import (
	"os"
	"reflect"
)

// Задание 2. Двусвязный (двунаправленный) список.
func equalValues(a, b int) bool {
	return reflect.DeepEqual(a, b)
}

type Node struct {
	prev  *Node
	next  *Node
	value int
}

type LinkedList2 struct {
	head *Node
	tail *Node
}

func (l *LinkedList2) AddInTail(item Node) {
	if l.head == nil {
		l.head = &item
		l.head.next = nil
		l.head.prev = nil
	} else {
		l.tail.next = &item
		item.prev = l.tail
	}

	l.tail = &item
	l.tail.next = nil
}

// 1. Поиск первого узла по его значению.
// Время: O(n)
// Память: O(1)
func (l *LinkedList2) Find(n int) (Node, error) {
	for node := l.head; node != nil; node = node.next {
		if equalValues(node.value, n) {
			return *node, nil
		}
	}

	return Node{value: -1, next: nil}, os.ErrNotExist
}

// 2. Поиск всех узлов по конкретному значению.
// Время: O(n)
// Память: O(k), где k — количество найденных узлов.
func (l *LinkedList2) FindAll(n int) []Node {
	var nodes []Node

	for node := l.head; node != nil; node = node.next {
		if equalValues(node.value, n) {
			nodes = append(nodes, *node)
		}
	}

	return nodes
}

// 3. Удаление одного узла по его значению (all = false).
// 4. Удаление всех узлов по конкретному значению (all = true).
// Время: O(n)
// Память: O(1)
func (l *LinkedList2) Delete(n int, all bool) {
	node := l.head

	for node != nil {
		next := node.next

		if !equalValues(node.value, n) {
			node = next
			continue
		}

		// Наличие prev избавляет от отдельного поиска предыдущего узла:
		// достаточно связать соседей удаляемого узла напрямую.
		if node.prev == nil {
			l.head = next
		} else {
			node.prev.next = next
		}

		if next == nil {
			l.tail = node.prev
		} else {
			next.prev = node.prev
		}

		node.prev = nil
		node.next = nil
		node = next

		if !all {
			return
		}
	}
}

// 5. Вставка нового узла после заданного узла.
// Если after == nil, узел становится первым элементом списка.
// Время: O(1)
// Память: O(1)
func (l *LinkedList2) Insert(after *Node, add Node) {
	if after == nil {
		l.InsertFirst(add)
		return
	}

	add.prev = after
	add.next = after.next

	if after.next == nil {
		l.tail = &add
	} else {
		after.next.prev = &add
	}

	after.next = &add
}

// 6. Вставка нового узла самым первым элементом списка.
// Время: O(1)
// Память: O(1)
func (l *LinkedList2) InsertFirst(first Node) {
	first.prev = nil
	first.next = l.head

	if l.head == nil {
		l.tail = &first
	} else {
		l.head.prev = &first
	}

	l.head = &first
}

// 7. Очистка всего содержимого списка.
// Время: O(1)
// Память: O(1)
func (l *LinkedList2) Clean() {
	l.head = nil
	l.tail = nil
}

// Вычисление длины списка.
// Время: O(n)
// Память: O(1)
func (l *LinkedList2) Count() int {
	count := 0

	for node := l.head; node != nil; node = node.next {
		count++
	}

	return count
}
