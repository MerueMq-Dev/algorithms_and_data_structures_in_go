package main

import (
	"os"
	"strconv"
)

// Задание 12. Кэш.
const (
	hashBase  = 31
	cacheStep = 1
)

type NativeCache[T any] struct {
	HashTable
	values []T
	hits   []int
}

func NewNativeCache[T any](sz int) *NativeCache[T] {
	sz = max(sz, 1)

	return &NativeCache[T]{
		HashTable: InitHashTable(sz, cacheStep),
		values:    make([]T, sz),
		hits:      make([]int, sz),
	}
}

// Новый ключ в полном кэше вытесняет ключ с наименьшим числом обращений.
// Время: O(n)
// Память: O(1)
func (nc *NativeCache[T]) Put(key string, value T) {
	index := nc.Find(key)

	if index < 0 {
		if nc.SeekSlot(key) < 0 {
			nc.evict()
		}

		index = nc.HashTable.Put(key)
		nc.hits[index] = 0
	}

	nc.values[index] = value
}

// Каждое чтение увеличивает счётчик обращений.
// Время: O(1) в среднем
// Память: O(1)
func (nc *NativeCache[T]) Get(key string) (T, error) {
	var result T

	index := nc.Find(key)
	if index < 0 {
		return result, os.ErrNotExist
	}

	nc.hits[index]++

	return nc.values[index], nil
}

// Время: O(1) в среднем
// Память: O(1)
func (nc *NativeCache[T]) IsKey(key string) bool {
	return nc.Find(key) >= 0
}

// Время: O(1) в среднем
// Память: O(1)
func (nc *NativeCache[T]) Hits(key string) (int, error) {
	index := nc.Find(key)
	if index < 0 {
		return 0, os.ErrNotExist
	}

	return nc.hits[index], nil
}

// Освобождает слот с наименьшим числом обращений.
// Время: O(n)
// Память: O(1)
func (nc *NativeCache[T]) evict() {
	victim := 0

	for i := range nc.hits {
		if nc.hits[i] < nc.hits[victim] {
			victim = i
		}
	}

	var empty T
	nc.slots[victim] = ""
	nc.values[victim] = empty
	nc.hits[victim] = 0
	nc.filled[victim] = false
}

// Хэш-таблица из задания 8.
type HashTable struct {
	size   int
	step   int
	slots  []string
	filled []bool
}

func InitHashTable(sz int, stp int) HashTable {
	ht := HashTable{size: sz, step: stp, slots: nil}
	ht.slots = make([]string, sz)
	ht.filled = make([]bool, sz)
	return ht
}

// Время: O(k), где k — длина строки
// Память: O(1)
func (ht *HashTable) HashFun(value string) int {
	if number, err := strconv.Atoi(value); err == nil {
		return (number%ht.size + ht.size) % ht.size
	}

	hash := 0

	for i := 0; i < len(value); i++ {
		hash = (hash*hashBase + int(value[i])) % ht.size
	}

	return hash
}

// Время: O(1) в среднем
// Память: O(1)
func (ht *HashTable) SeekSlot(value string) int {
	return slotOrMinus(ht.seekFrom(ht.HashFun(value), ht.step))
}

// Время: O(1) в среднем
// Память: O(1)
func (ht *HashTable) Put(value string) int {
	return ht.putFrom(ht.HashFun(value), ht.step, value)
}

// Время: O(1) в среднем
// Память: O(1)
func (ht *HashTable) Find(value string) int {
	return slotOrMinus(ht.findFrom(ht.HashFun(value), ht.step, value))
}

func (ht *HashTable) seekFrom(start int, step int) (int, error) {
	for i := 0; i < ht.size; i++ {
		index := (start + i*step) % ht.size

		if !ht.filled[index] {
			return index, nil
		}
	}

	return -1, os.ErrNotExist
}

func (ht *HashTable) findFrom(start int, step int, value string) (int, error) {
	for i := 0; i < ht.size; i++ {
		index := (start + i*step) % ht.size

		if !ht.filled[index] {
			break
		}

		if ht.slots[index] == value {
			return index, nil
		}
	}

	return -1, os.ErrNotExist
}

func (ht *HashTable) putFrom(start int, step int, value string) int {
	index, err := ht.seekFrom(start, step)
	if err != nil {
		return -1
	}

	ht.slots[index] = value
	ht.filled[index] = true

	return index
}

func slotOrMinus(index int, err error) int {
	if err != nil {
		return -1
	}

	return index
}
