package main

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

// Полиномиальный хэш по байтам строки.
// Время: O(k), где k — длина строки
// Память: O(1)
func (ht *HashTable) HashFun(value string) int {
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
	return ht.seekFrom(ht.HashFun(value), ht.step)
}

// Время: O(1) в среднем
// Память: O(1)
func (ht *HashTable) Put(value string) int {
	return ht.putAt(ht.SeekSlot(value), value)
}

// Слот со значением или -1.
// Время: O(1) в среднем
// Память: O(1)
func (ht *HashTable) Find(value string) int {
	return ht.findFrom(ht.HashFun(value), ht.step, value)
}

func (ht *HashTable) seekFrom(start int, step int) int {
	for i := 0; i < ht.size; i++ {
		index := (start + i*step) % ht.size

		if !ht.filled[index] {
			return index
		}
	}

	return -1
}

func (ht *HashTable) findFrom(start int, step int, value string) int {
	for i := 0; i < ht.size; i++ {
		index := (start + i*step) % ht.size

		if !ht.filled[index] {
			return -1
		}

		if ht.slots[index] == value {
			return index
		}
	}

	return -1
}

func (ht *HashTable) putAt(index int, value string) int {
	if index < 0 {
		return -1
	}

	ht.slots[index] = value
	ht.filled[index] = true

	return index
}
