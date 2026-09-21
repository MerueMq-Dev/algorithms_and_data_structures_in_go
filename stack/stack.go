package main

import "os"

// Задание 4. Стек.
type Stack[T any] struct {
	items []T
}

// 1. Количество элементов в стеке.
// Время: O(1)
// Память: O(1)
func (st *Stack[T]) Size() int {
	return len(st.items)
}

// 1. Верхний элемент без снятия со стека.
// Время: O(1)
// Память: O(1)
func (st *Stack[T]) Peek() (T, error) {
	var result T

	if len(st.items) == 0 {
		return result, os.ErrNotExist
	}

	return st.items[len(st.items)-1], nil
}

// 1. Снятие верхнего элемента.
// Время: O(1)
// Память: O(1)
func (st *Stack[T]) Pop() (T, error) {
	var result T

	if len(st.items) == 0 {
		return result, os.ErrNotExist
	}

	last := len(st.items) - 1
	result = st.items[last]

	var empty T
	st.items[last] = empty
	st.items = st.items[:last]

	return result, nil
}

// 1. Добавление элемента на вершину.
// Время: O(1) амортизированно, O(n) при расширении слайса
// Память: O(1) амортизированно
func (st *Stack[T]) Push(itm T) {
	st.items = append(st.items, itm)
}

type ListStack[T any] struct {
	head  *stackNode[T]
	count int
}

type stackNode[T any] struct {
	next  *stackNode[T]
	value T
}

// 2. Количество элементов в стеке.
// Время: O(1)
// Память: O(1)
func (st *ListStack[T]) Size() int {
	return st.count
}

// 2. Верхний элемент без снятия со стека.
// Время: O(1)
// Память: O(1)
func (st *ListStack[T]) Peek() (T, error) {
	var result T

	if st.head == nil {
		return result, os.ErrNotExist
	}

	return st.head.value, nil
}

// 2. Снятие верхнего элемента: голова заменяется следующим узлом.
// Время: O(1)
// Память: O(1)
func (st *ListStack[T]) Pop() (T, error) {
	var result T

	if st.head == nil {
		return result, os.ErrNotExist
	}

	node := st.head
	result = node.value

	st.head = node.next
	node.next = nil
	st.count--

	return result, nil
}

// 2. Добавление элемента новой головой списка.
// Время: O(1)
// Память: O(1)
func (st *ListStack[T]) Push(itm T) {
	st.head = &stackNode[T]{next: st.head, value: itm}
	st.count++
}

// 3. Как отработает цикл
//
//	while (stack.size() > 0)
//	    stack.pop()
//	    stack.pop()
//
// Чётный размер — стек опустеет за size/2 итераций, цикл завершится.
// Нечётный — на последней итерации второй pop сработает на пустом стеке:
// в Go это ошибка, а если её не проверить — нулевое значение типа.
// Пустой стек — цикл не выполнится ни разу.
