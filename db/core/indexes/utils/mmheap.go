package utils

import (
	"container/heap"
	"math"
)

// taken from https://github.com/dogmatiq/kyu

func NewMinMaxHeapFromSlice(s []Element) *Heap {
	h := Heap{
		Elements: s,
		less: func(elements []Element, i, j int) bool {
			return elements[i].Distance < elements[j].Distance
		},
	}

	Init(&h)

	return &h
}

// Init establishes the heap invariants required by the other routines in this
// package.
//
// Init is idempotent with respect to the heap invariants and may be called
// whenever the heap invariants may have been invalidated.
func Init(h heap.Interface) {
	n := h.Len()
	for i := n / 2; i > 0; i-- {
		down(h, i, n)
	}
}

// Push pushes the element x onto the heap.
func Push(h heap.Interface, x interface{}) {
	i := h.Len()
	h.Push(x)
	up(h, i, i+1)
}

// Pop removes and returns the minimum element (according to Less) from the
// heap.
func Pop(h heap.Interface) interface{} {
	i := h.Len() - 1
	h.Swap(0, i)
	down(h, 0, i)

	return h.Pop()
}

// PopMax removes and returns the maximum element (according to Less) from the
// heap.
func PopMax(h heap.Interface) interface{} {
	n := h.Len()
	i := indexOfMax(h, n)
	j := n - 1

	h.Swap(i, j)
	down(h, i, j)

	return h.Pop()
}

// Remove removes and returns the element at index i from the heap.
func Remove(h heap.Interface, i int) interface{} {
	n := h.Len()
	j := n - 1

	if i != j {
		h.Swap(i, j)
		if !down(h, i, j) {
			up(h, i, j)
		}
	}

	return h.Pop()
}

// Fix re-establishes the heap ordering after the element at index i has changed
// its value.
//
// Changing the value of the element at index i and then calling Fix is
// equivalent to, but less expensive than, calling Remove(h, i) followed by a
// Push of the new value.
func Fix(h heap.Interface, i int) {
	n := h.Len()
	if !down(h, i, n) {
		up(h, i, n)
	}
}

// Max returns index of the maximum element in h.
//
// Unlike the minimum element, which is the root node at index 0, the maximum
// element may be either the left-hand or right-hand child of the root node.
//
// If h is empty it returns -1.
//
// The complexity is O(1).
func Max(h heap.Interface) int {
	return indexOfMax(h, h.Len())
}

// indexOfMax returns index of the maximum element in h.
func indexOfMax(h heap.Interface, n int) int {
	if n <= 2 {
		return n - 1
	}

	if h.Less(2, 1) {
		return 1
	}

	return 2
}

// isOnMinLevel returns true if the depth of i within the heap is an even
// number, indicating a "min" level, as opposed to a "max" level.
func isOnMinLevel(i int) bool {
	n := int(math.Log2(float64(i) + 1))
	return n&1 == 0
}

type compare func(h heap.Interface, i, j int) bool

// isLess returns true if heap[i] < heap[j].
func isLess(h heap.Interface, i, j int) bool {
	return h.Less(i, j)
}

// isGreater returns true if heap[i] > heap[j].
func isGreater(h heap.Interface, i, j int) bool {
	return h.Less(j, i)
}

// swapIf swaps the elements at i and j if the element at i is less than the
// element at j.
//
// It returns true if the elements are swapped.
func swapIf(h heap.Interface, less compare, i, j int) bool {
	if less(h, i, j) {
		h.Swap(i, j)
		return true
	}

	return false
}

// up moves the element at i upwards within the heap until it occupies an
// appropriate node.
func up(h heap.Interface, i, n int) {
	parent := (i - 1) / 2

	if isOnMinLevel(i) {
		if i > 0 && swapIf(h, isGreater, i, parent) {
			upX(h, isGreater, parent, n)
		} else {
			upX(h, isLess, i, n)
		}
	} else {
		if i > 0 && swapIf(h, isLess, i, parent) {
			upX(h, isLess, parent, n)
		} else {
			upX(h, isGreater, i, n)
		}
	}
}

func upX(h heap.Interface, less compare, i, n int) {
	for i > 2 {
		grandparent := (((i - 1) / 2) - 1) / 2

		if !swapIf(h, less, i, grandparent) {
			return
		}

		i = grandparent
	}
}

// down moves the element at i downards within the heap until it occupies an
// appropriate node.
func down(h heap.Interface, i, n int) bool {
	if isOnMinLevel(i) {
		return downX(h, isLess, i, n)
	}

	return downX(h, isGreater, i, n)
}

func downX(h heap.Interface, less compare, i, n int) bool {
	recursed := false

	for {
		m := minDescendent(h, less, i, n)
		if m == -1 {
			// i has no children.
			return recursed
		}

		parent := (m - 1) / 2

		if i == parent {
			// m is a direct child of i.
			swapIf(h, less, m, i)
			return recursed
		}

		// m is a grandchild of i.
		if !swapIf(h, less, m, i) {
			return recursed
		}

		swapIf(h, less, parent, m)

		i = m
		recursed = true
	}
}

// minDescendent returns the index of the smallest child or grandchild of i.
//
// It returns -1 if i is a leaf node.
func minDescendent(h heap.Interface, less compare, i, n int) int {
	// check i's left-hand child.
	left := i*2 + 1
	if left >= n {
		// i has no children.
		return -1
	}

	// check i's right-hand child.
	right := left + 1
	min, done := least(h, less, left, right, n)
	if done {
		return min
	}

	// check i's left-hand child's own left-hand child.
	min, done = least(h, less, min, left*2+1, n)
	if done {
		return min
	}

	// check i's left-hand child's right-hand child.
	min, done = least(h, less, min, left*2+2, n)
	if done {
		return min
	}

	// check i's right-hand child's right-hand child.
	min, done = least(h, less, min, right*2+1, n)
	if done {
		return min
	}

	// check i's right-hand child's own right-hand child.
	min, _ = least(h, less, min, right*2+2, n)

	return min
}

// least returns the index of the smaller element of those elements at i and j.
//
// If j overruns the heap, done is true.
func least(h heap.Interface, less compare, i, j, n int) (_ int, done bool) {
	if j >= n {
		return i, true
	}

	if less(h, i, j) {
		return i, false
	}

	return j, false
}
