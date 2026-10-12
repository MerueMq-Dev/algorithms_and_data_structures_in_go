package main

import (
	"strconv"
	"testing"
)

// Задание 9. Ассоциативный массив.

// HashFun

func TestHashFunInRange(t *testing.T) {
	nd := Init[int](17)

	for _, key := range []string{"", "a", "key", "-5", "123456789"} {
		if index := nd.HashFun(key); index < 0 || index >= 17 {
			t.Errorf("HashFun(%q) = %d out of range", key, index)
		}
	}
}

// 4. Put

func TestPutNewKey(t *testing.T) {
	nd := Init[int](17)

	nd.Put("cat", 5)

	if !nd.IsKey("cat") {
		t.Fatal("expected cat to be a key")
	}

	if got, err := nd.Get("cat"); err != nil || got != 5 {
		t.Errorf("expected 5, got %d, %v", got, err)
	}

	if nd.count != 1 {
		t.Errorf("expected count 1, got %d", nd.count)
	}
}

func TestPutUsesHashTableSlots(t *testing.T) {
	nd := Init[int](17)

	nd.Put("cat", 5)

	index := nd.HashTable.Find("cat")

	if index < 0 || nd.slots[index] != "cat" || nd.values[index] != 5 {
		t.Errorf("expected key in table slot and value at the same index, got %d", index)
	}
}

func TestPutExistingKey(t *testing.T) {
	nd := Init[int](17)

	nd.Put("cat", 5)
	nd.Put("cat", 7)

	if got, err := nd.Get("cat"); err != nil || got != 7 {
		t.Errorf("expected 7, got %d, %v", got, err)
	}

	if nd.count != 1 {
		t.Errorf("expected count 1, got %d", nd.count)
	}
}

func TestPutCollidingKeys(t *testing.T) {
	nd := Init[string](17)

	// Числа 3, 20 и 37 попадают в один слот.
	nd.Put("3", "three")
	nd.Put("20", "twenty")
	nd.Put("37", "thirty seven")

	cases := map[string]string{"3": "three", "20": "twenty", "37": "thirty seven"}

	for key, want := range cases {
		if got, err := nd.Get(key); err != nil || got != want {
			t.Errorf("Get(%q): expected %q, got %q, %v", key, want, got, err)
		}
	}
}

func TestPutGrowsWhenFull(t *testing.T) {
	nd := Init[int](3)

	for i := 0; i < 100; i++ {
		nd.Put("key"+strconv.Itoa(i), i)
	}

	if nd.count != 100 {
		t.Errorf("expected count 100, got %d", nd.count)
	}

	if float64(nd.count) > float64(nd.size)*maxLoad {
		t.Errorf("expected load <= %.2f, got %d of %d", maxLoad, nd.count, nd.size)
	}

	for i := 0; i < 100; i++ {
		if got, err := nd.Get("key" + strconv.Itoa(i)); err != nil || got != i {
			t.Errorf("key%d: expected %d, got %d, %v", i, i, got, err)
		}
	}
}

func TestPutZeroSize(t *testing.T) {
	nd := Init[int](0)

	nd.Put("a", 1)
	nd.Put("b", 2)

	if got, _ := nd.Get("b"); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

func TestPutStructValues(t *testing.T) {
	type cat struct {
		name string
		age  int
	}

	nd := Init[cat](17)
	nd.Put("tom", cat{"Tom", 3})

	if got, _ := nd.Get("tom"); got.name != "Tom" || got.age != 3 {
		t.Errorf("unexpected value %+v", got)
	}
}

// 4. IsKey

func TestIsKeyPresentAndMissing(t *testing.T) {
	nd := Init[int](17)
	nd.Put("cat", 1)

	if !nd.IsKey("cat") {
		t.Error("expected cat to be a key")
	}

	if nd.IsKey("dog") {
		t.Error("expected dog to be missing")
	}
}

func TestIsKeyEmptyDictionary(t *testing.T) {
	nd := Init[int](17)

	if nd.IsKey("cat") {
		t.Error("expected empty dictionary to have no keys")
	}
}

func TestIsKeyEmptyString(t *testing.T) {
	nd := Init[int](17)

	if nd.IsKey("") {
		t.Error("expected empty key to be missing")
	}

	nd.Put("", 1)

	if !nd.IsKey("") {
		t.Error("expected empty key to be present")
	}
}

// 4. Get

func TestGetExisting(t *testing.T) {
	nd := Init[float64](17)
	nd.Put("pi", 3.14)

	if got, err := nd.Get("pi"); err != nil || got != 3.14 {
		t.Errorf("expected 3.14, got %v, %v", got, err)
	}
}

func TestGetMissing(t *testing.T) {
	nd := Init[int](17)
	nd.Put("cat", 1)

	got, err := nd.Get("dog")

	if err == nil {
		t.Error("expected error, got nil")
	}

	if got != 0 {
		t.Errorf("expected zero value, got %d", got)
	}
}

func TestGetZeroValueIsNotMissing(t *testing.T) {
	nd := Init[int](17)
	nd.Put("zero", 0)

	if _, err := nd.Get("zero"); err != nil {
		t.Errorf("expected stored zero, got error %v", err)
	}
}

// 5*. Словарь на упорядоченном списке

func TestOrderedDictionaryPutGet(t *testing.T) {
	od := &OrderedDictionary[int]{}

	for i, key := range []string{"pear", "apple", "fig", "kiwi"} {
		od.Put(key, i)
	}

	od.Put("fig", 10)

	want := []string{"apple", "fig", "kiwi", "pear"}

	for i, e := range od.entries {
		if e.key != want[i] {
			t.Errorf("index %d: expected %q, got %q", i, want[i], e.key)
		}
	}

	if got, err := od.Get("fig"); err != nil || got != 10 {
		t.Errorf("expected 10, got %d, %v", got, err)
	}

	if od.Count() != 4 {
		t.Errorf("expected count 4, got %d", od.Count())
	}
}

func TestOrderedDictionaryMissing(t *testing.T) {
	od := &OrderedDictionary[int]{}
	od.Put("b", 1)

	for _, key := range []string{"a", "c", ""} {
		if od.IsKey(key) {
			t.Errorf("expected %q to be missing", key)
		}

		if _, err := od.Get(key); err == nil {
			t.Errorf("expected error for %q", key)
		}
	}
}

func TestOrderedDictionaryDelete(t *testing.T) {
	od := &OrderedDictionary[int]{}

	for i, key := range []string{"a", "b", "c"} {
		od.Put(key, i)
	}

	if err := od.Delete("b"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if od.IsKey("b") || od.Count() != 2 {
		t.Error("expected b to be deleted")
	}

	if err := od.Delete("b"); err == nil {
		t.Error("expected error on second delete")
	}

	if got, _ := od.Get("c"); got != 2 {
		t.Errorf("expected 2, got %d", got)
	}
}

// 6*. Словарь с битовыми ключами

func TestNewBitDictionaryWidth(t *testing.T) {
	for _, width := range []int{0, -1, 65} {
		if _, err := NewBitDictionary[int](width); err == nil {
			t.Errorf("expected error for width %d", width)
		}
	}

	if _, err := NewBitDictionary[int](64); err != nil {
		t.Errorf("unexpected error for width 64: %v", err)
	}
}

func TestBitDictionaryPutGet(t *testing.T) {
	bd, _ := NewBitDictionary[string](4)

	bd.Put("0000", "zero")
	bd.Put("1010", "ten")
	bd.Put("1011", "eleven")
	bd.Put("1010", "TEN")

	cases := map[string]string{"0000": "zero", "1010": "TEN", "1011": "eleven"}

	for key, want := range cases {
		if got, err := bd.Get(key); err != nil || got != want {
			t.Errorf("Get(%q): expected %q, got %q, %v", key, want, got, err)
		}
	}

	if bd.Count() != 3 {
		t.Errorf("expected count 3, got %d", bd.Count())
	}

	if bd.IsKey("1111") || bd.IsKey("1000") {
		t.Error("expected missing keys")
	}
}

func TestBitDictionaryBadKeys(t *testing.T) {
	bd, _ := NewBitDictionary[int](4)

	for _, key := range []string{"", "101", "10101", "10a1", "-101", "+101"} {
		if err := bd.Put(key, 1); err == nil {
			t.Errorf("expected error for %q", key)
		}

		if bd.IsKey(key) {
			t.Errorf("expected %q to be missing", key)
		}
	}

	if bd.Count() != 0 {
		t.Errorf("expected count 0, got %d", bd.Count())
	}
}

func TestBitDictionaryDelete(t *testing.T) {
	bd, _ := NewBitDictionary[int](4)

	bd.Put("1010", 10)
	bd.Put("1011", 11)

	if err := bd.Delete("1010"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if bd.IsKey("1010") || !bd.IsKey("1011") || bd.Count() != 1 {
		t.Error("expected only 1011 to remain")
	}

	if err := bd.Delete("1010"); err == nil {
		t.Error("expected error on second delete")
	}

	if err := bd.Delete("0000"); err == nil {
		t.Error("expected error for missing key")
	}

	bd.Delete("1011")

	if bd.root.children[0] != nil || bd.root.children[1] != nil {
		t.Error("expected empty branches to be removed")
	}
}

func TestBitDictionaryWide(t *testing.T) {
	bd, _ := NewBitDictionary[int](64)

	ones := ""
	for i := 0; i < 64; i++ {
		ones += "1"
	}

	bd.Put(ones, 1)

	if got, err := bd.Get(ones); err != nil || got != 1 {
		t.Errorf("expected 1, got %d, %v", got, err)
	}
}
