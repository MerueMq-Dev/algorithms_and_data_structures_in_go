package main

import (
	"errors"
	"os"
)

// Задание 5. Очередь.

// 3*. Вращение очереди на n элементов: n раз переносим голову в хвост.
// При отрицательном n вращаем в другую сторону.
// Время: O(n mod size)
// Память: O(1)
func RotateQueue[T any](q *Queue[T], n int) {
	size := q.Size()
	if size == 0 {
		return
	}

	n %= size
	if n < 0 {
		n += size
	}

	for i := 0; i < n; i++ {
		item, _ := q.Dequeue()
		q.Enqueue(item)
	}
}

// Простой стек на слайсе для 4* и 5*.
type stack[T any] struct {
	items []T
}

func (s *stack[T]) push(itm T) {
	s.items = append(s.items, itm)
}

func (s *stack[T]) pop() (T, bool) {
	var result T

	if len(s.items) == 0 {
		return result, false
	}

	last := len(s.items) - 1
	result = s.items[last]

	var empty T
	s.items[last] = empty
	s.items = s.items[:last]

	return result, true
}

func (s *stack[T]) size() int {
	return len(s.items)
}

// 4*. Очередь на двух стеках.
type StackQueue[T any] struct {
	in  stack[T]
	out stack[T]
}

func (q *StackQueue[T]) Size() int {
	return q.in.size() + q.out.size()
}

// Время: O(1)
// Память: O(1)
func (q *StackQueue[T]) Enqueue(itm T) {
	q.in.push(itm)
}

// Время: O(1) амортизированно: каждый элемент переливается один раз,
// O(n) в момент переливания.
// Память: O(1)
func (q *StackQueue[T]) Dequeue() (T, error) {
	if q.out.size() == 0 {
		for q.in.size() > 0 {
			item, _ := q.in.pop()
			q.out.push(item)
		}
	}

	item, ok := q.out.pop()
	if !ok {
		return item, os.ErrNotExist
	}

	return item, nil
}

// 5*. Разворот очереди через стек.
// Время: O(n)
// Память: O(n)
func ReverseQueue[T any](q *Queue[T]) {
	var buffer stack[T]

	for q.Size() > 0 {
		item, _ := q.Dequeue()
		buffer.push(item)
	}

	for buffer.size() > 0 {
		item, _ := buffer.pop()
		q.Enqueue(item)
	}
}

// 6*. Круговая очередь на массиве фиксированного размера.
type CircularQueue[T any] struct {
	buffer []T
	head   int
	count  int
}

var errQueueFull = errors.New("queue is full")

// Время: O(capacity)
// Память: O(capacity)
func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
	if capacity < 0 {
		capacity = 0
	}

	return &CircularQueue[T]{buffer: make([]T, capacity)}
}

func (q *CircularQueue[T]) Size() int {
	return q.count
}

// Проверка, полна ли очередь.
// Время: O(1)
// Память: O(1)
func (q *CircularQueue[T]) IsFull() bool {
	return q.count == len(q.buffer)
}

// Время: O(1)
// Память: O(1)
func (q *CircularQueue[T]) Enqueue(itm T) error {
	if q.IsFull() {
		return errQueueFull
	}

	tail := (q.head + q.count) % len(q.buffer)
	q.buffer[tail] = itm
	q.count++

	return nil
}

// Время: O(1)
// Память: O(1)
func (q *CircularQueue[T]) Dequeue() (T, error) {
	var result T

	if q.count == 0 {
		return result, os.ErrNotExist
	}

	result = q.buffer[q.head]

	var empty T
	q.buffer[q.head] = empty

	q.head = (q.head + 1) % len(q.buffer)
	q.count--

	return result, nil
}

// Рефлексия по заданию 3.
//
// 6. Банковский метод. С добавлением сошлось с эталоном — 3 монеты: одна на
// саму запись, две в банк. Реаллокацию тоже делаю, когда массив заполнен
// целиком, и списываю из банка N. Проверял это тестом: баланс ни разу не ушёл в
// минус за тысячу вставок, и внесённых монет хватило на все копирования. А вот
// про удаление я просто не подумал. BankDynArray переопределяет только Append,
// Remove достаётся от DynArray как есть, и банк про него вообще не знает. По
// эталону удаление стоит 2 монеты (1 + 1), а при сжатии можно ничего не
// списывать или списать 10% от N. Выходит, мой вариант честно считает только
// рост массива. Insert в середину тоже прошёл мимо банка.
//
// 7. Многомерный массив. Тут разошёлся с эталоном сильнее всего. Эталон держит
// всё в одном одномерном массиве: для 3x4x5 это 60 элементов, а координата
// (i, j, k) просто пересчитывается в индекс. Я сделал дерево из динамических
// массивов — хотел, чтобы каждое измерение росло отдельно и остальное не
// приходилось перекладывать. Сейчас вижу, что решал не ту проблему. Рост оси в
// плоском массиве — та же реаллокация, что и в обычном динамическом массиве, и
// амортизируется она так же, причём банковский метод из задачи 6 как раз это и
// показывает. Зато плоский вариант проще, лежит в памяти подряд и находит
// элемент одной арифметикой. У меня же на каждое измерение переход по указателю,
// плюс узлы создаются лениво, и из-за этого в Get пришлось отдельно разбирать
// ячейки, до которых запись ещё не доходила. В язык тут упёрся только синтаксис:
// записи myArr[1,2,3] в Go нет, поэтому доступ идёт через Get(1, 2, 3). На выбор
// хранения это не влияло — плоский слайс с пересчётом индекса пишется в Go так же
// просто.
