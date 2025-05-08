package heap

import (
	"testing"
)

type testElement struct {
	value int
	index int
}

func (t *testElement) HeapLess(element *testElement) bool {
	return t.value < element.value
}

func (t *testElement) SetHeapIndex(i int) {
	t.index = i
}

func (t *testElement) HeapIndex() int {
	return t.index
}

func TestHeap(t *testing.T) {
	heap := NewHeap[*testElement](10)

	for i := 0; i < 1e2; i++ {
		heap.Push(&testElement{
			value: i,
			index: -1,
		})
	}
	t.Log(heap.Len(), cap(heap.list), heap.list)

	for heap.Len() > 0 {
		heap.Pop()
		t.Log(heap.Len(), cap(heap.list), heap.list)
	}
}

func BenchmarkHeap(b *testing.B) {
	heap := NewHeap[*testElement](10)

	b.Run("push", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			heap.Push(&testElement{
				value: i,
				index: -1,
			})
		}
	})

	b.Run("pop", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			heap.Pop()
		}
	})
}
