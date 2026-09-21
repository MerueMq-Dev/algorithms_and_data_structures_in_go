package main

import "testing"

// Задание 4. Стек.

func makeStack(values ...int) *Stack[int] {
	stack := &Stack[int]{}

	for _, value := range values {
		stack.Push(value)
	}

	return stack
}

func assertStack(t *testing.T, stack *Stack[int], expected []int) {
	t.Helper()

	if stack.Size() != len(expected) {
		t.Fatalf("expected size %d, got %d", len(expected), stack.Size())
	}

	// Снимаем элементы и сверяем с ожидаемым порядком, от вершины вниз.
	for i := len(expected) - 1; i >= 0; i-- {
		got, err := stack.Pop()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != expected[i] {
			t.Errorf("position %d: expected %d, got %d", i, expected[i], got)
		}
	}

	if stack.Size() != 0 {
		t.Errorf("expected empty stack, got size %d", stack.Size())
	}
}

// 1. Size

func TestSizeEmptyStack(t *testing.T) {
	stack := makeStack()

	if stack.Size() != 0 {
		t.Errorf("expected 0, got %d", stack.Size())
	}
}

func TestSizeAfterPush(t *testing.T) {
	stack := makeStack(10, 20, 30)

	if stack.Size() != 3 {
		t.Errorf("expected 3, got %d", stack.Size())
	}
}

func TestSizeAfterPop(t *testing.T) {
	stack := makeStack(10, 20, 30)

	stack.Pop()
	stack.Pop()

	if stack.Size() != 1 {
		t.Errorf("expected 1, got %d", stack.Size())
	}
}

func TestSizeIsNotChangedByPeek(t *testing.T) {
	stack := makeStack(10, 20)

	stack.Peek()

	if stack.Size() != 2 {
		t.Errorf("expected 2, got %d", stack.Size())
	}
}

// 1. Push

func TestPushOneItem(t *testing.T) {
	stack := makeStack(10)

	assertStack(t, stack, []int{10})
}

func TestPushKeepsOrder(t *testing.T) {
	stack := makeStack(10, 20, 30)

	assertStack(t, stack, []int{10, 20, 30})
}

func TestPushAfterEmptying(t *testing.T) {
	stack := makeStack(10)

	stack.Pop()
	stack.Push(20)

	assertStack(t, stack, []int{20})
}

func TestPushManyItems(t *testing.T) {
	stack := &Stack[int]{}
	expected := make([]int, 0, 100)

	for i := 0; i < 100; i++ {
		stack.Push(i)
		expected = append(expected, i)
	}

	assertStack(t, stack, expected)
}

// 1. Pop

func TestPopReturnsLastPushed(t *testing.T) {
	stack := makeStack(10, 20)

	got, err := stack.Pop()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 20 {
		t.Errorf("expected 20, got %d", got)
	}
}

func TestPopEmptyStack(t *testing.T) {
	stack := makeStack()

	if _, err := stack.Pop(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestPopAfterEmptying(t *testing.T) {
	stack := makeStack(10)

	stack.Pop()

	if _, err := stack.Pop(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestPopEveryItem(t *testing.T) {
	stack := makeStack(10, 20, 30)

	assertStack(t, stack, []int{10, 20, 30})
}

// 1. Peek

func TestPeekReturnsTop(t *testing.T) {
	stack := makeStack(10, 20)

	got, err := stack.Peek()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 20 {
		t.Errorf("expected 20, got %d", got)
	}
}

func TestPeekEmptyStack(t *testing.T) {
	stack := makeStack()

	if _, err := stack.Peek(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestPeekTwiceReturnsSameItem(t *testing.T) {
	stack := makeStack(10, 20)

	first, _ := stack.Peek()
	second, _ := stack.Peek()

	if first != second {
		t.Errorf("expected the same item, got %d and %d", first, second)
	}
}

func TestPeekAfterPop(t *testing.T) {
	stack := makeStack(10, 20)

	stack.Pop()

	got, err := stack.Peek()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 10 {
		t.Errorf("expected 10, got %d", got)
	}
}

// Стек строк: проверка шаблона типа

func TestStackOfStrings(t *testing.T) {
	stack := &Stack[string]{}

	stack.Push("a")
	stack.Push("b")

	got, err := stack.Pop()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != "b" {
		t.Errorf("expected %q, got %q", "b", got)
	}
}

// 2. Стек на односвязном списке

func TestListStackEmpty(t *testing.T) {
	stack := &ListStack[int]{}

	if stack.Size() != 0 {
		t.Errorf("expected 0, got %d", stack.Size())
	}

	if _, err := stack.Pop(); err == nil {
		t.Error("expected error on Pop, got nil")
	}

	if _, err := stack.Peek(); err == nil {
		t.Error("expected error on Peek, got nil")
	}
}

func TestListStackKeepsOrder(t *testing.T) {
	stack := &ListStack[int]{}

	for _, value := range []int{10, 20, 30} {
		stack.Push(value)
	}

	if stack.Size() != 3 {
		t.Fatalf("expected size 3, got %d", stack.Size())
	}

	for _, want := range []int{30, 20, 10} {
		got, err := stack.Pop()

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got != want {
			t.Errorf("expected %d, got %d", want, got)
		}
	}

	if stack.Size() != 0 {
		t.Errorf("expected empty stack, got size %d", stack.Size())
	}
}

func TestListStackPeek(t *testing.T) {
	stack := &ListStack[int]{}

	stack.Push(10)
	stack.Push(20)

	got, err := stack.Peek()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 20 {
		t.Errorf("expected 20, got %d", got)
	}

	if stack.Size() != 2 {
		t.Errorf("Peek must not change size, got %d", stack.Size())
	}
}

func TestListStackReusableAfterEmptying(t *testing.T) {
	stack := &ListStack[int]{}

	stack.Push(10)
	stack.Pop()
	stack.Push(20)

	got, err := stack.Peek()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 20 {
		t.Errorf("expected 20, got %d", got)
	}
}

// 3. Цикл с двумя pop подряд

func TestDoublePopOnEvenSize(t *testing.T) {
	stack := makeStack(10, 20, 30, 40)

	for stack.Size() > 0 {
		stack.Pop()

		if _, err := stack.Pop(); err != nil {
			t.Errorf("unexpected error on even size: %v", err)
		}
	}

	if stack.Size() != 0 {
		t.Errorf("expected empty stack, got size %d", stack.Size())
	}
}

func TestDoublePopOnOddSize(t *testing.T) {
	stack := makeStack(10, 20, 30)
	failed := false

	for stack.Size() > 0 {
		stack.Pop()

		if _, err := stack.Pop(); err != nil {
			failed = true
		}
	}

	if !failed {
		t.Error("expected the second Pop to fail on odd size")
	}
}

// 4*. Баланс круглых скобок

func TestIsBalanced(t *testing.T) {
	cases := []struct {
		expression string
		want       bool
	}{
		{"", true},
		{"()", true},
		{"(())", true},
		{"(()((())()))", true},
		{"()()", true},
		{"(", false},
		{")", false},
		{"(()()(()", false},
		{"())(", false},
		{"))((", false},
		{"((())", false},
	}

	for _, c := range cases {
		if got := IsBalanced(c.expression); got != c.want {
			t.Errorf("%q: expected %v, got %v", c.expression, c.want, got)
		}
	}
}

// 5*. Баланс скобок трёх типов

func TestIsBalancedAll(t *testing.T) {
	cases := []struct {
		expression string
		want       bool
	}{
		{"", true},
		{"()", true},
		{"[]{}", true},
		{"{[()]}", true},
		{"{{[[(())]]}}", true},
		{"([{}])", true},
		{"(]", false},
		{"{[}]", false},
		{"([)]", false},
		{"{", false},
		{"}", false},
		{"{[()]}}", false},
	}

	for _, c := range cases {
		if got := IsBalancedAll(c.expression); got != c.want {
			t.Errorf("%q: expected %v, got %v", c.expression, c.want, got)
		}
	}
}

// 6*. Минимум за O(1)

func TestMinStackEmpty(t *testing.T) {
	stack := &MinStack[int]{}

	if _, err := stack.Min(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestMinStackTracksMinimum(t *testing.T) {
	stack := &MinStack[int]{}

	for _, value := range []int{5, 3, 7, 1, 4} {
		stack.Push(value)
	}

	if got, _ := stack.Min(); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}

	// Снимаем 4 и 1: минимумом снова становится 3.
	stack.Pop()
	stack.Pop()

	if got, _ := stack.Min(); got != 3 {
		t.Errorf("expected 3, got %d", got)
	}
}

func TestMinStackWithDuplicates(t *testing.T) {
	stack := &MinStack[int]{}

	for _, value := range []int{2, 1, 1, 3} {
		stack.Push(value)
	}

	stack.Pop()
	stack.Pop()

	if got, _ := stack.Min(); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}

	stack.Pop()

	if got, _ := stack.Min(); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestMinStackPopReturnsItems(t *testing.T) {
	stack := &MinStack[int]{}

	stack.Push(10)
	stack.Push(20)

	if got, _ := stack.Pop(); got != 20 {
		t.Errorf("expected 20, got %d", got)
	}

	if stack.Size() != 1 {
		t.Errorf("expected size 1, got %d", stack.Size())
	}

	if got, _ := stack.Peek(); got != 10 {
		t.Errorf("expected 10, got %d", got)
	}

	stack.Pop()

	if _, err := stack.Pop(); err == nil {
		t.Error("expected error on empty stack, got nil")
	}
}

// 7*. Среднее за O(1)

func TestAvgStackEmpty(t *testing.T) {
	stack := &AvgStack[int]{}

	if _, err := stack.Average(); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestAvgStackAverage(t *testing.T) {
	stack := &AvgStack[int]{}

	for _, value := range []int{2, 4, 6} {
		stack.Push(value)
	}

	got, err := stack.Average()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 4 {
		t.Errorf("expected 4, got %v", got)
	}
}

func TestAvgStackAfterPop(t *testing.T) {
	stack := &AvgStack[int]{}

	for _, value := range []int{1, 2, 9} {
		stack.Push(value)
	}

	stack.Pop()

	got, _ := stack.Average()

	if got != 1.5 {
		t.Errorf("expected 1.5, got %v", got)
	}

	if stack.Size() != 2 {
		t.Errorf("expected size 2, got %d", stack.Size())
	}
}

func TestAvgStackNegativeValues(t *testing.T) {
	stack := &AvgStack[int]{}

	for _, value := range []int{-4, 0, 4} {
		stack.Push(value)
	}

	got, _ := stack.Average()

	if got != 0 {
		t.Errorf("expected 0, got %v", got)
	}

	if _, err := stack.Peek(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestAvgStackPopOnEmpty(t *testing.T) {
	stack := &AvgStack[int]{}

	if _, err := stack.Pop(); err == nil {
		t.Error("expected error, got nil")
	}
}

// 8*. Постфиксные выражения

func TestCalcPostfixFromTask(t *testing.T) {
	got, err := CalcPostfix("8 2 + 5 * 9 + =")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got != 59 {
		t.Errorf("expected 59, got %d", got)
	}
}

func TestCalcPostfixExpressions(t *testing.T) {
	cases := []struct {
		expression string
		want       int
	}{
		{"1 2 + 3 * =", 9},
		{"1 2 +", 3},
		{"5 =", 5},
		{"2 3 * 4 * =", 24},
		{"10 4 - =", 6},
		{"7 -3 + =", 4},
	}

	for _, c := range cases {
		got, err := CalcPostfix(c.expression)

		if err != nil {
			t.Fatalf("%q: unexpected error: %v", c.expression, err)
		}

		if got != c.want {
			t.Errorf("%q: expected %d, got %d", c.expression, c.want, got)
		}
	}
}

func TestCalcPostfixBadExpressions(t *testing.T) {
	cases := []string{"", "1 +", "+ =", "1 2", "1 2 ^ =", "1 2 + 3"}

	for _, expression := range cases {
		if _, err := CalcPostfix(expression); err == nil {
			t.Errorf("%q: expected error, got nil", expression)
		}
	}
}
