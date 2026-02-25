package skiplist

import (
	"math/rand"
	"time"
)

// Ordered represents types that can be compared with <, <=, >, >=.
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 |
		~string
}

const (
	defaultMaxLevel = 32
	defaultP        = 0.25
)

type node[K Ordered, V any] struct {
	key   K
	value V
	next  []*node[K, V]
	span  []int // span[i] is the distance to the next node at level i
}

// SkipList is a probabilistic data structure that allows O(log n) search complexity
// as well as O(log n) insertion complexity within an ordered sequence of elements.
// It also supports efficient rank calculation (finding the position of an element).
type SkipList[K Ordered, V any] struct {
	head     *node[K, V]
	maxLevel int
	level    int
	p        float64
	rand     *rand.Rand
	length   int
}

// New creates a new SkipList with default parameters.
func New[K Ordered, V any]() *SkipList[K, V] {
	return &SkipList[K, V]{
		head: &node[K, V]{
			next: make([]*node[K, V], defaultMaxLevel),
			span: make([]int, defaultMaxLevel),
		},
		maxLevel: defaultMaxLevel,
		level:    0,
		p:        defaultP,
		rand:     rand.New(rand.NewSource(time.Now().UnixNano())),
		length:   0,
	}
}

// randomLevel generates a random level for a new node.
func (s *SkipList[K, V]) randomLevel() int {
	level := 1
	for s.rand.Float64() < s.p && level < s.maxLevel {
		level++
	}
	return level
}

// Set inserts or updates a key-value pair in the SkipList.
func (s *SkipList[K, V]) Set(key K, value V) {
	update := make([]*node[K, V], s.maxLevel)
	rank := make([]int, s.maxLevel)
	current := s.head

	// Traverse the list to find the position
	for i := s.level - 1; i >= 0; i-- {
		// Store rank that is crossed to reach the insert position
		if i == s.level-1 {
			rank[i] = 0
		} else {
			rank[i] = rank[i+1]
		}
		for current.next[i] != nil && current.next[i].key < key {
			rank[i] += current.span[i]
			current = current.next[i]
		}
		update[i] = current
	}

	// If the key already exists, update the value
	if current.next[0] != nil && current.next[0].key == key {
		current.next[0].value = value
		return
	}

	// Determine the level for the new node
	newLevel := s.randomLevel()

	// If the new level is higher than the current level, initialize update array for new levels
	if newLevel > s.level {
		for i := s.level; i < newLevel; i++ {
			rank[i] = 0
			update[i] = s.head
			update[i].span[i] = s.length
		}
		s.level = newLevel
	}

	// Create the new node
	newNode := &node[K, V]{
		key:   key,
		value: value,
		next:  make([]*node[K, V], newLevel),
		span:  make([]int, newLevel),
	}

	// Insert the node by updating pointers and spans
	for i := 0; i < newLevel; i++ {
		newNode.next[i] = update[i].next[i]
		update[i].next[i] = newNode

		newNode.span[i] = update[i].span[i] - (rank[0] - rank[i])
		update[i].span[i] = (rank[0] - rank[i]) + 1
	}

	// Increment span for untouched levels
	for i := newLevel; i < s.level; i++ {
		update[i].span[i]++
	}

	s.length++
}

// Get retrieves the value associated with the given key.
// Returns the value and true if found, otherwise zero value and false.
func (s *SkipList[K, V]) Get(key K) (V, bool) {
	current := s.head
	for i := s.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
	}

	current = current.next[0]
	if current != nil && current.key == key {
		return current.value, true
	}

	var zero V
	return zero, false
}

// Remove deletes a key-value pair from the SkipList.
// Returns true if the key was found and removed, false otherwise.
func (s *SkipList[K, V]) Remove(key K) bool {
	update := make([]*node[K, V], s.maxLevel)
	current := s.head

	for i := s.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			current = current.next[i]
		}
		update[i] = current
	}

	current = current.next[0]

	if current != nil && current.key == key {
		for i := 0; i < s.level; i++ {
			if update[i].next[i] == current {
				update[i].span[i] += current.span[i] - 1
				update[i].next[i] = current.next[i]
			} else {
				update[i].span[i]--
			}
		}

		// Decrease the level of the skip list if necessary
		for s.level > 0 && s.head.next[s.level-1] == nil {
			s.level--
		}

		s.length--
		return true
	}

	return false
}

// GetRank returns the rank (1-based index) of the key.
// Returns 0 if the key is not found.
func (s *SkipList[K, V]) GetRank(key K) int {
	rank := 0
	current := s.head

	for i := s.level - 1; i >= 0; i-- {
		for current.next[i] != nil && current.next[i].key < key {
			rank += current.span[i]
			current = current.next[i]
		}
	}

	// Check if the next node is the key
	if current.next[0] != nil && current.next[0].key == key {
		return rank + 1
	}

	return 0
}

// GetByRank returns the key and value at the specified rank (1-based index).
// Returns zero values and false if the rank is out of bounds.
func (s *SkipList[K, V]) GetByRank(rank int) (K, V, bool) {
	var zeroK K
	var zeroV V

	if rank < 1 || rank > s.length {
		return zeroK, zeroV, false
	}

	current := s.head
	traversed := 0

	for i := s.level - 1; i >= 0; i-- {
		for current.next[i] != nil && (traversed+current.span[i]) <= rank {
			traversed += current.span[i]
			current = current.next[i]
		}
	}

	if traversed == rank {
		return current.key, current.value, true
	}

	return zeroK, zeroV, false
}

// Len returns the number of elements in the SkipList.
func (s *SkipList[K, V]) Len() int {
	return s.length
}

// Ascend iterates over the SkipList in ascending order.
// It calls the iterator function for each element.
// If the iterator returns false, the iteration stops.
func (s *SkipList[K, V]) Ascend(iterator func(key K, value V) bool) {
	current := s.head.next[0]
	for current != nil {
		if !iterator(current.key, current.value) {
			break
		}
		current = current.next[0]
	}
}
