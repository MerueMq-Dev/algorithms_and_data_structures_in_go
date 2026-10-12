package main

import (
	"os"
	"strconv"
)

// Задание 9. Ассоциативный массив.
const (
	hashBase = 31
	dictStep = 1
	maxLoad  = 0.75
)

type NativeDictionary[T any] struct {
	HashTable
	values []T
	count  int
}

func Init[T any](sz int) NativeDictionary[T] {
	sz = max(sz, 1)
	nd := NativeDictionary[T]{HashTable: InitHashTable(sz, dictStep), values: nil}
	nd.values = make([]T, sz)
	return nd
}

// 3. Время: O(1) в среднем
// Память: O(1)
func (nd *NativeDictionary[T]) IsKey(key string) bool {
	return nd.Find(key) >= 0
}

// 3. Время: O(1) в среднем
// Память: O(1)
func (nd *NativeDictionary[T]) Get(key string) (T, error) {
	var result T

	index := nd.Find(key)
	if index < 0 {
		return result, os.ErrNotExist
	}

	return nd.values[index], nil
}

// 3. Новый ключ добавляется, у старого меняется значение.
// Время: O(1) амортизированно
// Память: O(1) амортизированно
func (nd *NativeDictionary[T]) Put(key string, value T) {
	index := nd.Find(key)

	if index < 0 {
		if float64(nd.count+1) > float64(nd.size)*maxLoad {
			nd.grow()
		}

		index = nd.HashTable.Put(key)
		nd.count++
	}

	nd.values[index] = value
}

// Время: O(n)
// Память: O(n)
func (nd *NativeDictionary[T]) grow() {
	bigger := Init[T](nextPrime(nd.size * 2))

	for i, filled := range nd.filled {
		if filled {
			bigger.Put(nd.slots[i], nd.values[i])
		}
	}

	*nd = bigger
}

func nextPrime(n int) int {
	for !isPrime(n) {
		n++
	}

	return n
}

func isPrime(n int) bool {
	if n < 2 {
		return false
	}

	for d := 2; d*d <= n; d++ {
		if n%d == 0 {
			return false
		}
	}

	return true
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

// 1. Как устроен map в Go.
//
// С Go 1.24 map — это Swiss table, хэш-таблица с открытой адресацией. Слоты
// собраны в группы по 8, и у каждой группы хранятся короткие метки из хэшей
// ключей. По ним сразу видно, в каких слотах искать, и ключи сравниваются
// только там. Большой map разбит на таблицы до 1024 слотов, и при росте
// делится только переполненная таблица. До 1.24 map был массивом бакетов по
// 8 пар с цепочками переполнения. У каждого map своя случайная соль в хэше:
// это защищает от подобранных коллизий, и поэтому порядок обхода каждый раз
// разный. Писать в map из нескольких горутин одновременно нельзя — программа
// упадёт с "concurrent map writes". Для этого нужен sync.Mutex или sync.Map.
