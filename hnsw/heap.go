package hnsw

import "container/heap"

type element struct {
	id       int64
	distance float32
}

type Heap struct {
	elements []element
	less     func(elements []element, i, j int) bool
}

func (h Heap) Len() int { return len(h.elements) }

func (h Heap) Swap(i, j int) { h.elements[i], h.elements[j] = h.elements[j], h.elements[i] }

func (h Heap) Less(i, j int) bool { return h.less(h.elements, i, j) }

func (h *Heap) Push(x any) {
	h.elements = append(h.elements, x.(element))
}

func (h *Heap) Pop() any {
	old := h.elements
	n := len(old)
	x := old[n-1]
	h.elements = old[0 : n-1]
	return x
}

func (h *Heap) Peek() any {
	return h.elements[h.Len()-1]
}

func NewMinHeap(capacity int) *Heap {
	return &Heap{
		elements: make([]element, capacity),
		less: func(elements []element, i, j int) bool {
			return elements[i].distance < elements[j].distance
		},
	}
}

func NewMinHeapFromSlice(s []element) *Heap {
	h := Heap{
		elements: s,
		less: func(elements []element, i, j int) bool {
			return elements[i].distance < elements[j].distance
		},
	}

	heap.Init(&h)

	return &h
}

func NewMinHeapFromSliceDeep(s []element, capacity int) *Heap {
	h := Heap{
		less: func(elements []element, i, j int) bool {
			return elements[i].distance < elements[j].distance
		},
	}

	h.elements = make([]element, capacity)
	copy(h.elements, s)

	heap.Init(&h)

	return &h
}

func NewMaxHeap(capacity int) *Heap {
	return &Heap{
		elements: make([]element, capacity),
		less: func(elements []element, i, j int) bool {
			return elements[i].distance > elements[j].distance
		},
	}
}

func NewMaxHeapFromSliceDeep(s []element, capacity int) *Heap {
	h := Heap{
		less: func(elements []element, i, j int) bool {
			return elements[i].distance > elements[j].distance
		},
	}

	h.elements = make([]element, capacity)
	copy(h.elements, s)

	heap.Init(&h)

	return &h
}
