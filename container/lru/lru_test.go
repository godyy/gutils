package lru

import "testing"

type testValue struct {
	key   int
	value string
}

func (v testValue) LRUKey() int {
	return v.key
}

func TestLRU_PutGetAndEvict(t *testing.T) {
	lru := New[int, testValue](2)

	if _, rem := lru.Put(testValue{key: 1, value: "one"}); rem {
		t.Fatal("expected first put not to evict any item")
	}
	if _, rem := lru.Put(testValue{key: 2, value: "two"}); rem {
		t.Fatal("expected second put not to evict any item")
	}

	// 访问 key=1，确保后续淘汰的是 key=2。
	if v, ok := lru.Get(1); !ok || v.value != "one" {
		t.Fatalf("expected to get key=1, got value=%+v, ok=%v", v, ok)
	}

	evicted, rem := lru.Put(testValue{key: 3, value: "three"})
	if !rem {
		t.Fatal("expected put to evict one item")
	}
	if evicted.key != 2 || evicted.value != "two" {
		t.Fatalf("expected evicted value to be key=2, got %+v", evicted)
	}
	if lru.Size() != 2 {
		t.Fatalf("expected size 2, got %d", lru.Size())
	}
	if _, ok := lru.Get(2); ok {
		t.Fatal("expected key=2 to be evicted")
	}
	if v, ok := lru.Get(1); !ok || v.value != "one" {
		t.Fatalf("expected key=1 to remain in cache, got value=%+v, ok=%v", v, ok)
	}
	if v, ok := lru.Get(3); !ok || v.value != "three" {
		t.Fatalf("expected key=3 to be present, got value=%+v, ok=%v", v, ok)
	}
}

func TestLRU_UpdateMovesItemToFront(t *testing.T) {
	lru := New[int, testValue](2)

	lru.Put(testValue{key: 1, value: "one"})
	lru.Put(testValue{key: 2, value: "two"})

	// 更新已有项应同时刷新其最近使用顺序。
	if _, rem := lru.Put(testValue{key: 1, value: "one-new"}); rem {
		t.Fatal("expected updating existing key not to evict any item")
	}

	evicted, rem := lru.Put(testValue{key: 3, value: "three"})
	if !rem {
		t.Fatal("expected key=2 to be evicted after inserting key=3")
	}
	if evicted.key != 2 {
		t.Fatalf("expected key=2 to be evicted, got %+v", evicted)
	}
	if v, ok := lru.Get(1); !ok || v.value != "one-new" {
		t.Fatalf("expected updated key=1 to remain, got value=%+v, ok=%v", v, ok)
	}
}

func TestLRU_Clear(t *testing.T) {
	lru := New[int, testValue](2)

	lru.Put(testValue{key: 1, value: "one"})
	lru.Put(testValue{key: 2, value: "two"})
	lru.Clear()

	if lru.Size() != 0 {
		t.Fatalf("expected size 0 after clear, got %d", lru.Size())
	}
	if _, ok := lru.Get(1); ok {
		t.Fatal("expected key=1 to be absent after clear")
	}
	if _, ok := lru.Get(2); ok {
		t.Fatal("expected key=2 to be absent after clear")
	}

	// 清空后应可继续正常使用。
	if _, rem := lru.Put(testValue{key: 3, value: "three"}); rem {
		t.Fatal("expected put after clear not to evict any item")
	}
	if v, ok := lru.Get(3); !ok || v.value != "three" {
		t.Fatalf("expected key=3 to be present after clear, got value=%+v, ok=%v", v, ok)
	}
}

func TestLRU_ZeroCapacity(t *testing.T) {
	lru := New[int, testValue](0)

	evicted, rem := lru.Put(testValue{key: 1, value: "one"})
	if !rem {
		t.Fatal("expected zero-capacity cache to evict inserted item immediately")
	}
	if evicted.key != 1 || evicted.value != "one" {
		t.Fatalf("expected inserted item to be evicted immediately, got %+v", evicted)
	}
	if lru.Size() != 0 {
		t.Fatalf("expected size 0, got %d", lru.Size())
	}
	if _, ok := lru.Get(1); ok {
		t.Fatal("expected key=1 not to be cached when capacity is zero")
	}
}
