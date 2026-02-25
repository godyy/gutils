package skiplist

import (
	"math/rand"
	"sort"
	"testing"
)

func TestSkipList_Basic(t *testing.T) {
	sl := New[int, string]()

	// Test Insert
	sl.Set(10, "10")
	sl.Set(5, "5")
	sl.Set(20, "20")

	if sl.Len() != 3 {
		t.Errorf("Expected length 3, got %d", sl.Len())
	}

	// Test Get
	if val, ok := sl.Get(10); !ok || val != "10" {
		t.Errorf("Expected 10, got %v", val)
	}
	if val, ok := sl.Get(5); !ok || val != "5" {
		t.Errorf("Expected 5, got %v", val)
	}
	if val, ok := sl.Get(20); !ok || val != "20" {
		t.Errorf("Expected 20, got %v", val)
	}
	if _, ok := sl.Get(15); ok {
		t.Error("Expected key 15 not to exist")
	}

	// Test Update
	sl.Set(10, "10-new")
	if val, ok := sl.Get(10); !ok || val != "10-new" {
		t.Errorf("Expected 10-new, got %v", val)
	}

	// Test Remove
	if !sl.Remove(5) {
		t.Error("Expected to remove key 5")
	}
	if sl.Len() != 2 {
		t.Errorf("Expected length 2, got %d", sl.Len())
	}
	if _, ok := sl.Get(5); ok {
		t.Error("Expected key 5 to be removed")
	}
}

func TestSkipList_Ascend(t *testing.T) {
	sl := New[int, int]()
	expected := []int{}
	// Insert random values
	for i := 0; i < 100; i++ {
		val := rand.Intn(1000)
		sl.Set(val, val)
		expected = append(expected, val)
	}

	// Remove duplicates from expected and sort to match SkipList behavior
	sort.Ints(expected)
	unique := []int{}
	if len(expected) > 0 {
		unique = append(unique, expected[0])
		for i := 1; i < len(expected); i++ {
			if expected[i] != expected[i-1] {
				unique = append(unique, expected[i])
			}
		}
	}
	expected = unique

	var result []int
	sl.Ascend(func(key int, value int) bool {
		result = append(result, key)
		return true
	})

	if len(result) != len(expected) {
		t.Errorf("Expected length %d, got %d", len(expected), len(result))
	}

	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("Expected %d at index %d, got %d", expected[i], i, result[i])
		}
	}
}

func TestSkipList_Rank(t *testing.T) {
	sl := New[int, int]()
	// Insert 10, 20, 30, 40, 50
	elements := []int{10, 20, 30, 40, 50}
	for _, v := range elements {
		sl.Set(v, v)
	}

	// Check Ranks
	for i, v := range elements {
		rank := sl.GetRank(v)
		expectedRank := i + 1
		if rank != expectedRank {
			t.Errorf("Expected rank %d for key %d, got %d", expectedRank, v, rank)
		}
	}

	// Check non-existent key
	if rank := sl.GetRank(25); rank != 0 {
		t.Errorf("Expected rank 0 for non-existent key 25, got %d", rank)
	}

	// Check GetByRank
	for i, v := range elements {
		key, _, ok := sl.GetByRank(i + 1)
		if !ok {
			t.Errorf("Expected GetByRank(%d) to return true", i+1)
		}
		if key != v {
			t.Errorf("Expected key %d at rank %d, got %d", v, i+1, key)
		}
	}

	// Check out of bounds
	if _, _, ok := sl.GetByRank(0); ok {
		t.Error("Expected GetByRank(0) to return false")
	}
	if _, _, ok := sl.GetByRank(6); ok {
		t.Error("Expected GetByRank(6) to return false")
	}

	// Test Rank after Remove
	sl.Remove(30) // Remove middle element
	// Remaining: 10, 20, 40, 50
	expected := []int{10, 20, 40, 50}
	for i, v := range expected {
		rank := sl.GetRank(v)
		expectedRank := i + 1
		if rank != expectedRank {
			t.Errorf("Expected rank %d for key %d after removal, got %d", expectedRank, v, rank)
		}
	}
}

func BenchmarkSkipList_Set(b *testing.B) {
	sl := New[int, int]()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Set(i, i)
	}
}

func BenchmarkSkipList_Get(b *testing.B) {
	sl := New[int, int]()
	// Pre-fill
	limit := 100000
	for i := 0; i < limit; i++ {
		sl.Set(i, i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.Get(rand.Intn(limit))
	}
}

func BenchmarkSkipList_GetRank(b *testing.B) {
	sl := New[int, int]()
	limit := 100000
	for i := 0; i < limit; i++ {
		sl.Set(i, i)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sl.GetRank(rand.Intn(limit))
	}
}
