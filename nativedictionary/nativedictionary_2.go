package main

import (
	"errors"
	"os"
	"strconv"
)

// Задание 9. Ассоциативный массив.

// 5*. Словарь на упорядоченном списке.
// Поиск — O(log n), вставка и удаление — O(n).
type OrderedDictionary[T any] struct {
	entries []entry[T]
}

type entry[T any] struct {
	key   string
	value T
}

func (od *OrderedDictionary[T]) Count() int {
	return len(od.entries)
}

// Время: O(log n)
// Память: O(1)
func (od *OrderedDictionary[T]) IsKey(key string) bool {
	_, found := od.search(key)
	return found
}

// Время: O(log n)
// Память: O(1)
func (od *OrderedDictionary[T]) Get(key string) (T, error) {
	var result T

	i, found := od.search(key)
	if !found {
		return result, os.ErrNotExist
	}

	return od.entries[i].value, nil
}

// Время: O(n)
// Память: O(1) амортизированно
func (od *OrderedDictionary[T]) Put(key string, value T) {
	i, found := od.search(key)

	if found {
		od.entries[i].value = value
		return
	}

	od.entries = append(od.entries, entry[T]{})
	copy(od.entries[i+1:], od.entries[i:])
	od.entries[i] = entry[T]{key: key, value: value}
}

// Время: O(n)
// Память: O(1)
func (od *OrderedDictionary[T]) Delete(key string) error {
	i, found := od.search(key)
	if !found {
		return os.ErrNotExist
	}

	copy(od.entries[i:], od.entries[i+1:])
	od.entries[len(od.entries)-1] = entry[T]{}
	od.entries = od.entries[:len(od.entries)-1]

	return nil
}

// Позиция ключа или место, куда его вставить.
// Время: O(log n)
// Память: O(1)
func (od *OrderedDictionary[T]) search(key string) (int, bool) {
	low := 0
	high := len(od.entries)

	for low < high {
		mid := low + (high-low)/2

		if od.entries[mid].key < key {
			low = mid + 1
		} else {
			high = mid
		}
	}

	return low, low < len(od.entries) && od.entries[low].key == key
}

// 6*. Словарь с битовыми ключами.
type BitDictionary[T any] struct {
	width int
	root  *bitNode[T]
	count int
}

type bitNode[T any] struct {
	children [2]*bitNode[T]
	value    T
	filled   bool
}

var errBadKey = errors.New("bad bit key")

func NewBitDictionary[T any](width int) (*BitDictionary[T], error) {
	if width < 1 || width > 64 {
		return nil, errBadKey
	}

	return &BitDictionary[T]{width: width, root: &bitNode[T]{}}, nil
}

func (bd *BitDictionary[T]) Count() int {
	return bd.count
}

// Время: O(w), где w — длина ключа
// Память: O(w)
func (bd *BitDictionary[T]) Put(key string, value T) error {
	bits, err := bd.parse(key)
	if err != nil {
		return err
	}

	node := bd.root

	for level := 0; level < bd.width; level++ {
		b := bd.bit(bits, level)

		if node.children[b] == nil {
			node.children[b] = &bitNode[T]{}
		}

		node = node.children[b]
	}

	if !node.filled {
		bd.count++
	}

	node.value = value
	node.filled = true

	return nil
}

// Время: O(w), где w — длина ключа
// Память: O(1)
func (bd *BitDictionary[T]) Get(key string) (T, error) {
	var result T

	node, err := bd.leaf(key)
	if err != nil {
		return result, err
	}

	return node.value, nil
}

// Время: O(w), где w — длина ключа
// Память: O(1)
func (bd *BitDictionary[T]) IsKey(key string) bool {
	_, err := bd.leaf(key)
	return err == nil
}

// Пустые ветки после удаления срезаются.
// Время: O(w), где w — длина ключа
// Память: O(w)
func (bd *BitDictionary[T]) Delete(key string) error {
	bits, err := bd.parse(key)
	if err != nil {
		return err
	}

	path := []*bitNode[T]{bd.root}

	for level := 0; level < bd.width; level++ {
		next := path[level].children[bd.bit(bits, level)]

		if next == nil {
			return os.ErrNotExist
		}

		path = append(path, next)
	}

	last := path[bd.width]

	if !last.filled {
		return os.ErrNotExist
	}

	var empty T
	last.value = empty
	last.filled = false
	bd.count--

	for level := bd.width - 1; level >= 0; level-- {
		child := path[level+1]

		if child.filled || child.children[0] != nil || child.children[1] != nil {
			break
		}

		path[level].children[bd.bit(bits, level)] = nil
	}

	return nil
}

func (bd *BitDictionary[T]) leaf(key string) (*bitNode[T], error) {
	bits, err := bd.parse(key)
	if err != nil {
		return nil, err
	}

	node := bd.root

	for level := 0; level < bd.width && node != nil; level++ {
		node = node.children[bd.bit(bits, level)]
	}

	if node == nil || !node.filled {
		return nil, os.ErrNotExist
	}

	return node, nil
}

func (bd *BitDictionary[T]) parse(key string) (uint64, error) {
	if len(key) != bd.width {
		return 0, errBadKey
	}

	bits, err := strconv.ParseUint(key, 2, 64)
	if err != nil {
		return 0, errBadKey
	}

	return bits, nil
}

// Бит ключа на уровне level, считая от старшего.
func (bd *BitDictionary[T]) bit(bits uint64, level int) int {
	return int(bits >> (bd.width - 1 - level) & 1)
}

// Рефлексия по заданию 7.
//
// 9. Слияние. Сделал как в эталоне: иду по двум спискам и каждый раз беру
// меньший элемент. Сортировать заранее ничего не нужно, списки и так
// упорядочены. У меня есть одна добавка: списки могут быть упорядочены в
// разные стороны. Тогда второй список я прохожу с хвоста. В рекомендации об
// этом не сказано, но без этого слияние выдало бы неправильный порядок.
//
// 10. Под-список. Тут я сначала проверил рекомендацию кодом: написал её на
// слайсах и сравнил с простым перебором на 100 000 случайных списков. Ошибок
// не нашлось ни у эталона, ни у меня. Но у эталона больше способов закончить
// поиск раньше. Он останавливается, когда до конца списка осталось меньше
// элементов, чем в под-списке. И ещё одно правило, до которого я не додумался:
// если при несовпадении элемент под-списка меньше элемента основного списка,
// дальше искать бессмысленно. У меня поиск заканчивается, только когда
// значения в списке стали больше первого элемента под-списка. Ответ тот же,
// но в некоторых случаях я делаю лишние проверки.
//
// 11. Самое частое значение. Совпало: один проход и сравнение с предыдущим
// элементом. Небольшая разница в счётчике. Эталон при смене значения
// сбрасывает его в 0, а я в 1. Получается, у эталона счётчик считает повторы,
// а у меня — сами элементы. На результат это не влияет. Ещё я решил, что при
// равных частотах возвращаю значение, которое встретилось раньше. В эталоне
// этот случай не описан.
//
// 12. Индекс за O(log n). Сделал, как советует эталон: двоичный поиск по
// массиву. Разница есть с повторами. Эталон возвращает индекс, как только
// попал на равный элемент, и это может быть любой из одинаковых. Я ищу левую
// границу, поэтому всегда возвращаю первый. Что мне не нравится в моём
// решении: я сделал отдельный тип OrderedArray, а не поменял устройство
// OrderedList. Вышло две реализации одного списка, и вставка в массив стоит
// O(n) из-за сдвига.
