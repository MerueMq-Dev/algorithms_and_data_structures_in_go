package main

import (
	"math"
	"slices"
	"strconv"
	"testing"
	"time"
)

// Задание 10. Множество.

func makeSet(values ...string) PowerSet[string] {
	s := Init[string]()

	for _, value := range values {
		s.Put(value)
	}

	return s
}

func assertSet(t *testing.T, s PowerSet[string], expected ...string) {
	t.Helper()

	if s.Size() != len(expected) {
		t.Fatalf("expected size %d, got %d", len(expected), s.Size())
	}

	for _, value := range expected {
		if !s.Get(value) {
			t.Errorf("expected %q in set", value)
		}
	}
}

// 1. Put

func TestPutNewValue(t *testing.T) {
	s := makeSet()

	s.Put("cat")

	assertSet(t, s, "cat")
}

func TestPutExistingValue(t *testing.T) {
	s := makeSet("cat", "dog")

	s.Put("cat")

	assertSet(t, s, "cat", "dog")
}

func TestPutZeroValueSet(t *testing.T) {
	var s PowerSet[string]

	if s.Get("cat") || s.Remove("cat") {
		t.Error("expected empty zero-value set")
	}

	s.Put("cat")

	assertSet(t, s, "cat")
}

func TestPutOtherTypes(t *testing.T) {
	ints := Init[int]()
	ints.Put(-5)
	ints.Put(-5)
	ints.Put(7)

	if ints.Size() != 2 || !ints.Get(-5) || ints.Get(5) {
		t.Error("unexpected int set")
	}

	floats := Init[float64]()
	floats.Put(0.0)
	floats.Put(math.Copysign(0, -1))
	floats.Put(1.5)

	if floats.Size() != 2 || !floats.Get(1.5) {
		t.Errorf("unexpected float set, size %d", floats.Size())
	}
}

// Get

func TestGetMissing(t *testing.T) {
	s := makeSet("cat")

	if s.Get("dog") {
		t.Error("expected dog to be missing")
	}
}

// 2. Remove

func TestRemoveExisting(t *testing.T) {
	s := makeSet("cat", "dog", "fox")

	if !s.Remove("dog") {
		t.Error("expected true")
	}

	assertSet(t, s, "cat", "fox")

	if s.Get("dog") {
		t.Error("expected dog to be removed")
	}
}

func TestRemoveMissing(t *testing.T) {
	s := makeSet("cat")

	if s.Remove("dog") {
		t.Error("expected false")
	}

	assertSet(t, s, "cat")
}

func TestRemoveThenPut(t *testing.T) {
	s := makeSet("cat")

	s.Remove("cat")
	s.Put("cat")

	assertSet(t, s, "cat")
}

// 2. Intersection

func TestIntersectionNonEmpty(t *testing.T) {
	a := makeSet("a", "b", "c")
	b := makeSet("b", "c", "d")

	assertSet(t, a.Intersection(b), "b", "c")
}

func TestIntersectionEmpty(t *testing.T) {
	a := makeSet("a", "b")
	b := makeSet("c", "d")

	assertSet(t, a.Intersection(b))
	assertSet(t, a.Intersection(makeSet()))
}

func TestIntersectionKeepsOperands(t *testing.T) {
	a := makeSet("a", "b")
	b := makeSet("b")

	a.Intersection(b)

	assertSet(t, a, "a", "b")
	assertSet(t, b, "b")
}

// 2. Union

func TestUnionBothNonEmpty(t *testing.T) {
	a := makeSet("a", "b")
	b := makeSet("b", "c")

	assertSet(t, a.Union(b), "a", "b", "c")
}

func TestUnionWithEmpty(t *testing.T) {
	a := makeSet("a", "b")
	empty := makeSet()

	assertSet(t, a.Union(empty), "a", "b")
	assertSet(t, empty.Union(a), "a", "b")
}

// 2. Difference

func TestDifferenceNonEmpty(t *testing.T) {
	a := makeSet("a", "b", "c")
	b := makeSet("b", "d")

	assertSet(t, a.Difference(b), "a", "c")
}

func TestDifferenceEmpty(t *testing.T) {
	a := makeSet("a", "b")
	b := makeSet("a", "b", "c")

	empty := makeSet()

	assertSet(t, a.Difference(b))
	assertSet(t, empty.Difference(a))
}

// 2. IsSubset

func TestIsSubsetParamInside(t *testing.T) {
	a := makeSet("a", "b", "c")
	b := makeSet("a", "c")

	if !a.IsSubset(b) {
		t.Error("expected b to be a subset of a")
	}
}

func TestIsSubsetCurrentInside(t *testing.T) {
	a := makeSet("a", "c")
	b := makeSet("a", "b", "c")

	if a.IsSubset(b) {
		t.Error("expected b not to be a subset of a")
	}
}

func TestIsSubsetPartial(t *testing.T) {
	a := makeSet("a", "b", "c")
	b := makeSet("a", "d")

	if a.IsSubset(b) {
		t.Error("expected b not to be a subset of a")
	}
}

func TestIsSubsetEmptyAndSelf(t *testing.T) {
	a := makeSet("a", "b")

	if !a.IsSubset(makeSet()) || !a.IsSubset(a) {
		t.Error("expected empty set and itself to be subsets")
	}
}

// 2. Equals

func TestEquals(t *testing.T) {
	cases := []struct {
		name string
		a, b []string
		want bool
	}{
		{"same", []string{"a", "b"}, []string{"b", "a"}, true},
		{"both empty", nil, nil, true},
		{"different size", []string{"a"}, []string{"a", "b"}, false},
		{"same size different values", []string{"a", "b"}, []string{"a", "c"}, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			a := makeSet(c.a...)

			if got := a.Equals(makeSet(c.b...)); got != c.want {
				t.Errorf("expected %v, got %v", c.want, got)
			}
		})
	}
}

// Быстродействие

func TestPerformance(t *testing.T) {
	start := time.Now()

	a := Init[string]()
	b := Init[string]()

	for i := 0; i < 20000; i++ {
		a.Put("value" + strconv.Itoa(i))
		b.Put("value" + strconv.Itoa(i+10000))
	}

	assertSize := func(name string, s PowerSet[string], want int) {
		if s.Size() != want {
			t.Errorf("%s: expected size %d, got %d", name, want, s.Size())
		}
	}

	assertSize("intersection", a.Intersection(b), 10000)
	assertSize("union", a.Union(b), 30000)
	assertSize("difference", a.Difference(b), 10000)

	if a.IsSubset(b) || a.Equals(b) {
		t.Error("unexpected subset or equality")
	}

	for i := 0; i < 20000; i++ {
		a.Remove("value" + strconv.Itoa(i))
	}

	assertSize("after remove", a, 0)

	if elapsed := time.Since(start); elapsed > 2*time.Second {
		t.Errorf("expected under 2s, took %v", elapsed)
	}
}

// 4*. Декартово произведение

func TestCartesian(t *testing.T) {
	a := makeSet("a", "b")
	b := makeSet("1", "2", "3")

	pairs := a.Cartesian(b)

	if len(pairs) != 6 {
		t.Fatalf("expected 6 pairs, got %d", len(pairs))
	}

	for _, first := range []string{"a", "b"} {
		for _, second := range []string{"1", "2", "3"} {
			if !slices.Contains(pairs, Pair[string]{first, second}) {
				t.Errorf("expected pair (%s, %s)", first, second)
			}
		}
	}
}

func TestCartesianWithEmpty(t *testing.T) {
	a := makeSet("a")

	if len(a.Cartesian(makeSet())) != 0 {
		t.Error("expected no pairs")
	}
}

// 5*. Пересечение нескольких множеств

func TestIntersectAll(t *testing.T) {
	sets := []PowerSet[string]{
		makeSet("a", "b", "c", "d"),
		makeSet("b", "c", "d"),
		makeSet("c", "d", "e"),
		makeSet("d", "c"),
	}

	result, err := IntersectAll(sets)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assertSet(t, result, "c", "d")
}

func TestIntersectAllEmptyResult(t *testing.T) {
	sets := []PowerSet[string]{makeSet("a"), makeSet("a"), makeSet("b")}

	result, _ := IntersectAll(sets)

	assertSet(t, result)
}

func TestIntersectAllTooFew(t *testing.T) {
	if _, err := IntersectAll([]PowerSet[string]{makeSet("a"), makeSet("a")}); err == nil {
		t.Error("expected error for 2 sets")
	}
}

// 6*. Мультимножество

func TestBagAddAndCount(t *testing.T) {
	b := NewBag[string]()

	for _, value := range []string{"a", "b", "a", "c", "a"} {
		b.Add(value)
	}

	if b.Count("a") != 3 || b.Count("b") != 1 || b.Count("z") != 0 {
		t.Error("unexpected counts")
	}

	if b.Size() != 5 {
		t.Errorf("expected size 5, got %d", b.Size())
	}
}

func TestBagRemoveOne(t *testing.T) {
	b := NewBag[string]()
	b.Add("a")
	b.Add("a")

	if !b.Remove("a") || b.Count("a") != 1 {
		t.Error("expected one a to remain")
	}

	if !b.Remove("a") || b.Count("a") != 0 {
		t.Error("expected no a to remain")
	}

	if b.Remove("a") {
		t.Error("expected false when nothing to remove")
	}

	if b.Size() != 0 || len(b.Frequencies()) != 0 {
		t.Error("expected empty bag")
	}
}

func TestBagFrequencies(t *testing.T) {
	b := NewBag[string]()

	for _, value := range []string{"pear", "apple", "pear", "fig", "pear", "apple"} {
		b.Add(value)
	}

	want := []Frequency[string]{{"apple", 2}, {"fig", 1}, {"pear", 3}}

	if got := b.Frequencies(); !slices.Equal(got, want) {
		t.Errorf("expected %v, got %v", want, got)
	}
}
