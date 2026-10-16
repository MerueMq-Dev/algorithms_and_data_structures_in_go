package main

import (
	"errors"
	"os"
)

// Задание 11. Фильтр Блюма.

// 2*. Слияние фильтров.
// Ложных срабатываний становится больше: P = 0,6931^(m/N), где N — сколько
// значений во всех фильтрах вместе.
// Время: O(f), где f — число фильтров
// Память: O(1)
func Merge(filters ...*BloomFilter) (*BloomFilter, error) {
	if len(filters) == 0 {
		return nil, errors.New("no filters to merge")
	}

	result, err := NewBloomFilter(filters[0].filter_len)
	if err != nil {
		return nil, err
	}

	for _, filter := range filters {
		if filter.filter_len != result.filter_len {
			return nil, errors.New("filters have different sizes")
		}

		result.bits |= filter.bits
	}

	return result, nil
}

// 3*. Фильтр с удалением.
const maxCounter = 255

type CountingBloomFilter struct {
	filter   BloomFilter
	counters []uint8
}

func NewCountingBloomFilter(filterLen int) (*CountingBloomFilter, error) {
	filter, err := NewBloomFilter(filterLen)
	if err != nil {
		return nil, err
	}

	return &CountingBloomFilter{filter: *filter, counters: make([]uint8, filterLen)}, nil
}

// Время: O(k), где k — длина строки
// Память: O(1)
func (cf *CountingBloomFilter) Add(s string) {
	for _, i := range bitPositions(&cf.filter, s) {
		if cf.counters[i] < maxCounter {
			cf.counters[i]++
		}
	}
}

// Время: O(k), где k — длина строки
// Память: O(1)
func (cf *CountingBloomFilter) IsValue(s string) bool {
	for _, i := range bitPositions(&cf.filter, s) {
		if cf.counters[i] == 0 {
			return false
		}
	}

	return true
}

// Время: O(k), где k — длина строки
// Память: O(1)
func (cf *CountingBloomFilter) Remove(s string) error {
	if !cf.IsValue(s) {
		return os.ErrNotExist
	}

	for _, i := range bitPositions(&cf.filter, s) {
		if cf.counters[i] < maxCounter {
			cf.counters[i]--
		}
	}

	return nil
}

// 4*. Восстановление исходного множества.
// Если какой-то бит ставит только одна строка, она точно была добавлена.
// Время: O(c * k), где c — число кандидатов
// Память: O(c + m)
func Recover(bf *BloomFilter, candidates []string) (sure []string, possible []string, err error) {
	owners := make([]int, bf.filter_len)

	for _, candidate := range candidates {
		if bf.IsValue(candidate) {
			possible = append(possible, candidate)

			for _, i := range bitPositions(bf, candidate) {
				owners[i]++
			}
		}
	}

	for i := 0; i < bf.filter_len; i++ {
		if bf.getBit(i) && owners[i] == 0 {
			return nil, possible, errors.New("candidates do not cover the filter")
		}
	}

	for _, candidate := range possible {
		for _, i := range bitPositions(bf, candidate) {
			if owners[i] == 1 {
				sure = append(sure, candidate)
				break
			}
		}
	}

	return sure, possible, nil
}

func bitPositions(bf *BloomFilter, s string) []int {
	first := bf.Hash1(s)
	second := bf.Hash2(s)

	if first == second {
		return []int{first}
	}

	return []int{first, second}
}

// Рефлексия по заданию 9.
//
// 5. Словарь на упорядоченном списке. Эталон хранит значения отдельно, в
// конце обычного списка, а в упорядоченном списке держит пары [ключ, индекс
// значения]. У меня проще: ключ и значение лежат вместе в одной записи, и
// записи упорядочены по ключу. Индексы не нужно поддерживать в порядке, и при
// удалении в списке значений не остаётся дыр. Зато я не использовал свой
// OrderedList из задания 7, а эталон предлагает именно это, только с другим
// сравнением. Тут я заметил одну вещь: мой OrderedList сравнивает строки без
// пробелов по краям. Если взять его как есть, ключи "a" и " a" стали бы одним
// ключом. Так что сравнение там и правда пришлось бы менять. С минусами,
// которые называет эталон, я согласен: мой словарь тоже не потокобезопасен, а
// вставка стоит O(n) из-за сдвига. Дерево здесь подошло бы лучше.
