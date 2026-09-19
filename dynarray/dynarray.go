package main

import (
	"fmt"
	"os"
)

// Задание 3. Динамический массив.
// Буфер растёт в 2 раза, сжимается делением на 1.5, минимум — 16 элементов.
const (
	minCapacity     = 16
	growFactor      = 2
	shrinkFactor    = 1.5
	shrinkThreshold = 0.5
)

type DynArray[T any] struct {
	count    int
	capacity int
	array    []T
}

func (da *DynArray[T]) Init() {
	da.count = 0
	da.MakeArray(minCapacity)
}

// 1. Новый буфер заданного размера, старые элементы копируются в него.
// Время: O(n)
// Память: O(sz)
func (da *DynArray[T]) MakeArray(sz int) {
	if sz < minCapacity {
		sz = minCapacity
	}

	var arr = make([]T, sz)

	n := da.count
	if n > sz {
		n = sz
	}

	copy(arr, da.array[:n])

	da.count = n
	da.capacity = sz
	da.array = arr
}

// 1. Элемент по индексу.
// Время: O(1)
// Память: O(1)
func (da *DynArray[T]) GetItem(index int) (T, error) {
	var result T

	if index < 0 || index >= da.count {
		return result, badIndex(index)
	}

	return da.array[index], nil
}

// 1. Добавление элемента в конец.
// Время: O(1) амортизированно, O(n) при расширении буфера
// Память: O(1) амортизированно
func (da *DynArray[T]) Append(itm T) {
	da.ensureInit()

	if da.count == da.capacity {
		da.MakeArray(da.capacity * growFactor)
	}

	da.array[da.count] = itm
	da.count++
}

// 2. Вставка в позицию index, остальные элементы сдвигаются вправо.
// Если index равен количеству элементов, элемент добавляется в хвост.
// Время: O(n) — сдвигаем хвост, при нехватке места копируем буфер.
// Вставка в хвост — O(1) амортизированно.
// Память: O(1), не считая нового буфера
func (da *DynArray[T]) Insert(itm T, index int) error {
	da.ensureInit()

	if index < 0 || index > da.count {
		return badIndex(index)
	}

	if da.count == da.capacity {
		da.MakeArray(da.capacity * growFactor)
	}

	copy(da.array[index+1:da.count+1], da.array[index:da.count])

	da.array[index] = itm
	da.count++

	return nil
}

// 3. Удаление элемента из позиции index, буфер при необходимости сжимается.
// 4. Время: O(n) — сдвигаем хвост, при сжатии копируем буфер.
// Удаление последнего элемента без сжатия — O(1).
// Память: O(1), не считая нового буфера
func (da *DynArray[T]) Remove(index int) error {
	if index < 0 || index >= da.count {
		return badIndex(index)
	}

	copy(da.array[index:da.count-1], da.array[index+1:da.count])

	// Чистим освободившуюся ячейку, иначе буфер продолжит держать ссылку на объект.
	var empty T
	da.array[da.count-1] = empty
	da.count--

	da.shrink()

	return nil
}

// Сжатие буфера, если он заполнен меньше чем наполовину.
// Время: O(n) при сжатии, иначе O(1)
// Память: O(1)
func (da *DynArray[T]) shrink() {
	if da.capacity <= minCapacity {
		return
	}

	if float64(da.count)/float64(da.capacity) >= shrinkThreshold {
		return
	}

	newCapacity := int(float64(da.capacity) / shrinkFactor)

	if newCapacity < minCapacity {
		newCapacity = minCapacity
	}

	if newCapacity < da.count {
		newCapacity = da.count
	}

	if newCapacity != da.capacity {
		da.MakeArray(newCapacity)
	}
}

func (da *DynArray[T]) ensureInit() {
	if da.array == nil {
		da.Init()
	}
}

func badIndex(index int) error {
	return fmt.Errorf("bad index '%d': %w", index, os.ErrInvalid)
}
