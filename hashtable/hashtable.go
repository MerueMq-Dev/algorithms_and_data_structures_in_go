package main

import (
	"os"
	"strconv"
)

// Задание 8. Хэш-таблица.
const hashBase = 31

type HashTable struct {
	size   int
	step   int
	slots  []string
	filled []bool
}

func Init(sz int, stp int) HashTable {
	ht := HashTable{size: sz, step: stp, slots: nil}
	ht.slots = make([]string, sz)
	ht.filled = make([]bool, sz)
	return ht
}

// Числа хэшируются как числа, остальные строки — полиномиально по байтам.
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

// Пустой слот для значения или -1.
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

// Слот со значением или -1.
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
