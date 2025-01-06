package set

// Set 泛型的集合容器实现
type Set[T comparable] struct {
	values map[T]bool
}

// NewSet 构造集合
func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		values: make(map[T]bool),
	}
}

// NewSetWithValues 通过集合中的值构造集合
func NewSetWithValues[T comparable](values []T) *Set[T] {
	s := &Set[T]{
		values: make(map[T]bool, len(values)),
	}

	for _, v := range values {
		s.Add(v)
	}

	return s
}

// Size 集合的大小
func (s *Set[T]) Size() int {
	return len(s.values)
}

// Add 添加值
func (s *Set[T]) Add(v T) {
	s.values[v] = true
}

// Del 删除值
func (s *Set[T]) Del(v T) {
	delete(s.values, v)
}

// Contains 返回集合中是否包含提供的值
func (s *Set[T]) Contains(v T) bool {
	_, exist := s.values[v]
	return exist
}

// ToSlice 将集合中的值转换为切片
func (s *Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s.values))
	for v := range s.values {
		slice = append(slice, v)
	}
	return slice
}
