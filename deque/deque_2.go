package main

import "os"

// Задание 6. Двусторонняя очередь.

// 4*. Проверка палиндрома: снимаем символы с обоих концов и сравниваем.
// Время: O(n)
// Память: O(n)
func IsPalindrome(s string) bool {
	var deque Deque[rune]

	for _, symbol := range s {
		deque.AddTail(symbol)
	}

	for deque.Size() > 1 {
		first, _ := deque.RemoveFront()
		last, _ := deque.RemoveTail()

		if first != last {
			return false
		}
	}

	return true
}

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

// Стек, где рядом с элементом хранится минимум всех элементов под ним.
type minStack[T ordered] struct {
	items []minEntry[T]
}

type minEntry[T ordered] struct {
	value   T
	minimum T
}

func (s *minStack[T]) size() int {
	return len(s.items)
}

func (s *minStack[T]) push(itm T) {
	minimum := itm

	if last := len(s.items) - 1; last >= 0 && s.items[last].minimum < minimum {
		minimum = s.items[last].minimum
	}

	s.items = append(s.items, minEntry[T]{value: itm, minimum: minimum})
}

func (s *minStack[T]) pop() T {
	last := len(s.items) - 1
	result := s.items[last].value

	s.items[last] = minEntry[T]{}
	s.items = s.items[:last]

	return result
}

func (s *minStack[T]) min() T {
	return s.items[len(s.items)-1].minimum
}

// 5*. Дека с минимумом за O(1) на двух стеках.
// Вершина front — голова деки, вершина back — хвост.
type MinDeque[T ordered] struct {
	front minStack[T]
	back  minStack[T]
}

func (d *MinDeque[T]) Size() int {
	return d.front.size() + d.back.size()
}

// Время: O(1)
// Память: O(1)
func (d *MinDeque[T]) AddFront(itm T) {
	d.front.push(itm)
}

// Время: O(1)
// Память: O(1)
func (d *MinDeque[T]) AddTail(itm T) {
	d.back.push(itm)
}

// Время: O(1) амортизированно
// Память: O(1) амортизированно
func (d *MinDeque[T]) RemoveFront() (T, error) {
	var result T

	if d.Size() == 0 {
		return result, os.ErrNotExist
	}

	if d.front.size() == 0 {
		moveHalf(&d.back, &d.front)
	}

	return d.front.pop(), nil
}

// Время: O(1) амортизированно
// Память: O(1) амортизированно
func (d *MinDeque[T]) RemoveTail() (T, error) {
	var result T

	if d.Size() == 0 {
		return result, os.ErrNotExist
	}

	if d.back.size() == 0 {
		moveHalf(&d.front, &d.back)
	}

	return d.back.pop(), nil
}

// Время: O(1)
// Память: O(1)
func (d *MinDeque[T]) Min() (T, error) {
	var result T

	switch {
	case d.Size() == 0:
		return result, os.ErrNotExist
	case d.front.size() == 0:
		return d.back.min(), nil
	case d.back.size() == 0:
		return d.front.min(), nil
	}

	return min(d.front.min(), d.back.min()), nil
}

// Нижняя половина from переезжает в пустой to, дно from становится вершиной to.
// Время: O(n)
// Память: O(n)
func moveHalf[T ordered](from, to *minStack[T]) {
	values := make([]T, from.size())

	for i, entry := range from.items {
		values[i] = entry.value
	}

	half := (len(values) + 1) / 2
	from.items = nil

	for _, value := range values[half:] {
		from.push(value)
	}

	for i := half - 1; i >= 0; i-- {
		to.push(values[i])
	}
}

// 6*. Дека на динамическом массиве — кольцевой буфер.
const (
	minCapacity     = 16
	growFactor      = 2
	shrinkFactor    = 1.5
	shrinkThreshold = 0.5
)

type ArrayDeque[T any] struct {
	buffer []T
	head   int
	count  int
}

func (d *ArrayDeque[T]) Size() int {
	return d.count
}

// Время: O(1) амортизированно
// Память: O(1) амортизированно
func (d *ArrayDeque[T]) AddFront(itm T) {
	d.grow()

	d.head = (d.head - 1 + len(d.buffer)) % len(d.buffer)
	d.buffer[d.head] = itm
	d.count++
}

// Время: O(1) амортизированно
// Память: O(1) амортизированно
func (d *ArrayDeque[T]) AddTail(itm T) {
	d.grow()

	d.buffer[d.index(d.count)] = itm
	d.count++
}

// Время: O(1) амортизированно
// Память: O(1)
func (d *ArrayDeque[T]) RemoveFront() (T, error) {
	var result T

	if d.count == 0 {
		return result, os.ErrNotExist
	}

	result = d.buffer[d.head]

	var empty T
	d.buffer[d.head] = empty

	d.head = d.index(1)
	d.count--
	d.shrink()

	return result, nil
}

// Время: O(1) амортизированно
// Память: O(1)
func (d *ArrayDeque[T]) RemoveTail() (T, error) {
	var result T

	if d.count == 0 {
		return result, os.ErrNotExist
	}

	last := d.index(d.count - 1)
	result = d.buffer[last]

	var empty T
	d.buffer[last] = empty

	d.count--
	d.shrink()

	return result, nil
}

func (d *ArrayDeque[T]) index(offset int) int {
	return (d.head + offset) % len(d.buffer)
}

// Время: O(n)
// Память: O(n)
func (d *ArrayDeque[T]) resize(capacity int) {
	buffer := make([]T, capacity)

	for i := 0; i < d.count; i++ {
		buffer[i] = d.buffer[d.index(i)]
	}

	d.buffer = buffer
	d.head = 0
}

func (d *ArrayDeque[T]) grow() {
	if d.buffer == nil {
		d.buffer = make([]T, minCapacity)
		return
	}

	if d.count == len(d.buffer) {
		d.resize(len(d.buffer) * growFactor)
	}
}

// Сжатие, если буфер заполнен меньше чем наполовину.
// Время: O(n)
// Память: O(n)
func (d *ArrayDeque[T]) shrink() {
	capacity := len(d.buffer)

	if capacity <= minCapacity {
		return
	}

	if float64(d.count)/float64(capacity) >= shrinkThreshold {
		return
	}

	d.resize(max(int(float64(capacity)/shrinkFactor), minCapacity))
}

// Рефлексия по заданию 4.
//
// 4. Скобки одного типа. Тут всё совпало, и это было ожидаемо: открывающую в
// стек, на закрывающую снимаю, а если снимать нечего — баланс нарушен. В конце
// стек должен остаться пустым. Честно, для одного типа меня подмывало обойтись
// счётчиком. Но тогда 5-я задача превратилась бы в переписывание с нуля, а со
// стеком она выросла из этой почти без правок.
//
// 5. Скобки трёх типов. Словарь у меня есть, только смотрит в другую сторону:
// ключ — закрывающая, значение — открывающая. В эталоне наоборот. Мне мой
// вариант по-прежнему нравится: снял открывающую и сразу сравнил с
// pairs[symbol]. А вот чего я не заметил, пока не прочитал про сто типов скобок:
// они у меня перечислены ещё и в case. Добавишь новую пару в словарь — и надо не
// забыть поправить switch. Надо было спрашивать у самого словаря, что перед
// тобой — открывающая или закрывающая.
//
// 6. Минимум за O(1). Вот тут разошлись по-настоящему. Я кладу минимум во
// второй стек на каждом Push, и он растёт вровень с основным. Эталон экономнее:
// кладёт элемент, только если тот меньше или равен текущему минимуму, а
// снимает, только если снятый с ним совпал. Если минимум обновляется редко,
// второй стек у эталона почти пустой, а у меня — копия основного. В худшем
// случае оба O(n), так что выигрыш эталона только в среднем, но он есть. Зато у
// меня нечего сломать: Pop снимает с обоих стеков без всяких сравнений. В
// эталоне легко промахнуться, написать < вместо <= и получить сломанный стек на
// первом же повторе минимума. В 5* этого задания я остался при своей схеме и
// даже её ужал: минимум лежит прямо в паре с элементом, и при перекладывании
// половины стека они переезжают вместе.
//
// 7. Среднее за O(1). Сошлось: сумма под рукой, Push прибавляет, Pop вычитает,
// среднее — одно деление. А потом я посмотрел на тип суммы, и стало неуютно.
// Она у меня того же типа T, что и элементы. Для AvgStack[int8] хватит двух
// элементов по 100, чтобы сумма переполнилась и среднее стало мусором. С float
// беда тише: после длинной серии Push/Pop в сумме копится ошибка округления.
// Сумму стоило держать в int64 или float64, не привязывая к типу элементов.
//
// 8. Постфиксная запись. Ловушку с a.pop() - b.pop() в одной строке я обошёл,
// но заслуги тут мало — выручил Go. Pop возвращает ещё и ошибку, так что
// операнды хочешь не хочешь снимаешь по одному. Кстати, порядок вызовов в
// выражении Go как раз гарантирует — слева направо. Но right и left с именами
// всё равно читаются лучше, чем расчёт на спецификацию. Где я проиграл, так это
// в switch вместо словаря с функциями. Дело не только в том, что каждая новая
// операция — это ещё один case. Сравнивая с эталоном, я нашёл у себя настоящий
// баг: неизвестный токен распознаётся только после снятия двух операндов. На
// "5 x" функция жалуется на нехватку операндов, хотя проблема в том, что x — не
// операция. Со словарём такого бы не случилось: сначала ищем токен, нет его —
// сразу ошибка, есть — снимаем операнды и вызываем функцию.
