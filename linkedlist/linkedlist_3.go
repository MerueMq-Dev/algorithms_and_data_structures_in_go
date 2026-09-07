package main

import "testing"

func makeList(values ...int) LinkedList {
	var list LinkedList
	for _, value := range values {
		list.AddInTail(Node{value: value})
	}
	return list
}

func listValues(list *LinkedList) []int {
	var values []int
	current := list.head

	for current != nil {
		values = append(values, current.value)
		current = current.next
	}

	return values
}

func assertList(t *testing.T, list *LinkedList, expected []int) {
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

	if len(expected) == 0 {
		if list.head != nil {
			t.Error("expected head to be nil")
		}
		if list.tail != nil {
			t.Error("expected tail to be nil")
		}
		return
	}

	if list.head == nil {
		t.Fatal("expected head to be non-nil")
	}

	if list.tail == nil {
		t.Fatal("expected tail to be non-nil")
	}

	if list.tail.value != expected[len(expected)-1] {
		t.Errorf("expected tail value %d, got %d", expected[len(expected)-1], list.tail.value)
	}

	if list.tail.next != nil {
		t.Error("expected tail.next to be nil")
	}
}

// AddInTail

func TestAddInTailEmptyList(t *testing.T) {
	var list LinkedList

	list.AddInTail(Node{value: 10})

	assertList(t, &list, []int{10})

	if list.head != list.tail {
		t.Error("head and tail should point to the same node")
	}
}

func TestAddInTailMultipleNodes(t *testing.T) {
	var list LinkedList

	list.AddInTail(Node{value: 10})
	list.AddInTail(Node{value: 20})
	list.AddInTail(Node{value: 30})

	assertList(t, &list, []int{10, 20, 30})
}

func TestAddInTailNodeNextIsReset(t *testing.T) {
	var list LinkedList
	existing := &Node{value: 100}

	list.AddInTail(Node{value: 10, next: existing})

	assertList(t, &list, []int{10})

	if list.tail.next != nil {
		t.Error("tail.next should be nil")
	}
}

// Count

func TestCountEmpty(t *testing.T) {
	var list LinkedList

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

// Find

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

	if node.value != 20 {
		t.Errorf("expected value 20, got %d", node.value)
	}
}

func TestFindMissingNode(t *testing.T) {
	list := makeList(10, 20, 30)
	_, err := list.Find(40)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

func TestFindEmptyList(t *testing.T) {
	var list LinkedList
	_, err := list.Find(10)

	if err == nil {
		t.Error("expected error, got nil")
	}
}

// FindAll

func TestFindAllEmptyList(t *testing.T) {
	var list LinkedList
	nodes := list.FindAll(10)

	if len(nodes) != 0 {
		t.Errorf("expected 0 nodes, got %d", len(nodes))
	}
}

func TestFindAllNoMatches(t *testing.T) {
	list := makeList(10, 20, 30)
	nodes := list.FindAll(40)

	if len(nodes) != 0 {
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

// Delete

func TestDeleteFirstFromEmptyList(t *testing.T) {
	var list LinkedList
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

// Insert

func TestInsertIntoEmptyList(t *testing.T) {
	var list LinkedList
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

func TestInsertAfterNodeNotInList(t *testing.T) {
	list := makeList(10, 20)
	foreign := &Node{value: 100}
	list.Insert(foreign, Node{value: 30})
	assertList(t, &list, []int{10, 20})
}

// InsertFirst

func TestInsertFirstEmptyList(t *testing.T) {
	var list LinkedList
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

func TestInsertFirstMultipleTimes(t *testing.T) {
	list := makeList(30)
	list.InsertFirst(Node{value: 20})
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

// Clean

func TestCleanEmptyList(t *testing.T) {
	var list LinkedList
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

	if list.head != list.tail {
		t.Error("head and tail should point to the same node")
	}
}
