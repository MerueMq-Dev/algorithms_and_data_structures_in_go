package main

// Задание 2. Двусвязный (двунаправленный) список.

// 9*. Разворот порядка элементов списка на противоположный.
// Время: O(n)
// Память: O(1)
func (l *LinkedList2) Reverse() {
	node := l.head

	for node != nil {
		node.prev, node.next = node.next, node.prev
		node = node.prev
	}

	l.head, l.tail = l.tail, l.head
}

// 10*. Проверка наличия цикла внутри списка.
// Время: O(n)
// Память: O(1)
func (l *LinkedList2) HasCycle() bool {
	slow := l.head
	fast := l.head

	for fast != nil && fast.next != nil {
		slow = slow.next
		fast = fast.next.next

		if slow == fast {
			return true
		}
	}

	return false
}

// 11*. Сортировка списка по возрастанию значений.
// Время: O(n log n)
// Память: O(log n) — глубина рекурсии.
func (l *LinkedList2) Sort() {
	l.head = sortNodes(l.head)
	l.tail = restoreLinks(l.head)
}

// sortNodes сортирует цепочку узлов, опираясь только на поле next.
func sortNodes(head *Node) *Node {
	if head == nil || head.next == nil {
		return head
	}

	// Делим цепочку пополам двумя указателями.
	slow := head
	fast := head.next

	for fast != nil && fast.next != nil {
		slow = slow.next
		fast = fast.next.next
	}

	second := slow.next
	slow.next = nil

	return mergeNodes(sortNodes(head), sortNodes(second))
}

// mergeNodes сливает две отсортированные цепочки в одну.
func mergeNodes(first *Node, second *Node) *Node {
	var head *Node
	var tail *Node

	appendNode := func(node *Node) {
		if head == nil {
			head = node
		} else {
			tail.next = node
		}
		tail = node
	}

	for first != nil && second != nil {
		if first.value <= second.value {
			next := first.next
			appendNode(first)
			first = next
		} else {
			next := second.next
			appendNode(second)
			second = next
		}
	}

	rest := first
	if rest == nil {
		rest = second
	}

	for rest != nil {
		next := rest.next
		appendNode(rest)
		rest = next
	}

	return head
}

// restoreLinks проставляет prev по всей цепочке и возвращает последний узел.
func restoreLinks(head *Node) *Node {
	var prev *Node

	for node := head; node != nil; node = node.next {
		node.prev = prev
		prev = node
	}

	return prev
}

// 12*. Объединение двух отсортированных списков в третий, также отсортированный.
// Входные списки должны быть отсортированы заранее; метод сортировки здесь
// не используется — результат собирается слиянием за один проход.
// Исходные списки не изменяются: в новый список копируются значения.
// Время: O(n + m)
// Память: O(n + m)
func MergeSorted(list1 LinkedList2, list2 LinkedList2) LinkedList2 {
	var result LinkedList2

	first := list1.head
	second := list2.head

	for first != nil && second != nil {
		if first.value <= second.value {
			result.AddInTail(Node{value: first.value})
			first = first.next
		} else {
			result.AddInTail(Node{value: second.value})
			second = second.next
		}
	}

	for ; first != nil; first = first.next {
		result.AddInTail(Node{value: first.value})
	}

	for ; second != nil; second = second.next {
		result.AddInTail(Node{value: second.value})
	}

	return result
}

// 13*. Список с фиктивными (dummy) узлами.

type DummyLinkedList struct {
	head *Node
	tail *Node
}

// NewDummyLinkedList создаёт пустой список из двух фиктивных узлов.
// Время: O(1)
// Память: O(1)
func NewDummyLinkedList() *DummyLinkedList {
	head := &Node{}
	tail := &Node{}

	head.next = tail
	tail.prev = head

	return &DummyLinkedList{head: head, tail: tail}
}

// InsertBefore — единственная операция вставки: никаких проверок на пустой
// список, на голову и на хвост.
// Время: O(1)
// Память: O(1)
func (l *DummyLinkedList) InsertBefore(before *Node, add Node) *Node {
	node := &add

	node.prev = before.prev
	node.next = before
	before.prev.next = node
	before.prev = node

	return node
}

// AddInTail добавляет узел в конец — это вставка перед фиктивным хвостом.
// Время: O(1)
// Память: O(1)
func (l *DummyLinkedList) AddInTail(item Node) *Node {
	return l.InsertBefore(l.tail, item)
}

// InsertFirst добавляет узел в начало — вставка перед первым видимым узлом.
// Время: O(1)
// Память: O(1)
func (l *DummyLinkedList) InsertFirst(item Node) *Node {
	return l.InsertBefore(l.head.next, item)
}

// Remove — единственная операция удаления, тоже без краевых случаев.
// Время: O(1)
// Память: O(1)
func (l *DummyLinkedList) Remove(node *Node) {
	node.prev.next = node.next
	node.next.prev = node.prev
	node.prev = nil
	node.next = nil
}

// Delete удаляет узлы по значению, обходя только видимую часть списка.
// Время: O(n)
// Память: O(1)
func (l *DummyLinkedList) Delete(n int, all bool) {
	for node := l.head.next; node != l.tail; {
		next := node.next

		if equalValues(node.value, n) {
			l.Remove(node)

			if !all {
				return
			}
		}

		node = next
	}
}

// Values возвращает значения видимых узлов по порядку.
// Время: O(n)
// Память: O(n)
func (l *DummyLinkedList) Values() []int {
	var values []int

	for node := l.head.next; node != l.tail; node = node.next {
		values = append(values, node.value)
	}

	return values
}

// Count возвращает количество видимых узлов.
// Время: O(n)
// Память: O(1)
func (l *DummyLinkedList) Count() int {
	count := 0

	for node := l.head.next; node != l.tail; node = node.next {
		count++
	}

	return count
}
