package main

import (
	"hash/maphash"
	"strconv"
)

// Задание 8. Хэш-таблица.

// 3*. Таблица, которая растёт при заполнении больше чем на 75%.
const maxLoad = 0.75

type DynamicHashTable struct {
	table HashTable
	count int
}

func NewDynamicHashTable(sz int, stp int) *DynamicHashTable {
	return &DynamicHashTable{table: Init(sz, stp)}
}

func (d *DynamicHashTable) Size() int {
	return d.table.size
}

func (d *DynamicHashTable) Count() int {
	return d.count
}

// Время: O(1) амортизированно
// Память: O(1) амортизированно
func (d *DynamicHashTable) Put(value string) int {
	if float64(d.count+1) > float64(d.table.size)*maxLoad {
		d.grow()
	}

	index := d.table.Put(value)

	if index < 0 {
		d.grow()
		index = d.table.Put(value)
	}

	if index >= 0 {
		d.count++
	}

	return index
}

// Время: O(1) в среднем
// Память: O(1)
func (d *DynamicHashTable) Find(value string) int {
	return d.table.Find(value)
}

// Новая таблица простого размера не меньше чем вдвое больше, все значения
// раскладываются заново.
// Время: O(n)
// Память: O(n)
func (d *DynamicHashTable) grow() {
	bigger := Init(nextPrime(max(d.table.size*2, 2)), d.table.step)

	for i, filled := range d.table.filled {
		if filled {
			bigger.Put(d.table.slots[i])
		}
	}

	d.table = bigger
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

// 4*. Двойное хэширование: первая функция даёт слот, вторая — шаг.
// Значения, попавшие в один слот, расходятся разными шагами, и длинные цепочки
// не собираются. Вероятность совпадения первого слота та же, зато проб меньше.
// Цена — вторая хэш-функция на каждую операцию.
const stepBase = 37

type DoubleHashTable struct {
	table HashTable
}

func NewDoubleHashTable(sz int) *DoubleHashTable {
	return &DoubleHashTable{table: Init(sz, 1)}
}

// Шаг от 1 до size - 1: при простом размере обходит все слоты.
// Время: O(k), где k — длина строки
// Память: O(1)
func (d *DoubleHashTable) StepFun(value string) int {
	if d.table.size < 2 {
		return 1
	}

	hash := 0

	for i := 0; i < len(value); i++ {
		hash = (hash*stepBase + int(value[i])) % (d.table.size - 1)
	}

	return hash + 1
}

// Время: O(1) в среднем
// Память: O(1)
func (d *DoubleHashTable) Put(value string) int {
	return d.table.putFrom(d.table.HashFun(value), d.StepFun(value), value)
}

// Время: O(1) в среднем
// Память: O(1)
func (d *DoubleHashTable) Find(value string) int {
	return slotOrMinus(d.table.findFrom(d.table.HashFun(value), d.StepFun(value), value))
}

// 5*. Атака: подбираем ключи, которые попадают в один слот.
// Каждая следующая вставка обходит всех предыдущих — O(n^2) на n ключей.
// Время: O(count * size)
// Память: O(count)
func CollidingKeys(ht *HashTable, count int) []string {
	var keys []string

	for i := 0; len(keys) < count; i++ {
		key := strconv.Itoa(i)

		if ht.HashFun(key) == 0 {
			keys = append(keys, key)
		}
	}

	return keys
}

// 5*. Защита: хэш со случайной солью, которая выбирается при создании таблицы.
// Приписать соль к строке не помогает: у полиномиального хэша строки одной
// длины так и останутся в одном слоте.
type SaltedHashTable struct {
	table HashTable
	salt  maphash.Seed
}

func NewSaltedHashTable(sz int, stp int) *SaltedHashTable {
	return &SaltedHashTable{table: Init(sz, stp), salt: maphash.MakeSeed()}
}

// Время: O(k), где k — длина строки
// Память: O(1)
func (st *SaltedHashTable) HashFun(value string) int {
	return int(maphash.String(st.salt, value) % uint64(st.table.size))
}

// Время: O(1) в среднем
// Память: O(1)
func (st *SaltedHashTable) Put(value string) int {
	return st.table.putFrom(st.HashFun(value), st.table.step, value)
}

// Время: O(1) в среднем
// Память: O(1)
func (st *SaltedHashTable) Find(value string) int {
	return slotOrMinus(st.table.findFrom(st.HashFun(value), st.table.step, value))
}

// Рефлексия по заданию 6.
//
// 4. Палиндром. Сделал так же, как в эталоне: кладу строку в деку и снимаю по
// символу с обоих концов, пока они совпадают. Тут всё сошлось.
//
// 5. Минимум деки. Эталон предлагает вторую деку, где элементы лежат по
// возрастанию. Я написал её по описанию и проверил на случайных данных. Для
// обычной очереди она работает, если выкидывать только строго бо́льшие
// элементы. С «больше или равно» минимум теряется на повторах. А когда
// элементы снимают с той же стороны, куда добавляли, минимум теряется в
// большинстве случаев. Например: AddTail(3), AddFront(1), RemoveFront — вторая дека пустая,
// а в основной лежит 3. По-моему, этот способ хорош для очереди, а для деки
// его надо дорабатывать. Я сделал по-другому: два стека с минимумами, один для
// головы, другой для хвоста. Когда один пустеет, перекладываю в него
// половину другого. Так работают все четыре операции. Минус в том, что удаление
// у меня O(1) только в среднем: иногда приходится перекладывать половину деки.
//
// 6. Дека на динамическом массиве. Тут эталон прав, а я ошибся. Он советует не
// смешивать в одном классе деку и динамический массив. У меня в ArrayDeque всё
// вместе: и кольцевой буфер, и рост, и сжатие. Константы роста я вообще
// скопировал из DynArray. Лучше было взять готовый динамический массив и
// работать с ним только как с декой. В этом задании я так и сделал: таблицы в
// 3*–5* не копируют HashTable, а хранят её внутри и меняют только то, как
// считается слот или шаг.
