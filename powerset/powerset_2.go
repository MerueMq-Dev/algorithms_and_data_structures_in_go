package main

import (
	"cmp"
	"constraints"
	"errors"
	"slices"
)

// Задание 10. Множество.

// 4*. Декартово произведение.
type Pair[T constraints.Ordered] struct {
	First  T
	Second T
}

// Время: O(n * m)
// Память: O(n * m)
func (s *PowerSet[T]) Cartesian(set2 PowerSet[T]) []Pair[T] {
	result := make([]Pair[T], 0, s.count*set2.count)

	for _, first := range s.values() {
		for _, second := range set2.values() {
			result = append(result, Pair[T]{First: first, Second: second})
		}
	}

	return result
}

// 5*. Пересечение трёх и более множеств: идём по самому маленькому.
// Время: O(k * min), где k — число множеств
// Память: O(min)
func IntersectAll[T constraints.Ordered](sets []PowerSet[T]) (PowerSet[T], error) {
	if len(sets) < 3 {
		return Init[T](), errors.New("need at least 3 sets")
	}

	smallest := 0

	for i := range sets {
		if sets[i].count < sets[smallest].count {
			smallest = i
		}
	}

	result := Init[T]()

	for _, value := range sets[smallest].values() {
		inAll := true

		for i := range sets {
			if i != smallest && !sets[i].Get(value) {
				inAll = false
				break
			}
		}

		if inAll {
			result.Put(value)
		}
	}

	return result, nil
}

// 6*. Мультимножество: рядом со значением хранится, сколько раз оно есть.
type Bag[T constraints.Ordered] struct {
	buckets [][]Frequency[T]
	size    int
}

type Frequency[T constraints.Ordered] struct {
	Value T
	Count int
}

func NewBag[T constraints.Ordered]() *Bag[T] {
	return &Bag[T]{buckets: make([][]Frequency[T], setSize)}
}

// Общее число экземпляров.
func (b *Bag[T]) Size() int {
	return b.size
}

// Время: O(1) в среднем
// Память: O(1) амортизированно
func (b *Bag[T]) Add(value T) {
	b.size++

	if entry := b.entry(value); entry != nil {
		entry.Count++
		return
	}

	i := bucketIndex(value)
	b.buckets[i] = append(b.buckets[i], Frequency[T]{Value: value, Count: 1})
}

// Удаление одного экземпляра.
// Время: O(1) в среднем
// Память: O(1)
func (b *Bag[T]) Remove(value T) bool {
	i := bucketIndex(value)
	bucket := b.buckets[i]

	for j := range bucket {
		if bucket[j].Value != value {
			continue
		}

		b.size--
		bucket[j].Count--

		if bucket[j].Count == 0 {
			last := len(bucket) - 1
			bucket[j] = bucket[last]
			bucket[last] = Frequency[T]{}
			b.buckets[i] = bucket[:last]
		}

		return true
	}

	return false
}

// Время: O(1) в среднем
// Память: O(1)
func (b *Bag[T]) Count(value T) int {
	if entry := b.entry(value); entry != nil {
		return entry.Count
	}

	return 0
}

// Все значения с частотами, по возрастанию значения.
// Время: O(n log n)
// Память: O(n)
func (b *Bag[T]) Frequencies() []Frequency[T] {
	var result []Frequency[T]

	for _, bucket := range b.buckets {
		result = append(result, bucket...)
	}

	slices.SortFunc(result, func(x, y Frequency[T]) int {
		return cmp.Compare(x.Value, y.Value)
	})

	return result
}

func (b *Bag[T]) entry(value T) *Frequency[T] {
	bucket := b.buckets[bucketIndex(value)]

	for j := range bucket {
		if bucket[j].Value == value {
			return &bucket[j]
		}
	}

	return nil
}

// Рефлексия по заданию 8.
//
// 3. Динамическая хэш-таблица. Композицию я использовал, как и советует
// эталон: новая таблица хранит обычную HashTable внутри. А вот динамический
// массив я не взял. Когда таблица заполнена больше чем на 75%, я создаю новую
// таблицу примерно вдвое больше и перекладываю туда все значения. Размер
// выбираю простым числом. Так шаг поиска гарантированно обходит все слоты.
// С обычным удвоением массива это могло бы сломаться: например, при размере
// 18 и шаге 3 поиск видел бы только треть слотов. Явного метода расширения,
// как в эталоне, у меня нет — таблица растёт сама, когда это нужно.
//
// 5. Соль. Согласен с эталоном, что статическая соль — плохая идея: её видно в
// коде, и под неё легко подобрать ключи. Совет добавлять к хэшу небольшое
// случайное число я проверил на 300 ключах, которые попадают в один слот. Без
// соли вставка стоила 45 150 проб. С добавкой от 0 до 3 — 18 853, а поиск
// вырос до 74 944: соль заранее неизвестна, и приходится пробовать все
// варианты. У меня соль — случайный seed для
// каждой таблицы. На тех же ключах он дал 371 пробу и на вставке, и на поиске.
