package hnsw

import (
	"math"
	"math/rand"
)

func (h *Hnsw) calculateDistance(v1, v2 []float32) float32 {
	return h.distFunc(v1, v2)
}

func (h *Hnsw) isEmpty() bool {
	return len(h.nodes) == 0
}

func (h *Hnsw) calculateLevelForVertex() int64 {
	return int64(math.Floor(-math.Log(rand.Float64()) * h.mL))
}

func min(a, b int64) int64 {
	m := b
	if a < b {
		m = a
	}

	return m
}

type Set[T comparable] map[T]struct{}

func newSet[T comparable]() Set[T] {
	return Set[T]{}
}

func (s Set[T]) Add(elem T) {
	s[elem] = struct{}{}
	return
}

func (s Set[T]) Contains(elem T) bool {
	_, ok := s[elem]
	return ok
}
