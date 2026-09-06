package main

import (
	"errors"
)

type Node struct {
	next  *Node
	value int
}

type LinkedList struct {
	head *Node
	tail *Node
}

func (l *LinkedList) AddInTail(item Node) {
	item.next = nil

	if l.head == nil {
		l.head = &item
		l.tail = &item
		return
	}

	l.tail.next = &item
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
func (l *LinkedList) Delete(n int, all bool) {
	var prev *Node
	current := l.head

	for current != nil {
		if current.value == n {
			if prev == nil {
				l.head = current.next
			} else {
				prev.next = current.next
			}

			if current == l.tail {
				l.tail = prev
			}

			if !all {
				return
			}

			current = current.next
			continue
		}

		prev = current
		current = current.next
	}

	if l.head == nil {
		l.tail = nil
	}
}

// 2. Удаление всех узлов по конкретному значению.
// Время: O(n)
// Память: O(1)
func (l *LinkedList) DeleteAll(n int) {
	var prev *Node
	current := l.head

	for current != nil {
		if current.value == n {
			if prev == nil {
				l.head = current.next
			} else {
				prev.next = current.next
			}

			if current == l.tail {
				l.tail = prev
			}

			current = current.next
			continue
		}

		prev = current
		current = current.next
	}

	if l.head == nil {
		l.tail = nil
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
	nodes := []Node{}
	current := l.head

	for current != nil {
		if current.value == n {
			nodes = append(nodes, *current)
		}

		current = current.next
	}

	return nodes
}

// 5. Вычисление длины списка.
// Время: O(n)
// Память: O(1)
func (l *LinkedList) Count() int {
	current := l.head
	counter := 0

	for current != nil {
		counter++
		current = current.next
	}

	return counter
}

// 6. Вставка нового узла после заданного узла.
// Время: O(1)
// Память: O(1)
func (l *LinkedList) Insert(after *Node, add Node) {
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
	current := l.head

	for current != nil {
		if current.value == n {
			return *current, nil
		}

		current = current.next
	}

	return Node{}, errors.New("node not found")
}
