package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Задание 4. Стек.
// 4*. Проверка баланса круглых скобок.
// Время: O(n)
// Память: O(n)
func IsBalanced(expression string) bool {
	var stack Stack[rune]

	for _, symbol := range expression {
		switch symbol {
		case '(':
			stack.Push(symbol)
		case ')':
			if _, err := stack.Pop(); err != nil {
				return false
			}
		}
	}

	return stack.Size() == 0
}

// 5*. Проверка баланса скобок трёх типов: (), {}, [].
// Время: O(n)
// Память: O(n)
func IsBalancedAll(expression string) bool {
	pairs := map[rune]rune{
		')': '(',
		']': '[',
		'}': '{',
	}

	var stack Stack[rune]

	for _, symbol := range expression {
		switch symbol {
		case '(', '[', '{':
			stack.Push(symbol)
		case ')', ']', '}':
			opened, err := stack.Pop()

			// Закрывающая должна совпасть по типу с последней открытой.
			if err != nil || opened != pairs[symbol] {
				return false
			}
		}
	}

	return stack.Size() == 0
}

type ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 | ~string
}

// 6*. Стек с минимальным элементом за O(1).
// На вершине второго стека всегда лежит минимум основного.
type MinStack[T ordered] struct {
	items Stack[T]
	mins  Stack[T]
}

func (st *MinStack[T]) Size() int {
	return st.items.Size()
}

func (st *MinStack[T]) Peek() (T, error) {
	return st.items.Peek()
}

// Во второй стек кладём меньшее из нового элемента и текущего минимума.
// Время: O(1)
// Память: O(1)
func (st *MinStack[T]) Push(itm T) {
	minimum := itm

	if current, err := st.mins.Peek(); err == nil && current < minimum {
		minimum = current
	}

	st.items.Push(itm)
	st.mins.Push(minimum)
}

// Снимаем сразу с обоих стеков.
// Время: O(1)
// Память: O(1)
func (st *MinStack[T]) Pop() (T, error) {
	result, err := st.items.Pop()
	if err != nil {
		return result, err
	}

	st.mins.Pop()

	return result, nil
}

// Минимум — вершина второго стека.
// Время: O(1)
// Память: O(1)
func (st *MinStack[T]) Min() (T, error) {
	return st.mins.Peek()
}

type number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// 7*. Стек со средним значением за O(1).
// Сумму держим готовой, среднее — это деление.
type AvgStack[T number] struct {
	items Stack[T]
	sum   T
}

func (st *AvgStack[T]) Size() int {
	return st.items.Size()
}

func (st *AvgStack[T]) Peek() (T, error) {
	return st.items.Peek()
}

// Время: O(1)
// Память: O(1)
func (st *AvgStack[T]) Push(itm T) {
	st.items.Push(itm)
	st.sum += itm
}

// Время: O(1)
// Память: O(1)
func (st *AvgStack[T]) Pop() (T, error) {
	result, err := st.items.Pop()
	if err != nil {
		return result, err
	}

	st.sum -= result

	return result, nil
}

// Среднее всех элементов стека.
// Время: O(1)
// Память: O(1)
func (st *AvgStack[T]) Average() (float64, error) {
	if st.items.Size() == 0 {
		return 0, fmt.Errorf("stack is empty")
	}

	return float64(st.sum) / float64(st.items.Size()), nil
}

// 8*. Вычисление постфиксного выражения двумя стеками.
// Первый стек — токены выражения (верхушка слева), второй — числа.
// Операция снимает два числа со второго стека и кладёт результат обратно.
// Время: O(n)
// Память: O(n)
func CalcPostfix(expression string) (int, error) {
	tokens := strings.Fields(expression)

	var input Stack[string]
	var values Stack[int]

	for i := len(tokens) - 1; i >= 0; i-- {
		input.Push(tokens[i])
	}

	for input.Size() > 0 {
		token, _ := input.Pop()

		if token == "=" {
			break
		}

		if number, err := strconv.Atoi(token); err == nil {
			values.Push(number)
			continue
		}

		right, err := values.Pop()
		if err != nil {
			return 0, fmt.Errorf("not enough operands for '%s'", token)
		}

		left, err := values.Pop()
		if err != nil {
			return 0, fmt.Errorf("not enough operands for '%s'", token)
		}

		switch token {
		case "+":
			values.Push(left + right)
		case "*":
			values.Push(left * right)
		case "-":
			values.Push(left - right)
		default:
			return 0, fmt.Errorf("unknown token '%s'", token)
		}
	}

	if values.Size() != 1 {
		return 0, fmt.Errorf("bad expression '%s'", expression)
	}

	return values.Pop()
}

// Рефлексия по заданию 2.
//
// 9. Разворот списка. Сделал как в эталоне — меняю указатели у каждого узла, в
// конце меняю местами head и tail. Отдельная переменная под предыдущий узел не
// понадобилась, prev и так лежит в узле, обмен занял одну строку. Зато потерял
// время на другом: после обмена следующий узел оказывается в prev, а не в next.
// Пока не сообразил, цикл разворачивал только первый узел и на этом заканчивался.
//
// 10. Проверка циклов. Эталон советует пройти столько шагов, сколько элементов в
// списке, и посмотреть, попали ли в хвост. Но длину брать неоткуда. Count() сам
// обходит список, а на зациклённом он просто не вернётся. Получается, чтобы найти
// цикл, надо сначала по нему пройти. Такой вариант работал бы, только если длина
// хранится отдельным полем и всегда актуальна. В итоге взял Флойда — ему длина не
// нужна вообще. Думал ещё про множество посещённых узлов, но это O(n) памяти там,
// где хватает двух указателей.
//
// 11. Сортировка. Пузырёк на списке выглядит естественно: соседи рядом, значения
// меняются на месте. Я всё равно взял слияние, потому что пузырёк — это O(n^2), и
// на тысяче узлов это уже чувствуется. Заодно оказалось, что слияние на списке
// ложится даже удобнее, чем на массиве: разрезать пополам двумя указателями и
// склеить обратно можно вообще без произвольного доступа, которого у списка нет.
// Плата — рекурсия глубиной log n. Связи prev во время сортировки не поддерживаю.
// Сначала пытался чинить их по ходу слияния, потом бросил: дешевле пройти один
// раз в конце.
//
// 12. Слияние. Тут сделал меньше, чем в эталоне: слил два списка, а не
// произвольное количество с выбором минимального из текущих. Перебором это
// O(n*k). Заодно понял, зачем дальше идёт куча — она отдаёт минимум за O(log k).
// Мой вариант остаётся частным случаем k = 2, переписать его логично уже после
// кучи. Запрет вызывать сортировку в результате поначалу казался придиркой, а по
// сути это O((n+m) log(n+m)) вместо честного O(n+m).
//
// 13. Dummy. Главное поймал: краевые случаи исчезают, потому что у каждого
// видимого узла всегда есть соседи слева и справа, и вставка с удалением
// становятся безусловными. После возни с head и tail в основном списке разница
// заметная. А вот фиктивные узлы отличаю не так, как советует эталон — сравниваю
// указатель с head и tail вместо отдельного класса-наследника. Флажок в Node не
// добавлял, так что ошибки, о которой предупреждает эталон, у меня не возникло.
// Проверка типа всё же выразительнее сравнения по указателю, и в Go её можно
// получить через интерфейс, раз наследования нет. Ещё отметил для себя про
// круговой список: замкнёшь хвост на голову — и одного фиктивного узла хватает
// сразу на обе роли.
