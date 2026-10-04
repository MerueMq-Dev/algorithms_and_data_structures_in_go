package main

import (
	"constraints"
	"os"
)

// Задание 7. Упорядоченный список.

// 8*. Удаление всех дубликатов: равные значения стоят рядом.
// Время: O(n)
// Память: O(1)
func (l *OrderedList[T]) RemoveDuplicates() {
	node := l.head

	for node != nil && node.next != nil {
		if l.Compare(node.value, node.next.value) == 0 {
			l.unlink(node.next)
		} else {
			node = node.next
		}
	}
}

// 9*. Слияние в новый список с порядком первого: каждый список обходим
// в нужную сторону и забираем меньший из двух текущих элементов.
// Время: O(n + m)
// Память: O(n + m)
func Merge[T constraints.Ordered](l1 *OrderedList[T], l2 *OrderedList[T]) *OrderedList[T] {
	result := &OrderedList[T]{}
	result.Clear(l1._ascending)

	node1, step1 := l1.start(result._ascending)
	node2, step2 := l2.start(result._ascending)

	for node1 != nil || node2 != nil {
		if node2 == nil || node1 != nil && result.order(node1.value, node2.value) <= 0 {
			result.insertBefore(nil, &Node[T]{value: node1.value})
			node1 = step1(node1)
		} else {
			result.insertBefore(nil, &Node[T]{value: node2.value})
			node2 = step2(node2)
		}
	}

	return result
}

func (l *OrderedList[T]) start(asc bool) (*Node[T], func(*Node[T]) *Node[T]) {
	if asc == l._ascending {
		return l.head, func(node *Node[T]) *Node[T] { return node.next }
	}

	return l.tail, func(node *Node[T]) *Node[T] { return node.prev }
}

// 10*. Есть ли sub в списке подряд идущими элементами.
// Начало ищем только среди узлов, равных первому значению sub.
// Время: O(n * m)
// Память: O(1)
func (l *OrderedList[T]) HasSublist(sub *OrderedList[T]) bool {
	first, step := sub.start(l._ascending)

	if first == nil {
		return true
	}

	for node := l.head; node != nil; node = node.next {
		c := l.order(node.value, first.value)

		if c > 0 {
			return false
		}

		if c == 0 && l.matchesFrom(node, first, step) {
			return true
		}
	}

	return false
}

func (l *OrderedList[T]) matchesFrom(node *Node[T], subNode *Node[T], step func(*Node[T]) *Node[T]) bool {
	for ; subNode != nil; subNode = step(subNode) {
		if node == nil || l.Compare(node.value, subNode.value) != 0 {
			return false
		}

		node = node.next
	}

	return true
}

// 11*. Самое частое значение: равные стоят рядом, считаем длину серий.
// При равных частотах — то, что в списке раньше.
// Время: O(n)
// Память: O(1)
func (l *OrderedList[T]) MostFrequent() (T, error) {
	var best T

	if l.head == nil {
		return best, os.ErrNotExist
	}

	bestCount := 0
	runCount := 0

	for node := l.head; node != nil; node = node.next {
		if node.prev != nil && l.Compare(node.prev.value, node.value) == 0 {
			runCount++
		} else {
			runCount = 1
		}

		if runCount > bestCount {
			best = node.value
			bestCount = runCount
		}
	}

	return best, nil
}

// 12*. Упорядоченный список на слайсе: индекс ищется двоичным поиском.
type OrderedArray[T constraints.Ordered] struct {
	items     []T
	ascending bool
}

func NewOrderedArray[T constraints.Ordered](asc bool) *OrderedArray[T] {
	return &OrderedArray[T]{ascending: asc}
}

func (a *OrderedArray[T]) Count() int {
	return len(a.items)
}

// Время: O(n)
// Память: O(1) амортизированно
func (a *OrderedArray[T]) Add(item T) {
	i := a.lowerBound(item)

	var empty T
	a.items = append(a.items, empty)

	copy(a.items[i+1:], a.items[i:])
	a.items[i] = item
}

// Время: O(n)
// Память: O(1)
func (a *OrderedArray[T]) Delete(n T) {
	i, err := a.Index(n)
	if err != nil {
		return
	}

	copy(a.items[i:], a.items[i+1:])

	var empty T
	a.items[len(a.items)-1] = empty
	a.items = a.items[:len(a.items)-1]
}

// Индекс первого элемента со значением n.
// Время: O(log n)
// Память: O(1)
func (a *OrderedArray[T]) Index(n T) (int, error) {
	i := a.lowerBound(n)

	if i == len(a.items) || compareValues(a.items[i], n) != 0 {
		return -1, os.ErrNotExist
	}

	return i, nil
}

// Первая позиция, где элемент не должен стоять раньше n.
// Время: O(log n)
// Память: O(1)
func (a *OrderedArray[T]) lowerBound(n T) int {
	low := 0
	high := len(a.items)

	for low < high {
		mid := low + (high-low)/2

		if a.order(a.items[mid], n) < 0 {
			low = mid + 1
		} else {
			high = mid
		}
	}

	return low
}

func (a *OrderedArray[T]) order(x T, y T) int {
	if a.ascending {
		return compareValues(x, y)
	}

	return compareValues(y, x)
}

// Рефлексия по заданию 5.
//
// 3. Вращение очереди. Сделал так же, как в эталоне: снимаю с головы, кладу в
// хвост, и так n раз. От себя добавил только остаток от деления на размер и
// отрицательный сдвиг. Просто представил, что кто-то попросит повернуть очередь
// из трёх элементов на миллион шагов, и стало жалко крутить её впустую.
//
// 4. Очередь на двух стеках. Тут я угадал полностью, и, честно, было приятно.
// Кладу во входной стек, беру из выходного, переливаю, только когда выходной
// пуст. В эталоне сказано «близко к O(1)». Я бы назвал это амортизированным
// O(1): каждый элемент переливается всего один раз. Немного расстроило другое.
// Стек из прошлого задания я так и не переиспользовал: каждая папка у меня —
// отдельный package main, и подключить его оттуда нельзя. Пришлось написать
// маленький стек заново прямо рядом с очередью. Работает, но осадок остался —
// это же копипаста.
//
// 5. Разворот через стек. Здесь вообще без открытий: всё из очереди в стек,
// потом обратно. Стек сам переворачивает порядок, ничего придумывать не надо.
//
// 6. Кольцевой буфер. Самый интересный пункт, потому что мы с эталоном пошли
// разными дорогами. В эталоне два указателя, head и tail, и одна ячейка всегда
// остаётся пустой, чтобы отличить «пусто» от «полно». Я храню head и счётчик, а
// хвост считаю как (head + count) % размер. Тогда пусто — это count == 0,
// полно — count == размеру, и ни одна ячейка не пропадает. Сначала я даже
// засомневался, не упростил ли я что-то лишнее. Но эталон сам упоминает
// счётчик как вариант, так что всё честно. Зато у меня очередь на 3 элемента
// действительно вмещает 3, а в эталоне — только 2. Ещё я отдельно проверил
// тестом нулевой размер. Деления на ноль там не будет: проверки на пустоту и
// заполненность срабатывают раньше, чем дело доходит до остатка.
