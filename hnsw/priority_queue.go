package hnsw

type element struct {
	id       int64
	distance int64
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

func NewMinHeap(capacity int) *Heap {
	return &Heap{
		elements: make([]element, capacity),
		less: func(elements []element, i, j int) bool {
			return elements[i].distance < elements[j].distance
		},
	}
}

func NewMaxHeap(capacity int) *Heap {
	return &Heap{
		elements: make([]element, capacity),
		less: func(elements []element, i, j int) bool {
			return elements[i].distance > elements[j].distance
		},
	}
}
