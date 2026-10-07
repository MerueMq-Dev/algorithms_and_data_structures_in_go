package main

import (
	"strconv"
	"testing"
)

// Задание 8. Хэш-таблица.

// Сколько проб от start с шагом step нужно, чтобы дойти до index.
func probesTo(size int, start int, step int, index int) int {
	for i := 0; i < size; i++ {
		if (start+i*step)%size == index {
			return i + 1
		}
	}

	return size
}

// HashFun

func TestHashFunInRange(t *testing.T) {
	ht := Init(17, 3)

	for _, value := range []string{"", "a", "hello", "очень длинная строка"} {
		if index := ht.HashFun(value); index < 0 || index >= 17 {
			t.Errorf("HashFun(%q) = %d out of range", value, index)
		}
	}
}

func TestHashFunDeterministic(t *testing.T) {
	ht := Init(17, 3)

	if ht.HashFun("cat") != ht.HashFun("cat") {
		t.Error("expected same hash for same value")
	}
}

func TestHashFunNumbers(t *testing.T) {
	ht := Init(17, 3)

	cases := map[string]int{"0": 0, "5": 5, "22": 5, "-1": 16, "+3": 3}

	for value, want := range cases {
		if got := ht.HashFun(value); got != want {
			t.Errorf("HashFun(%q): expected %d, got %d", value, want, got)
		}
	}
}

func TestHashFunSpreadsValues(t *testing.T) {
	ht := Init(17, 3)
	used := map[int]bool{}

	for i := 0; i < 100; i++ {
		used[ht.HashFun("key"+strconv.Itoa(i))] = true
	}

	if len(used) != 17 {
		t.Errorf("expected all 17 slots used, got %d", len(used))
	}
}

// SeekSlot

func TestSeekSlotEmptyTable(t *testing.T) {
	ht := Init(17, 3)

	if got := ht.SeekSlot("cat"); got != ht.HashFun("cat") {
		t.Errorf("expected %d, got %d", ht.HashFun("cat"), got)
	}
}

func TestSeekSlotSkipsByStep(t *testing.T) {
	ht := Init(17, 3)
	keys := CollidingKeys(&ht, 3)

	ht.Put(keys[0])
	ht.Put(keys[1])

	if got := ht.SeekSlot(keys[2]); got != 6 {
		t.Errorf("expected 6, got %d", got)
	}
}

func TestSeekSlotWrapsAround(t *testing.T) {
	ht := Init(5, 3)
	keys := CollidingKeys(&ht, 3)

	ht.Put(keys[0])
	ht.Put(keys[1])

	if got := ht.SeekSlot(keys[2]); got != 1 {
		t.Errorf("expected 1, got %d", got)
	}
}

func TestSeekSlotFullTable(t *testing.T) {
	ht := Init(5, 3)

	for _, key := range CollidingKeys(&ht, 5) {
		ht.Put(key)
	}

	if got := ht.SeekSlot("cat"); got != -1 {
		t.Errorf("expected -1, got %d", got)
	}
}

// Put

func TestPutStoresValue(t *testing.T) {
	ht := Init(17, 3)

	index := ht.Put("cat")

	if index != ht.HashFun("cat") {
		t.Errorf("expected %d, got %d", ht.HashFun("cat"), index)
	}

	if !ht.filled[index] || ht.slots[index] != "cat" {
		t.Errorf("expected cat in slot %d", index)
	}
}

func TestPutCollision(t *testing.T) {
	ht := Init(17, 3)
	keys := CollidingKeys(&ht, 2)

	first := ht.Put(keys[0])
	second := ht.Put(keys[1])

	if first != 0 || second != 3 {
		t.Errorf("expected slots 0 and 3, got %d and %d", first, second)
	}
}

func TestPutEmptyString(t *testing.T) {
	ht := Init(17, 3)

	index := ht.Put("")

	if index < 0 || ht.Find("") != index {
		t.Errorf("expected empty string to be stored, got %d", index)
	}
}

func TestPutFullTable(t *testing.T) {
	ht := Init(5, 3)

	for i := 0; i < 5; i++ {
		if ht.Put("key"+strconv.Itoa(i)) < 0 {
			t.Fatalf("unexpected failure on key%d", i)
		}
	}

	if got := ht.Put("extra"); got != -1 {
		t.Errorf("expected -1, got %d", got)
	}
}

// Find

func TestFindExisting(t *testing.T) {
	ht := Init(17, 3)

	index := ht.Put("cat")

	if got := ht.Find("cat"); got != index {
		t.Errorf("expected %d, got %d", index, got)
	}
}

func TestFindMissing(t *testing.T) {
	ht := Init(17, 3)
	ht.Put("cat")

	if got := ht.Find("dog"); got != -1 {
		t.Errorf("expected -1, got %d", got)
	}
}

func TestFindEmptyTable(t *testing.T) {
	ht := Init(17, 3)

	if got := ht.Find("cat"); got != -1 {
		t.Errorf("expected -1, got %d", got)
	}
}

func TestFindAfterCollisions(t *testing.T) {
	ht := Init(17, 3)
	keys := CollidingKeys(&ht, 4)

	for _, key := range keys {
		ht.Put(key)
	}

	for i, key := range keys {
		if got := ht.Find(key); got != i*3 {
			t.Errorf("Find(%q): expected %d, got %d", key, i*3, got)
		}
	}
}

func TestFindMissingInFullTable(t *testing.T) {
	ht := Init(5, 3)

	for i := 0; i < 5; i++ {
		ht.Put("key" + strconv.Itoa(i))
	}

	if got := ht.Find("cat"); got != -1 {
		t.Errorf("expected -1, got %d", got)
	}
}

// 3*. Динамическая таблица

func TestDynamicHashTableGrows(t *testing.T) {
	d := NewDynamicHashTable(5, 3)

	for i := 0; i < 100; i++ {
		if d.Put("key"+strconv.Itoa(i)) < 0 {
			t.Fatalf("unexpected failure on key%d", i)
		}
	}

	if d.Count() != 100 {
		t.Errorf("expected count 100, got %d", d.Count())
	}

	if d.Size() < 134 || !isPrime(d.Size()) {
		t.Errorf("expected prime size with load <= 0.75, got %d", d.Size())
	}

	for i := 0; i < 100; i++ {
		key := "key" + strconv.Itoa(i)

		if index := d.Find(key); index < 0 || d.table.slots[index] != key {
			t.Errorf("expected to find %s after growth", key)
		}
	}

	if d.Find("missing") != -1 {
		t.Error("expected -1 for missing value")
	}
}

func TestNextPrime(t *testing.T) {
	cases := map[int]int{0: 2, 2: 2, 10: 11, 34: 37, 97: 97}

	for n, want := range cases {
		if got := nextPrime(n); got != want {
			t.Errorf("nextPrime(%d): expected %d, got %d", n, want, got)
		}
	}
}

// 4*. Двойное хэширование

func TestDoubleHashTablePutFind(t *testing.T) {
	d := NewDoubleHashTable(17)

	for i := 0; i < 17; i++ {
		if d.Put("key"+strconv.Itoa(i)) < 0 {
			t.Fatalf("unexpected failure on key%d", i)
		}
	}

	if d.Put("extra") != -1 {
		t.Error("expected -1 on full table")
	}

	for i := 0; i < 17; i++ {
		if d.Find("key"+strconv.Itoa(i)) < 0 {
			t.Errorf("expected to find key%d", i)
		}
	}

	if d.Find("extra") != -1 {
		t.Error("expected -1 for missing value")
	}
}

func TestDoubleHashingNeedsFewerProbes(t *testing.T) {
	linear := Init(101, 3)
	double := NewDoubleHashTable(101)
	keys := CollidingKeys(&linear, 30)

	linearProbes := 0
	doubleProbes := 0

	for _, key := range keys {
		index := linear.Put(key)
		linearProbes += probesTo(101, linear.HashFun(key), 3, index)

		index = double.Put(key)
		doubleProbes += probesTo(101, double.table.HashFun(key), double.StepFun(key), index)
	}

	if doubleProbes*3 >= linearProbes {
		t.Errorf("expected far fewer probes: linear %d, double %d", linearProbes, doubleProbes)
	}
}

// 5*. Атака и соль

func TestCollidingKeysHitOneSlot(t *testing.T) {
	ht := Init(1009, 3)

	for _, key := range CollidingKeys(&ht, 50) {
		if ht.HashFun(key) != 0 {
			t.Fatalf("expected %q to hash to 0", key)
		}
	}
}

func TestAttackQuadraticProbes(t *testing.T) {
	ht := Init(1009, 3)
	keys := CollidingKeys(&ht, 300)

	probes := 0

	for _, key := range keys {
		index := ht.Put(key)
		probes += probesTo(1009, 0, 3, index)
	}

	if want := 300 * 301 / 2; probes != want {
		t.Errorf("expected %d probes, got %d", want, probes)
	}
}

func TestSaltedHashTableResistsAttack(t *testing.T) {
	plain := Init(1009, 3)
	salted := NewSaltedHashTable(1009, 3)
	keys := CollidingKeys(&plain, 300)

	probes := 0

	for _, key := range keys {
		index := salted.Put(key)

		if index < 0 {
			t.Fatalf("unexpected failure on %q", key)
		}

		probes += probesTo(1009, salted.HashFun(key), 3, index)
	}

	if probes > 300*3 {
		t.Errorf("expected about one probe per key, got %d for 300 keys", probes)
	}

	for _, key := range keys {
		if salted.Find(key) < 0 {
			t.Errorf("expected to find %q", key)
		}
	}
}

func TestSaltDiffersBetweenTables(t *testing.T) {
	first := NewSaltedHashTable(1009, 3)
	second := NewSaltedHashTable(1009, 3)

	same := 0

	for i := 0; i < 50; i++ {
		key := strconv.Itoa(i)

		if first.HashFun(key) == second.HashFun(key) {
			same++
		}
	}

	if same == 50 {
		t.Error("expected different salts to give different slots")
	}
}
