package hnsw

import (
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"sync"
)

type Hnsw struct {
	// Number of established connections
	m int

	// Maximum number of connections for each element per layer
	mMax int

	// Maximum number of connections for each element at layer zero
	mMax0 int

	// Size of the dynamic candidate list
	efConstruction int

	// Normalization factor for level generation
	mL float64

	entrypointID    int64
	currentMaxLayer int64

	nodes    map[int64]*Vertex
	distFunc func([]float32, []float32) float32

	initialInsertion *sync.Once
}

func (h *Hnsw) insert(v *Vertex) error {
	var (
		first bool
		err   error
	)

	h.initialInsertion.Do(func() {
		if h.isEmpty() {
			err = h.insertFirstVertex(v)
			if err != nil {
				return
			}

			first = true
		}

	})

	if err != nil {
		return fmt.Errorf("initialInsertion: %s", err.Error())
	}

	if first {
		return nil
	}

	var nearestNeighbors heap.Interface

	dist := h.calculateDistance(h.nodes[h.entrypointID].vector, v.vector)
	eps := []element{{
		id:       h.entrypointID,
		distance: dist,
	}}
	currentMaxLayer := h.currentMaxLayer
	vertexLayer := h.calculateLevelForVertex()

	v.Init(vertexLayer, h.mMax, h.mMax0)

	// Lookup Phase
	for i := currentMaxLayer; i > vertexLayer; i-- {
		nearestNeighbors = h.searchLayer(v, eps, 1, i)
		eps[0] = nearestNeighbors.Pop().(element)
	}

	// Construction Phase
	for i := min(currentMaxLayer, vertexLayer); i >= 0; i-- {

	}

	return nil
}

// TODO: test
func (h *Hnsw) searchLayer(v *Vertex, eps []element, ef int, level int64) heap.Interface {
	visited := NewSet[int64]()
	for _, e := range eps {
		visited.Add(e.id)
	}

	candidates := NewMinHeapFromSlice(eps)
	nearestNeighbors := NewMaxHeapFromSlice(eps)

	for candidates.Len() > 0 {
		c := candidates.Pop().(element)
		f := nearestNeighbors.Peek().(element)

		if c.distance > f.distance {
			break
		}

		cVertex := h.nodes[c.id]
		for _, nid := range cVertex.GetConnections(level) {
			if !visited.Contains(nid) {
				visited.Add(nid)

				f = nearestNeighbors.Peek().(element)
				neighbour := h.nodes[nid]

				dist := h.calculateDistance(neighbour.vector, v.vector)
				if dist < f.distance || nearestNeighbors.Len() < ef {
					e := element{id: nid, distance: dist}

					candidates.Push(e)
					nearestNeighbors.Push(e)

					if nearestNeighbors.Len() > ef {
						nearestNeighbors.Pop()
					}
				}
			}
		}

	}

	return nearestNeighbors
}

// TODO: implement
func (h *Hnsw) calculateDistance(v1, v2 []float32) float32 {
	return h.distFunc(v1, v2)
}

func (h *Hnsw) insertFirstVertex(v *Vertex) error {
	v.Init(1, -1, h.mMax0)
	h.entrypointID = v.id
	h.nodes[v.id] = v

	return nil
}

func (h *Hnsw) isEmpty() bool {
	return len(h.nodes) > 0
}

func (h *Hnsw) calculateLevelForVertex() int64 {
	return int64(math.Floor(-math.Log(rand.Float64()) * h.mL))
}