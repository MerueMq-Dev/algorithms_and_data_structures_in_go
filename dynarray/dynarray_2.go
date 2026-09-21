package main

import "fmt"

// Задание 3. Динамический массив.
// 6*. Динамический массив на основе банковского метода.
const appendCost = 3

type BankDynArray[T any] struct {
	DynArray[T]

	bank        int
	totalPaid   int
	totalCopied int
}

func NewBankDynArray[T any]() *BankDynArray[T] {
	ba := &BankDynArray[T]{}
	ba.Init()
	return ba
}

// Добавление элемента в конец со счётом монет.
// Время: O(1) амортизированно
// Память: O(1) амортизированно
func (ba *BankDynArray[T]) Append(itm T) {
	ba.ensureInit()

	ba.bank += appendCost
	ba.totalPaid += appendCost
	ba.bank--

	if ba.count == ba.capacity {
		ba.bank -= ba.count
		ba.totalCopied += ba.count
	}

	ba.DynArray.Append(itm)
}

func (ba *BankDynArray[T]) Bank() int {
	return ba.bank
}

func (ba *BankDynArray[T]) TotalPaid() int {
	return ba.totalPaid
}

func (ba *BankDynArray[T]) TotalCopied() int {
	return ba.totalCopied
}

// 7*. Многомерный динамический массив.
type MultiArray[T any] struct {
	dims []int
	root *multiNode[T]
}

type multiNode[T any] struct {
	children *DynArray[*multiNode[T]]
	values   *DynArray[T]
}

func newMultiNode[T any](levelsLeft int) *multiNode[T] {
	node := &multiNode[T]{}

	if levelsLeft <= 1 {
		node.values = &DynArray[T]{}
		node.values.Init()
	} else {
		node.children = &DynArray[*multiNode[T]]{}
		node.children.Init()
	}

	return node
}

// Новый массив: число измерений и размер по каждому из них.
// Время: O(1)
// Память: O(1)
func NewMultiArray[T any](dims ...int) (*MultiArray[T], error) {
	if len(dims) == 0 {
		return nil, fmt.Errorf("bad dimensions count '%d'", len(dims))
	}

	for i, size := range dims {
		if size < 0 {
			return nil, fmt.Errorf("bad dimension size '%d' at %d", size, i)
		}
	}

	ma := &MultiArray[T]{
		dims: append([]int(nil), dims...),
		root: newMultiNode[T](len(dims)),
	}

	return ma, nil
}

func (ma *MultiArray[T]) Dims() []int {
	return append([]int(nil), ma.dims...)
}

// Элемент по набору индексов.
// Время: O(d), где d — число измерений
// Память: O(1)
func (ma *MultiArray[T]) Get(indices ...int) (T, error) {
	var result T

	if err := ma.checkIndices(indices); err != nil {
		return result, err
	}

	node := ma.root

	for level, index := range indices {
		if level == len(ma.dims)-1 {
			if index >= node.values.count {
				return result, nil
			}

			return node.values.GetItem(index)
		}

		if index >= node.children.count {
			return result, nil
		}

		child, err := node.children.GetItem(index)
		if err != nil {
			return result, err
		}

		if child == nil {
			return result, nil
		}

		node = child
	}

	return result, nil
}

// Запись элемента по набору индексов.
// Индекс больше текущей границы расширяет измерение.
// Время: O(d + s), где s — число новых ячеек
// Память: O(s)
func (ma *MultiArray[T]) Set(item T, indices ...int) error {
	if len(indices) != len(ma.dims) {
		return fmt.Errorf("bad indices count '%d'", len(indices))
	}

	for i, index := range indices {
		if index < 0 {
			return badIndex(index)
		}

		if index >= ma.dims[i] {
			ma.dims[i] = index + 1
		}
	}

	node := ma.root

	for level, index := range indices {
		if level == len(ma.dims)-1 {
			var empty T
			for node.values.count <= index {
				node.values.Append(empty)
			}

			node.values.array[index] = item

			return nil
		}

		for node.children.count <= index {
			node.children.Append(nil)
		}

		child, err := node.children.GetItem(index)
		if err != nil {
			return err
		}

		if child == nil {
			child = newMultiNode[T](len(ma.dims) - level - 1)
			node.children.array[index] = child
		}

		node = child
	}

	return nil
}

func (ma *MultiArray[T]) checkIndices(indices []int) error {
	if len(indices) != len(ma.dims) {
		return fmt.Errorf("bad indices count '%d'", len(indices))
	}

	for i, index := range indices {
		if index < 0 || index >= ma.dims[i] {
			return badIndex(index)
		}
	}

	return nil
}

// Рефлексия по заданию 1.
//
// 8. Сложение двух списков. Сошлось с эталоном один в один: сверил длины, пошёл
// по обоим спискам сразу, суммы — в хвост. А вот от совета проверять в цикле
// только первый список меня сначала передёрнуло: второй-то идёт вслепую! Потом
// дошло — он и держится на той самой сверке длин парой строк выше. Красиво, но
// хрупко: уберёшь проверку — и словишь nil на первом же неравном входе.
// Эталон ругает исключения, а в Go их попросту нет — возвращай ошибку и не
// выдумывай. И мне так спокойнее: разная длина списков это не катастрофа, а
// обычный вход, который бывает.
// Что сделал бы иначе: два Count() — это два лишних прохода просто чтобы узнать
// то, что указатели выяснят сами. Достаточно идти вдвоём и поймать момент, когда
// один список кончился раньше.
