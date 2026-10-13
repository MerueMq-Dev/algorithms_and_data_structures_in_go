package main

import (
	"constraints"
	"os"
	"reflect"
	"strconv"
)

// Задание 10. Множество.
const setSize = 20000

type PowerSet[T constraints.Ordered] struct {
	buckets [][]T
	count   int
}

func Init[T constraints.Ordered]() PowerSet[T] {
	return PowerSet[T]{buckets: make([][]T, setSize)}
}

// Время: O(1)
// Память: O(1)
func (s *PowerSet[T]) Size() int {
	return s.count
}

// 1. Уже имеющееся значение не добавляется.
// Время: O(1) в среднем
// Память: O(1) амортизированно
func (s *PowerSet[T]) Put(value T) {
	if s.Get(value) {
		return
	}

	if s.buckets == nil {
		s.buckets = make([][]T, setSize)
	}

	i := bucketIndex(value)
	s.buckets[i] = append(s.buckets[i], value)
	s.count++
}

// Время: O(1) в среднем
// Память: O(1)
func (s *PowerSet[T]) Get(value T) bool {
	_, _, err := s.find(value)
	return err == nil
}

// 2. Время: O(1) в среднем
// Память: O(1)
func (s *PowerSet[T]) Remove(value T) bool {
	i, j, err := s.find(value)
	if err != nil {
		return false
	}

	bucket := s.buckets[i]
	last := len(bucket) - 1
	bucket[j] = bucket[last]

	var empty T
	bucket[last] = empty
	s.buckets[i] = bucket[:last]
	s.count--

	return true
}

// 2. Время: O(min(n, m))
// Память: O(min(n, m))
func (s *PowerSet[T]) Intersection(set2 PowerSet[T]) PowerSet[T] {
	small, big := s, &set2
	if small.count > big.count {
		small, big = big, small
	}

	result := Init[T]()

	for _, value := range small.values() {
		if big.Get(value) {
			result.Put(value)
		}
	}

	return result
}

// 2. Время: O(n + m)
// Память: O(n + m)
func (s *PowerSet[T]) Union(set2 PowerSet[T]) PowerSet[T] {
	result := Init[T]()

	for _, value := range s.values() {
		result.Put(value)
	}

	for _, value := range set2.values() {
		result.Put(value)
	}

	return result
}

// 2. Время: O(n)
// Память: O(n)
func (s *PowerSet[T]) Difference(set2 PowerSet[T]) PowerSet[T] {
	result := Init[T]()

	for _, value := range s.values() {
		if !set2.Get(value) {
			result.Put(value)
		}
	}

	return result
}

// 2. Время: O(m)
// Память: O(m)
func (s *PowerSet[T]) IsSubset(set2 PowerSet[T]) bool {
	if set2.count > s.count {
		return false
	}

	for _, value := range set2.values() {
		if !s.Get(value) {
			return false
		}
	}

	return true
}

// 2. Время: O(n)
// Память: O(n)
func (s *PowerSet[T]) Equals(set2 PowerSet[T]) bool {
	return s.count == set2.count && s.IsSubset(set2)
}

func (s *PowerSet[T]) find(value T) (int, int, error) {
	if s.buckets == nil {
		return -1, -1, os.ErrNotExist
	}

	i := bucketIndex(value)

	for j, item := range s.buckets[i] {
		if item == value {
			return i, j, nil
		}
	}

	return -1, -1, os.ErrNotExist
}

func (s *PowerSet[T]) values() []T {
	result := make([]T, 0, s.count)

	for _, bucket := range s.buckets {
		result = append(result, bucket...)
	}

	return result
}

// Значение переводится в байты и хэшируется по FNV-1a.
// Время: O(k), где k — длина значения
// Память: O(k)
func bucketIndex[T constraints.Ordered](value T) int {
	var buf [32]byte

	hash := uint32(2166136261)

	for _, b := range appendKey(buf[:0], value) {
		hash ^= uint32(b)
		hash *= 16777619
	}

	return int(hash % setSize)
}

func appendKey[T constraints.Ordered](buf []byte, value T) []byte {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String:
		return append(buf, v.String()...)
	case reflect.Float32, reflect.Float64:
		f := v.Float()
		if f == 0 {
			f = 0
		}

		return strconv.AppendFloat(buf, f, 'g', -1, 64)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.AppendUint(buf, v.Uint(), 10)
	default:
		return strconv.AppendInt(buf, v.Int(), 10)
	}
}
