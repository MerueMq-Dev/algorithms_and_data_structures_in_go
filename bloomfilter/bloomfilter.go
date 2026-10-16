package main

import (
	"os"
)

// Задание 11. Фильтр Блюма.
const (
	random1  = 17
	random2  = 223
	wordBits = 32
)

type BloomFilter struct {
	filter_len int
	bits       []uint32
}

func NewBloomFilter(filterLen int) (*BloomFilter, error) {
	if filterLen < 1 {
		return nil, os.ErrInvalid
	}

	return &BloomFilter{filter_len: filterLen}, nil
}

// Время: O(k), где k — длина строки
// Память: O(1)
func (bf *BloomFilter) Hash1(s string) int {
	return bf.hash(s, random1)
}

// Время: O(k), где k — длина строки
// Память: O(1)
func (bf *BloomFilter) Hash2(s string) int {
	return bf.hash(s, random2)
}

// Время: O(k), где k — длина строки
// Память: O(1)
func (bf *BloomFilter) Add(s string) {
	bf.setBit(bf.Hash1(s))
	bf.setBit(bf.Hash2(s))
}

// Время: O(k), где k — длина строки
// Память: O(1)
func (bf *BloomFilter) IsValue(s string) bool {
	return bf.getBit(bf.Hash1(s)) && bf.getBit(bf.Hash2(s))
}

func (bf *BloomFilter) hash(s string, random int) int {
	if bf.filter_len < 1 {
		return 0
	}

	sum := 0

	for _, char := range s {
		code := int(char)
		sum = (sum*random + code) % bf.filter_len
	}

	return sum
}

func (bf *BloomFilter) setBit(i int) {
	bf.ensureBits()
	bf.bits[i/wordBits] |= 1 << (i % wordBits)
}

func (bf *BloomFilter) clearBit(i int) {
	bf.ensureBits()
	bf.bits[i/wordBits] &^= 1 << (i % wordBits)
}

func (bf *BloomFilter) getBit(i int) bool {
	if bf.bits == nil {
		return false
	}

	return bf.bits[i/wordBits]&(1<<(i%wordBits)) != 0
}

func (bf *BloomFilter) ensureBits() {
	if bf.bits == nil {
		bf.bits = make([]uint32, (bf.filter_len+wordBits-1)/wordBits)
	}
}
