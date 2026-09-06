package main

import (
	"errors"
)

// 8.*  Сложение двух списков одинаковой длины.
// Время: O(n)
// Память: O(n)
func SumLists(list1 LinkedList, list2 LinkedList) (LinkedList, error) {
	if list1.Count() != list2.Count() {
		return LinkedList{}, errors.New("lists have different lengths")
	}

	var result LinkedList
	current1 := list1.head
	current2 := list2.head

	for current1 != nil {
		result.AddInTail(Node{
			value: current1.value + current2.value,
		})

		current1 = current1.next
		current2 = current2.next
	}

	return result, nil
}
