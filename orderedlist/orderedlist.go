package main

import (
	"constraints"
	"os"
	"strings"
)

// Задание 7. Упорядоченный список.
type Node[T constraints.Ordered] struct {
	prev  *Node[T]
	next  *Node[T]
	value T
}

type OrderedList[T constraints.Ordered] struct {
	head       *Node[T]
	tail       *Node[T]
	_ascending bool
	count      int
}

// Время: O(1)
// Память: O(1)
func (l *OrderedList[T]) Count() int {
	return l.count
}

// 3. Вставка перед первым узлом, который должен стоять после нового.
// Время: O(n)
// Память: O(1)
func (l *OrderedList[T]) Add(item T) {
	node := l.head

	for node != nil && l.order(node.value, item) <= 0 {
		node = node.next
	}

	l.insertBefore(node, &Node[T]{value: item})
}

// 6. Поиск с ранней остановкой.
// Время: O(n)
// Память: O(1)
func (l *OrderedList[T]) Find(n T) (Node[T], error) {
	node := l.find(n)

	if node == nil {
		return Node[T]{}, os.ErrNotExist
	}

	return *node, nil
}

// 4. Удаление первого узла с заданным значением.
// Время: O(n)
// Память: O(1)
func (l *OrderedList[T]) Delete(n T) {
	if node := l.find(n); node != nil {
		l.unlink(node)
	}
}

// 1. Очистка и выбор порядка: asc — по возрастанию, иначе по убыванию.
// Время: O(1)
// Память: O(1)
func (l *OrderedList[T]) Clear(asc bool) {
	l.head = nil
	l.tail = nil
	l.count = 0
	l._ascending = asc
}

// 2, 5. -1, если v1 < v2; 0, если равны; +1, если v1 > v2.
// Время: O(1), для строк — O(k), где k — длина строк
// Память: O(1)
func (l *OrderedList[T]) Compare(v1 T, v2 T) int {
	return compareValues(v1, v2)
}

func compareValues[T constraints.Ordered](v1 T, v2 T) int {
	if s1, ok := any(v1).(string); ok {
		s2 := any(v2).(string)

		return strings.Compare(strings.TrimSpace(s1), strings.TrimSpace(s2))
	}

	if v1 < v2 {
		return -1
	}

	if v1 > v2 {
		return +1
	}

	return 0
}

// Больше нуля, если a должно стоять в списке после b.
func (l *OrderedList[T]) order(a T, b T) int {
	if l._ascending {
		return l.Compare(a, b)
	}

	return l.Compare(b, a)
}

func (l *OrderedList[T]) find(n T) *Node[T] {
	for node := l.head; node != nil; node = node.next {
		c := l.order(node.value, n)

		if c == 0 {
			return node
		}

		if c > 0 {
			return nil
		}
	}

	return nil
}

func (l *OrderedList[T]) insertBefore(before *Node[T], node *Node[T]) {
	node.next = before

	if before == nil {
		node.prev = l.tail
		l.tail = node
	} else {
		node.prev = before.prev
		before.prev = node
	}

	if node.prev == nil {
		l.head = node
	} else {
		node.prev.next = node
	}

	l.count++
}

func (l *OrderedList[T]) unlink(node *Node[T]) {
	if node.prev == nil {
		l.head = node.next
	} else {
		node.prev.next = node.next
	}

	if node.next == nil {
		l.tail = node.prev
	} else {
		node.next.prev = node.prev
	}

	node.prev = nil
	node.next = nil
	l.count--
}

// 6. Сложность поиска.
//
// Не изменилась: O(n). Ранняя остановка экономит проход, когда значения нет, но
// если искомое в конце списка или должно стоять после всех, обходим всё.
