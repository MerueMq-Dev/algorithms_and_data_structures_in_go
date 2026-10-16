package main

import (
	"slices"
	"strconv"
	"testing"
)

// Задание 11. Фильтр Блюма.

func testStrings() []string {
	digits := "0123456789"
	var result []string

	for i := 0; i < 10; i++ {
		result = append(result, digits[i:]+digits[:i])
	}

	return result
}

func makeFilter(values ...string) *BloomFilter {
	bf := &BloomFilter{filter_len: 32}

	for _, value := range values {
		bf.Add(value)
	}

	return bf
}

// 1. Биты

func TestSetAndClearBit(t *testing.T) {
	bf := makeFilter()

	bf.setBit(0)
	bf.setBit(5)
	bf.setBit(31)

	if bf.bits != 1<<0|1<<5|1<<31 {
		t.Errorf("unexpected bits %032b", bf.bits)
	}

	bf.clearBit(5)

	if bf.getBit(5) || !bf.getBit(0) || !bf.getBit(31) {
		t.Errorf("unexpected bits %032b", bf.bits)
	}
}

func TestNewBloomFilter(t *testing.T) {
	for _, length := range []int{0, 33} {
		if _, err := NewBloomFilter(length); err == nil {
			t.Errorf("expected error for length %d", length)
		}
	}

	if bf, err := NewBloomFilter(32); err != nil || bf.filter_len != 32 {
		t.Errorf("unexpected result %v, %v", bf, err)
	}
}

// 1. Хэш-функции

func TestHashSmallStrings(t *testing.T) {
	bf := makeFilter()

	cases := []struct {
		s            string
		hash1, hash2 int
	}{
		{"", 0, 0},
		{"a", 1, 1},
		{"ab", 19, 1},
	}

	for _, c := range cases {
		if got := bf.Hash1(c.s); got != c.hash1 {
			t.Errorf("Hash1(%q): expected %d, got %d", c.s, c.hash1, got)
		}

		if got := bf.Hash2(c.s); got != c.hash2 {
			t.Errorf("Hash2(%q): expected %d, got %d", c.s, c.hash2, got)
		}
	}
}

func TestHashInRange(t *testing.T) {
	bf := makeFilter()

	for _, s := range append(testStrings(), "hello", "привет") {
		if h := bf.Hash1(s); h < 0 || h >= 32 {
			t.Errorf("Hash1(%q) = %d out of range", s, h)
		}

		if h := bf.Hash2(s); h < 0 || h >= 32 {
			t.Errorf("Hash2(%q) = %d out of range", s, h)
		}
	}
}

// 1. Add и IsValue

func TestEmptyFilter(t *testing.T) {
	bf := makeFilter()

	for _, s := range testStrings() {
		if bf.IsValue(s) {
			t.Errorf("expected %q to be missing", s)
		}
	}
}

func TestValueCopies(t *testing.T) {
	var bf BloomFilter
	bf.filter_len = 32

	copied := bf
	copied.Add("0123456789")

	if !copied.IsValue("0123456789") || bf.IsValue("0123456789") {
		t.Error("expected copies to be independent")
	}
}

func TestAddSetsHashBits(t *testing.T) {
	bf := makeFilter("0123456789")

	if bf.bits != 1<<13|1<<5 {
		t.Errorf("expected bits 13 and 5, got %032b", bf.bits)
	}
}

func TestAllTestStringsFound(t *testing.T) {
	bf := makeFilter(testStrings()...)

	for _, s := range testStrings() {
		if !bf.IsValue(s) {
			t.Errorf("expected %q to be found", s)
		}
	}
}

func TestMissingStrings(t *testing.T) {
	bf := makeFilter(testStrings()...)

	for _, s := range []string{"", "a", "abc", "hello", "0000000000"} {
		if bf.IsValue(s) {
			t.Errorf("expected %q to be missing", s)
		}
	}
}

func TestFalsePositiveExists(t *testing.T) {
	bf := makeFilter(testStrings()...)

	if !bf.IsValue("9876543210") {
		t.Error("expected a false positive")
	}
}

// 2*. Слияние

func TestMergeEqualsSingleFilter(t *testing.T) {
	values := testStrings()

	first := makeFilter(values[:3]...)
	second := makeFilter(values[3:7]...)
	third := makeFilter(values[7:]...)

	merged, err := Merge(first, second, third)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if merged.bits != makeFilter(values...).bits {
		t.Errorf("expected merged bits to match one filter")
	}

	for _, s := range values {
		if !merged.IsValue(s) {
			t.Errorf("expected %q in merged filter", s)
		}
	}
}

func TestMergeKeepsSources(t *testing.T) {
	first := makeFilter("a")
	second := makeFilter("abc")

	Merge(first, second)

	if first.IsValue("abc") || second.IsValue("a") {
		t.Error("expected source filters to stay unchanged")
	}
}

func TestMergeEmptyFilters(t *testing.T) {
	merged, err := Merge(makeFilter(), makeFilter())

	if err != nil || merged.IsValue("a") {
		t.Errorf("expected empty merged filter, got %v", err)
	}
}

func TestMergeErrors(t *testing.T) {
	if _, err := Merge(); err == nil {
		t.Error("expected error for no filters")
	}

	if _, err := Merge(makeFilter(), &BloomFilter{filter_len: 16}); err == nil {
		t.Error("expected error for different sizes")
	}
}

// 3*. Фильтр с удалением

func TestCountingAddRemove(t *testing.T) {
	cf, _ := NewCountingBloomFilter(32)

	cf.Add("a")
	cf.Add("abc")

	if err := cf.Remove("a"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cf.IsValue("a") || !cf.IsValue("abc") {
		t.Error("expected only abc to remain")
	}
}

func TestCountingSharedBits(t *testing.T) {
	cf, _ := NewCountingBloomFilter(32)

	cf.Add("0123456789")
	cf.Add("2345678901")

	cf.Remove("0123456789")

	if !cf.IsValue("2345678901") {
		t.Error("expected shared bits to survive one removal")
	}
}

func TestCountingDuplicates(t *testing.T) {
	cf, _ := NewCountingBloomFilter(32)

	cf.Add("a")
	cf.Add("a")
	cf.Remove("a")

	if !cf.IsValue("a") {
		t.Error("expected one copy of a to remain")
	}

	cf.Remove("a")

	if cf.IsValue("a") {
		t.Error("expected a to be removed")
	}
}

func TestCountingRemoveMissing(t *testing.T) {
	cf, _ := NewCountingBloomFilter(32)
	cf.Add("abc")

	if err := cf.Remove("hello"); err == nil {
		t.Error("expected error for missing value")
	}

	if !cf.IsValue("abc") {
		t.Error("expected abc to survive")
	}
}

func TestCountingFalsePositiveBreaksFilter(t *testing.T) {
	cf, _ := NewCountingBloomFilter(32)
	cf.Add("1234567890")

	if err := cf.Remove("9876543210"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if cf.IsValue("1234567890") {
		t.Error("expected the known limitation: added value is lost")
	}
}

func TestCountingSaturation(t *testing.T) {
	cf, _ := NewCountingBloomFilter(32)

	for i := 0; i < 300; i++ {
		cf.Add("a")
	}

	for i := 0; i < 300; i++ {
		cf.Remove("a")
	}

	if !cf.IsValue("a") {
		t.Error("expected saturated counter to stay set")
	}
}

// 4*. Восстановление

func TestRecoverSureAndPossible(t *testing.T) {
	bf := makeFilter("a", "abc")

	candidates := []string{"a", "abc", "hello", "ab", "zzz"}

	sure, possible, err := Recover(bf, candidates)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, s := range []string{"a", "abc"} {
		if !slices.Contains(possible, s) || !slices.Contains(sure, s) {
			t.Errorf("expected %q to be recovered for sure", s)
		}
	}

	if slices.Contains(possible, "hello") {
		t.Error("expected hello to be excluded")
	}
}

func TestRecoverAmbiguous(t *testing.T) {
	bf := makeFilter(testStrings()...)

	var candidates []string
	candidates = append(candidates, testStrings()...)

	for i := 0; i < 1000; i++ {
		candidates = append(candidates, strconv.Itoa(1000000000+i*7919))
	}

	sure, possible, err := Recover(bf, candidates)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for _, s := range testStrings() {
		if !slices.Contains(possible, s) {
			t.Errorf("expected %q among possible", s)
		}
	}

	if len(sure) != 0 {
		t.Errorf("expected no sure values, got %v", sure)
	}
}

func TestRecoverIncompleteCandidates(t *testing.T) {
	bf := makeFilter("a", "abc")

	if _, _, err := Recover(bf, []string{"a"}); err == nil {
		t.Error("expected error when candidates miss added values")
	}
}
