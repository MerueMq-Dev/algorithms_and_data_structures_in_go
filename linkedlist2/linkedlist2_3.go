package main

import "testing"

// Задание 2. Двусвязный (двунаправленный) список.

func makeList(values ...int) LinkedList2 {
	var list LinkedList2
	for _, value := range values {
		list.AddInTail(Node{value: value})
	}
	return list
}

func listValues(list *LinkedList2) []int {
	var values []int
	for node := list.head; node != nil; node = node.next {
		values = append(values, node.value)
	}
	return values
}

func assertList(t *testing.T, list *LinkedList2, expected []int) {
	t.Helper()
	actual := listValues(list)

	if len(actual) != len(expected) {
		t.Fatalf("expected list %v, got %v", expected, actual)
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("index %d: expected %d, got %d", i, expected[i], actual[i])
		}
	}

	var reversed []int
	for node := list.tail; node != nil; node = node.prev {
		reversed = append(reversed, node.value)
	}

	if len(reversed) != len(expected) {
		t.Fatalf("backward traversal gave %v, expected %d nodes", reversed, len(expected))
	}

	for i := range expected {
		if want := expected[len(expected)-1-i]; reversed[i] != want {
			t.Errorf("backward index %d: expected %d, got %d", i, want, reversed[i])
		}
	}

	if got := list.Count(); got != len(expected) {
		t.Errorf("Count() returned %d, but list holds %d nodes", got, len(expected))
	}

	if len(expected) == 0 {
		if list.head != nil {
			t.Error("expected head to be nil")
		}
		if list.tail != nil {
			t.Error("expected tail to be nil")
		}
		return
	}

	if list.head == nil || list.tail == nil {
		t.Fatal("expected head and tail to be non-nil")
	}

	if list.head.prev != nil {
		t.Error("expected head.prev to be nil")
	}

	if list.tail.next != nil {
		t.Error("expected tail.next to be nil")
	}

	for node := list.head; node.next != nil; node = node.next {
		if node.next.prev != node {
			t.Errorf("broken back link between %d and %d", node.value, node.next.value)
		}
	}
}

// AddInTail

func TestAddInTailEmptyList(t *testing.T) {
	var list LinkedList2
	list.AddInTail(Node{value: 10})
	assertList(t, &list, []int{10})

	if list.head != list.tail {
		t.Error("head and tail should point to the same node")
	}
}

func TestAddInTailMultipleNodes(t *testing.T) {
	list := makeList(10, 20, 30)
	assertList(t, &list, []int{10, 20, 30})
}

// 1. Find

func TestFindExistingNode(t *testing.T) {
	list := makeList(10, 20, 30)
	node, err := list.Find(20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if node.value != 20 {
		t.Errorf("expected value 20, got %d", node.value)
	}
}

func TestFindFirstMatchingNode(t *testing.T) {
	list := makeList(10, 20, 20, 30)
	node, err := list.Find(20)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if node.prev == nil || node.prev.value != 10 {
		t.Error("expected the first of the two matching nodes")
	}
}

func TestFindMissingNode(t *testing.T) {
	list := makeList(10, 20, 30)

	if _, err := list.Find(40); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestFindEmptyList(t *testing.T) {
	var list LinkedList2

	if _, err := list.Find(10); err == nil {
		t.Error("expected error, got nil")
	}
}

// 2. FindAll

func TestFindAllEmptyList(t *testing.T) {
	var list LinkedList2

	if nodes := list.FindAll(10); len(nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(nodes))
	}
}

func TestFindAllNoMatches(t *testing.T) {
	list := makeList(10, 20, 30)

	if nodes := list.FindAll(40); len(nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(nodes))
	}
}

func TestFindAllOneMatch(t *testing.T) {
	list := makeList(10, 20, 30)
	nodes := list.FindAll(20)

	if len(nodes) != 1 {
		t.Fatalf("expected 1 node, got %d", len(nodes))
	}

	if nodes[0].value != 20 {
		t.Errorf("expected value 20, got %d", nodes[0].value)
	}
}

func TestFindAllMultipleMatches(t *testing.T) {
	list := makeList(10, 20, 20, 30, 20)
	nodes := list.FindAll(20)

	if len(nodes) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(nodes))
	}

	for i, node := range nodes {
		if node.value != 20 {
			t.Errorf("node %d: expected value 20, got %d", i, node.value)
		}
	}
}

// 3. Delete одного узла

func TestDeleteFromEmptyList(t *testing.T) {
	var list LinkedList2
	list.Delete(10, false)
	assertList(t, &list, []int{})
}

func TestDeleteFirstNode(t *testing.T) {
	list := makeList(10, 20, 30)
	list.Delete(10, false)
	assertList(t, &list, []int{20, 30})
}

func TestDeleteMiddleNode(t *testing.T) {
	list := makeList(10, 20, 30)
	list.Delete(20, false)
	assertList(t, &list, []int{10, 30})
}

func TestDeleteLastNode(t *testing.T) {
	list := makeList(10, 20, 30)
	list.Delete(30, false)
	assertList(t, &list, []int{10, 20})
}

func TestDeleteOnlyNode(t *testing.T) {
	list := makeList(10)
	list.Delete(10, false)
	assertList(t, &list, []int{})
}

func TestDeleteMissingNode(t *testing.T) {
	list := makeList(10, 20, 30)
	list.Delete(40, false)
	assertList(t, &list, []int{10, 20, 30})
}

func TestDeleteFirstOccurrenceOnly(t *testing.T) {
	list := makeList(10, 20, 20, 30)
	list.Delete(20, false)
	assertList(t, &list, []int{10, 20, 30})
}

// 4. Delete всех узлов

func TestDeleteAllOccurrences(t *testing.T) {
	list := makeList(10, 20, 20, 30, 20)
	list.Delete(20, true)
	assertList(t, &list, []int{10, 30})
}

func TestDeleteAllOccurrencesFromBeginning(t *testing.T) {
	list := makeList(20, 20, 20, 30)
	list.Delete(20, true)
	assertList(t, &list, []int{30})
}

func TestDeleteAllOccurrencesFromEnd(t *testing.T) {
	list := makeList(10, 20, 20, 20)
	list.Delete(20, true)
	assertList(t, &list, []int{10})
}

func TestDeleteAllOccurrencesOfOnlyValue(t *testing.T) {
	list := makeList(20, 20, 20)
	list.Delete(20, true)
	assertList(t, &list, []int{})
}

func TestDeleteAllThenAddInTail(t *testing.T) {
	list := makeList(20, 20)
	list.Delete(20, true)
	list.AddInTail(Node{value: 10})
	assertList(t, &list, []int{10})
}

// 5. Insert

func TestInsertIntoEmptyListWithNil(t *testing.T) {
	var list LinkedList2
	list.Insert(nil, Node{value: 10})
	assertList(t, &list, []int{10})
}

func TestInsertAtBeginningWithNil(t *testing.T) {
	list := makeList(20, 30)
	list.Insert(nil, Node{value: 10})
	assertList(t, &list, []int{10, 20, 30})
}

func TestInsertAfterHead(t *testing.T) {
	list := makeList(10, 30)
	list.Insert(list.head, Node{value: 20})
	assertList(t, &list, []int{10, 20, 30})
}

func TestInsertAfterMiddle(t *testing.T) {
	list := makeList(10, 20, 40)
	list.Insert(list.head.next, Node{value: 30})
	assertList(t, &list, []int{10, 20, 30, 40})
}

func TestInsertAfterTail(t *testing.T) {
	list := makeList(10, 20)
	list.Insert(list.tail, Node{value: 30})
	assertList(t, &list, []int{10, 20, 30})
}

// 6. InsertFirst

func TestInsertFirstEmptyList(t *testing.T) {
	var list LinkedList2
	list.InsertFirst(Node{value: 10})
	assertList(t, &list, []int{10})

	if list.head != list.tail {
		t.Error("head and tail should point to the same node")
	}
}

func TestInsertFirstNonEmptyList(t *testing.T) {
	list := makeList(20, 30)
	list.InsertFirst(Node{value: 10})
	assertList(t, &list, []int{10, 20, 30})
}

func TestInsertFirstDoesNotChangeTail(t *testing.T) {
	list := makeList(20, 30)
	tail := list.tail

	list.InsertFirst(Node{value: 10})

	if list.tail != tail {
		t.Error("tail should not change after InsertFirst")
	}
}

// 7. Clean

func TestCleanEmptyList(t *testing.T) {
	var list LinkedList2
	list.Clean()
	assertList(t, &list, []int{})
}

func TestCleanNonEmptyList(t *testing.T) {
	list := makeList(10, 20, 30)
	list.Clean()
	assertList(t, &list, []int{})
}

func TestCleanThenAdd(t *testing.T) {
	list := makeList(10, 20, 30)
	list.Clean()
	list.AddInTail(Node{value: 40})
	assertList(t, &list, []int{40})
}

// Count

func TestCountEmpty(t *testing.T) {
	var list LinkedList2

	if got := list.Count(); got != 0 {
		t.Errorf("expected 0, got %d", got)
	}
}

func TestCountOneElement(t *testing.T) {
	list := makeList(10)

	if got := list.Count(); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestCountMultipleElements(t *testing.T) {
	list := makeList(10, 20, 30, 40, 50)

	if got := list.Count(); got != 5 {
		t.Errorf("expected 5, got %d", got)
	}
}

// 9*. Reverse

func TestReverseEmptyList(t *testing.T) {
	var list LinkedList2
	list.Reverse()
	assertList(t, &list, []int{})
}

func TestReverseSingleElement(t *testing.T) {
	list := makeList(10)
	list.Reverse()
	assertList(t, &list, []int{10})
}

func TestReverseManyElements(t *testing.T) {
	list := makeList(10, 20, 30, 40)
	list.Reverse()
	assertList(t, &list, []int{40, 30, 20, 10})
}

func TestReverseTwiceRestoresOrder(t *testing.T) {
	list := makeList(10, 20, 30)
	list.Reverse()
	list.Reverse()
	assertList(t, &list, []int{10, 20, 30})
}

// 10*. HasCycle

func TestHasCycleEmptyList(t *testing.T) {
	var list LinkedList2

	if list.HasCycle() {
		t.Error("empty list has no cycle")
	}
}

func TestHasCycleNormalList(t *testing.T) {
	list := makeList(10, 20, 30)

	if list.HasCycle() {
		t.Error("expected no cycle in a normal list")
	}
}

func TestHasCycleTailPointsToHead(t *testing.T) {
	list := makeList(10, 20, 30)
	list.tail.next = list.head

	if !list.HasCycle() {
		t.Error("expected cycle when tail points to head")
	}
}

func TestHasCycleSelfLoop(t *testing.T) {
	list := makeList(10)
	list.head.next = list.head

	if !list.HasCycle() {
		t.Error("expected cycle when node points to itself")
	}
}

// 11*. Sort

func TestSortEmptyAndSingle(t *testing.T) {
	var empty LinkedList2
	empty.Sort()
	assertList(t, &empty, []int{})

	single := makeList(10)
	single.Sort()
	assertList(t, &single, []int{10})
}

func TestSortReversedList(t *testing.T) {
	list := makeList(4, 3, 2, 1)
	list.Sort()
	assertList(t, &list, []int{1, 2, 3, 4})
}

func TestSortWithDuplicates(t *testing.T) {
	list := makeList(3, 1, 3, 1, 2)
	list.Sort()
	assertList(t, &list, []int{1, 1, 2, 3, 3})
}

func TestSortNegativeValues(t *testing.T) {
	list := makeList(0, -5, 10, -1)
	list.Sort()
	assertList(t, &list, []int{-5, -1, 0, 10})
}

// 12*. MergeSorted

func TestMergeSortedBothEmpty(t *testing.T) {
	var list1, list2 LinkedList2
	result := MergeSorted(list1, list2)
	assertList(t, &result, []int{})
}

func TestMergeSortedOneEmpty(t *testing.T) {
	var empty LinkedList2
	list := makeList(1, 2, 3)

	result := MergeSorted(empty, list)

	assertList(t, &result, []int{1, 2, 3})
}

func TestMergeSortedDifferentLengths(t *testing.T) {
	list1 := makeList(1, 3)
	list2 := makeList(2, 4, 5, 6)

	result := MergeSorted(list1, list2)

	assertList(t, &result, []int{1, 2, 3, 4, 5, 6})
}

func TestMergeSortedDoesNotModifyInputs(t *testing.T) {
	list1 := makeList(1, 3)
	list2 := makeList(2, 4)

	MergeSorted(list1, list2)

	assertList(t, &list1, []int{1, 3})
	assertList(t, &list2, []int{2, 4})
}

// 13*. Список с фиктивными узлами

func TestDummyListAddAndInsertFirst(t *testing.T) {
	list := NewDummyLinkedList()

	list.AddInTail(Node{value: 20})
	list.AddInTail(Node{value: 30})
	list.InsertFirst(Node{value: 10})

	assertInts(t, list.Values(), []int{10, 20, 30})
}

func TestDummyListInsertBeforeMiddle(t *testing.T) {
	list := NewDummyLinkedList()

	list.AddInTail(Node{value: 10})
	middle := list.AddInTail(Node{value: 30})
	list.InsertBefore(middle, Node{value: 20})

	assertInts(t, list.Values(), []int{10, 20, 30})
}

func TestDummyListRemoveNode(t *testing.T) {
	list := NewDummyLinkedList()

	list.AddInTail(Node{value: 10})
	middle := list.AddInTail(Node{value: 20})
	list.AddInTail(Node{value: 30})

	list.Remove(middle)

	assertInts(t, list.Values(), []int{10, 30})
}

func TestDummyListDeleteAllOccurrences(t *testing.T) {
	list := NewDummyLinkedList()

	for _, value := range []int{7, 10, 7} {
		list.AddInTail(Node{value: value})
	}

	list.Delete(7, true)

	assertInts(t, list.Values(), []int{10})

	if list.Count() != 1 {
		t.Errorf("expected 1 node, got %d", list.Count())
	}
}

func assertInts(t *testing.T, actual []int, expected []int) {
	t.Helper()

	if len(actual) != len(expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}

	for i := range expected {
		if actual[i] != expected[i] {
			t.Errorf("index %d: expected %d, got %d", i, expected[i], actual[i])
		}
	}
}
