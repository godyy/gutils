package heap

import "container/heap"

// Element 实现Element即可作为Heap元素
type Element interface {
	// Less Element比较
	Less(Element) bool

	// SetIndex 设置索引
	SetIndex(int)

	// Index 获取索引
	Index() int
}

// heapList 堆元素列表
// 实现heap.Interface的容器
type heapList []Element

func (h heapList) Len() int {
	return len(h)
}

func (h heapList) Less(i, j int) bool {
	return h[i].Less(h[j])
}

func (h heapList) Swap(i, j int) {
	h[i].SetIndex(j)
	h[j].SetIndex(i)
	h[i], h[j] = h[j], h[i]
}

func (h *heapList) Push(x any) {
	e := x.(Element)
	e.SetIndex(h.Len())
	*h = append(*h, x.(Element))
}

func (h *heapList) Pop() any {
	n := h.Len() - 1
	e := (*h)[n]
	e.SetIndex(-1)
	*h = (*h)[:n]
	return e
}

// Heap 堆容器
type Heap struct {
	minCap int      // 堆的最小容量
	list   heapList // 堆元素列表
}

func NewHeap(minCap ...int) *Heap {
	h := &Heap{}
	if len(minCap) > 0 && minCap[0] > 0 {
		h.minCap = minCap[0]
	}
	h.list = make(heapList, 0, h.minCap)
	return h
}

// Init 使用所给元素初始化堆
func (h *Heap) Init(e ...Element) {
	c := len(e)
	if c < h.minCap {
		c = h.minCap
	}
	h.list = make(heapList, 0, c)
	for _, v := range e {
		h.list.Push(v)
	}
	heap.Init(&h.list)
}

// Len 获取堆的当前长度
func (h *Heap) Len() int {
	return h.list.Len()
}

// Push 元素入堆
func (h *Heap) Push(e Element) {
	if e.Index() >= 0 {
		panic("Element in other Heap")
	}
	heap.Push(&h.list, e)
}

// Pop 堆顶元素出堆
func (h *Heap) Pop() Element {
	if h.list.Len() <= 0 {
		panic("empty Heap")
	}
	x := heap.Pop(&h.list)
	h.adjustCap()
	return x.(Element)
}

// Remove 移除堆中指定位置的元素
func (h *Heap) Remove(i int) Element {
	if i < 0 || i >= h.list.Len() {
		panic("index out of range")
	}
	x := heap.Remove(&h.list, i)
	h.adjustCap()
	return x.(Element)
}

// Fix 修正指定位置元素的位置
func (h *Heap) Fix(i int) {
	if i < 0 || i >= h.list.Len() {
		panic("index out of range")
	}
	heap.Fix(&h.list, i)
}

// Top 返回堆顶元素
func (h *Heap) Top() Element {
	if h.list.Len() <= 0 {
		panic("empty Heap")
	}
	x := h.list[0]
	return x.(Element)
}

// adjustCap 调整堆容量
func (h *Heap) adjustCap() {
	c := cap(h.list)
	if c <= h.minCap {
		return
	}
	n := h.list.Len()
	if n < c/4 {
		c = c / 2
		if c < h.minCap {
			c = h.minCap
		}
		list := make(heapList, n, c)
		copy(list, h.list)
		h.list = list
	}
}
