package main

import (
	"strconv"
	"testing"
)

// Задание 12. Кэш.

// Числа, кратные 5, в таблице на 5 слотов все попадают в слот 0.
func makeFullCache() *NativeCache[string] {
	nc := NewNativeCache[string](5)

	for _, key := range []string{"0", "5", "10", "15", "20"} {
		nc.Put(key, "v"+key)
	}

	return nc
}

func assertHits(t *testing.T, nc *NativeCache[string], key string, want int) {
	t.Helper()

	got, err := nc.Hits(key)

	if err != nil {
		t.Fatalf("Hits(%q): unexpected error: %v", key, err)
	}

	if got != want {
		t.Errorf("Hits(%q): expected %d, got %d", key, want, got)
	}
}

// Put и Get

func TestPutGet(t *testing.T) {
	nc := NewNativeCache[int](7)

	nc.Put("cat", 1)
	nc.Put("dog", 2)

	if got, err := nc.Get("cat"); err != nil || got != 1 {
		t.Errorf("expected 1, got %d, %v", got, err)
	}

	if _, err := nc.Get("fox"); err == nil {
		t.Error("expected error for missing key")
	}
}

func TestPutExistingKeyUpdatesValue(t *testing.T) {
	nc := makeFullCache()

	nc.Get("5")
	nc.Put("5", "new")

	if got, _ := nc.Get("5"); got != "new" {
		t.Errorf("expected new, got %q", got)
	}

	assertHits(t, nc, "5", 2)

	for _, key := range []string{"0", "10", "15", "20"} {
		if !nc.IsKey(key) {
			t.Errorf("expected %q to stay, no eviction on update", key)
		}
	}
}

// Счётчик обращений

func TestHitsCounted(t *testing.T) {
	nc := makeFullCache()

	for i := 0; i < 3; i++ {
		nc.Get("10")
	}

	nc.Get("20")

	assertHits(t, nc, "10", 3)
	assertHits(t, nc, "20", 1)
	assertHits(t, nc, "0", 0)
}

func TestHitsNotChangedByIsKeyAndMisses(t *testing.T) {
	nc := makeFullCache()

	nc.IsKey("0")
	nc.Get("missing")
	nc.Hits("0")

	assertHits(t, nc, "0", 0)

	if _, err := nc.Hits("missing"); err == nil {
		t.Error("expected error for missing key")
	}
}

// Вытеснение

func TestTableFullByCollisions(t *testing.T) {
	nc := makeFullCache()

	for i, key := range []string{"0", "5", "10", "15", "20"} {
		if nc.Find(key) != i {
			t.Errorf("expected %q in slot %d, got %d", key, i, nc.Find(key))
		}
	}

	if nc.SeekSlot("25") != -1 {
		t.Error("expected table to be full")
	}
}

func TestEvictsLeastHit(t *testing.T) {
	nc := makeFullCache()

	for _, key := range []string{"0", "5", "15", "20"} {
		nc.Get(key)
	}

	nc.Put("25", "v25")

	if nc.IsKey("10") {
		t.Error("expected 10 to be evicted")
	}

	if got, err := nc.Get("25"); err != nil || got != "v25" {
		t.Errorf("expected v25, got %q, %v", got, err)
	}

	if nc.Find("25") != 2 {
		t.Errorf("expected 25 in freed slot 2, got %d", nc.Find("25"))
	}

	for _, key := range []string{"0", "5", "15", "20"} {
		if got, err := nc.Get(key); err != nil || got != "v"+key {
			t.Errorf("expected %q to stay, got %q, %v", key, got, err)
		}
	}
}

func TestNewKeyStartsWithZeroHits(t *testing.T) {
	nc := makeFullCache()

	for _, key := range []string{"5", "10", "15", "20"} {
		nc.Get(key)
		nc.Get(key)
	}

	nc.Get("0")
	nc.Put("25", "v25")

	assertHits(t, nc, "25", 0)

	if nc.IsKey("0") {
		t.Error("expected 0 to be evicted")
	}
}

func TestEvictionTieTakesFirstSlot(t *testing.T) {
	nc := makeFullCache()

	nc.Put("25", "v25")

	if nc.IsKey("0") || nc.Find("25") != 0 {
		t.Error("expected key in slot 0 to be evicted on tie")
	}
}

func TestRepeatedEvictions(t *testing.T) {
	nc := makeFullCache()

	for i := 0; i < 5; i++ {
		nc.Get("15")
	}

	for i := 5; i < 30; i++ {
		key := strconv.Itoa(i * 5)
		nc.Put(key, "v"+key)
		nc.Get(key)
	}

	if !nc.IsKey("15") {
		t.Error("expected most used key to survive")
	}

	count := 0

	for _, filled := range nc.filled {
		if filled {
			count++
		}
	}

	if count != 5 {
		t.Errorf("expected 5 filled slots, got %d", count)
	}
}

func TestEvictionWithDifferentHashes(t *testing.T) {
	nc := NewNativeCache[int](3)

	nc.Put("a", 1)
	nc.Put("b", 2)
	nc.Put("c", 3)

	nc.Get("a")
	nc.Get("c")

	nc.Put("d", 4)

	if nc.IsKey("b") {
		t.Error("expected b to be evicted")
	}

	for key, want := range map[string]int{"a": 1, "c": 3, "d": 4} {
		if got, err := nc.Get(key); err != nil || got != want {
			t.Errorf("Get(%q): expected %d, got %d, %v", key, want, got, err)
		}
	}
}
