package main

import (
	"os"
	"reflect"
)

// Задание 1. Связный (связанный) список.
type Node struct {
	next  *Node
	value int
}

type LinkedList struct {
	head *Node
	tail *Node
}

func (l *LinkedList) AddInTail(item Node) {
	if l.head == nil {
		l.head = &item
	} else {
		l.tail.next = &item
	}

	l.tail = &item
}

func (l *LinkedList) InsertFirst(first Node) {
	first.next = l.head
	l.head = &first

	if l.tail == nil {
		l.tail = &first
	}
}

// 1. Удаление одного узла по его значению.
// Время: O(n)
// Память: O(1)
// 2. Удаление всех узлов по конкретному значению.
// Время: O(n)
// Память: O(1)
func (l *LinkedList) Delete(n int, all bool) {
	var prev *Node
	node := l.head

	for node != nil {
		if node.value != n {
			prev = node
			node = node.next
			continue
		}

		next := node.next

		if prev == nil {
			l.head = next
		} else {
			prev.next = next
		}

		if next == nil {
			l.tail = prev
		}

		node.next = nil
		node = next

		if !all {
			return
		}
	}
}

// 3. Очистка всего содержимого списка.
// Время: O(1)
// Память: O(1)
func (l *LinkedList) Clean() {
	l.head = nil
	l.tail = nil
}

// 4. Поиск всех узлов по конкретному значению.
// Время: O(n)
// Память: O(k), где k — количество найденных узлов.
func (l *LinkedList) FindAll(n int) []Node {
	var nodes []Node

	for node := l.head; node != nil; node = node.next {
		if node.value == n {
			nodes = append(nodes, *node)
		}
	}

	return nodes
}

// 5. Вычисление длины списка.
// Время: O(n)
// Память: O(1)
func (l *LinkedList) Count() int {
	count := 0

	for node := l.head; node != nil; node = node.next {
		count++
	}

	return count
}

// 6. Вставка нового узла после заданного узла.
// Время: O(1)
// Память: O(1)
func (l *LinkedList) Insert(after *Node, add Node) {
	if after == nil {
		return
	}

	add.next = after.next
	after.next = &add

	if after == l.tail {
		l.tail = &add
	}
}

// Поиск первого узла по его значению.
// Время: O(n)
// Память: O(1)
func (l *LinkedList) Find(n int) (Node, error) {
	for node := l.head; node != nil; node = node.next {
		if node.value == n {
			return *node, nil
		}
	}

	return Node{value: -1, next: nil}, os.ErrNotExist
}
