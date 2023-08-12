package utils

import "container/heap"

type Element struct {
	Id       uint64
	Distance float32
}

type Heap struct {
	Elements []Element
	less     func(elements []Element, i, j int) bool
}

func (h Heap) Len() int { return len(h.Elements) }

func (h Heap) Swap(i, j int) { h.Elements[i], h.Elements[j] = h.Elements[j], h.Elements[i] }

func (h Heap) Less(i, j int) bool { return h.less(h.Elements, i, j) }

func (h *Heap) Push(x any) {
	h.Elements = append(h.Elements, x.(Element))
}

func (h *Heap) Pop() any {
	old := h.Elements
	n := len(old)
	x := old[n-1]
	h.Elements = old[0 : n-1]

	return x
}

func (h *Heap) Peek() any {
	return h.Elements[h.Len()-1]
}

func NewMinHeap(capacity int) *Heap {
	return &Heap{
		Elements: make([]Element, capacity),
		less: func(elements []Element, i, j int) bool {
			return elements[i].Distance < elements[j].Distance
		},
	}
}

func NewMinHeapFromSlice(s []Element) *Heap {
	h := Heap{
		Elements: s,
		less: func(elements []Element, i, j int) bool {
			return elements[i].Distance < elements[j].Distance
		},
	}

	heap.Init(&h)

	return &h
}

func NewMinHeapFromSliceDeep(s []Element, capacity int) *Heap {
	h := Heap{
		less: func(elements []Element, i, j int) bool {
			return elements[i].Distance < elements[j].Distance
		},
	}

	h.Elements = make([]Element, len(s), capacity)
	copy(h.Elements, s)

	heap.Init(&h)

	return &h
}

func NewMaxHeap(capacity int) *Heap {
	return &Heap{
		Elements: make([]Element, capacity),
		less: func(elements []Element, i, j int) bool {
			return elements[i].Distance > elements[j].Distance
		},
	}
}

func NewMaxHeapFromSliceDeep(s []Element, capacity int) *Heap {
	h := Heap{
		less: func(elements []Element, i, j int) bool {
			return elements[i].Distance > elements[j].Distance
		},
	}

	h.Elements = make([]Element, len(s), capacity)
	copy(h.Elements, s)

	heap.Init(&h)

	return &h
}
