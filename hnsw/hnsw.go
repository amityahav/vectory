package hnsw

import (
	"math"
	"math/rand"
)

type Hnsw struct {
	// Number of established connections
	m int64

	// Maximum number of connections for each element per layer
	mMax int64

	// Maximum number of connections for each element at layer zero
	mMax0 int64

	// Size of the dynamic candidate list
	efConstruction int64

	// Normalization factor for level generation
	mL float64

	entrypointID    int64
	currentMaxLayer int64

	nodes map[int64]Vertex
}

func (h *Hnsw) insert(vec *Vertex) error {
	return nil
}

func (h *Hnsw) calculateLevelForVector() int64 {
	return int64(math.Floor(-math.Log(rand.Float64()) * h.mL))
}
