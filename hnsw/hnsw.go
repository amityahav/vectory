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

	// Flags for the heuristic neighbors selection
	extendCandidates      bool
	keepPrunedConnections bool

	entrypointID    int64
	currentMaxLayer int64

	nodes           map[int64]*Vertex
	distFunc        func([]float32, []float32) float32
	selectNeighbors func(*Vertex, []element, int, int64, bool, bool) []int64

	initialInsertion *sync.Once
}

// TODO: currently just for testing
func NewHnsw() *Hnsw {
	hnsw := Hnsw{
		m:               0,
		mMax:            0,
		mMax0:           0,
		efConstruction:  0,
		mL:              0,
		entrypointID:    0,
		currentMaxLayer: 0,
		nodes:           make(map[int64]*Vertex),
	}

	selectedNeighborsSimple := true
	hnsw.SetSelectNeighborsFunc(selectedNeighborsSimple)

	return &hnsw
}

func (h *Hnsw) SetSelectNeighborsFunc(simple bool) {
	h.selectNeighbors = h.selectNeighborsHeuristic

	if simple {
		h.selectNeighbors = h.selectNeighborsSimple
	}
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

	var nearestNeighbors []element

	dist := h.calculateDistance(h.nodes[h.entrypointID].vector, v.vector)
	eps := []element{{
		id:       h.entrypointID,
		distance: dist,
	}}

	currentMaxLayer := h.currentMaxLayer
	vertexLayer := h.calculateLevelForVertex()

	v.Init(vertexLayer, h.mMax, h.mMax0)

	// Lookup Phase
	for l := currentMaxLayer; l > vertexLayer; l-- {
		nearestNeighbors = h.searchLayer(v, eps, 1, l)
		eps[0] = nearestNeighbors[0]
	}

	// Construction Phase
	for l := min(currentMaxLayer, vertexLayer); l >= 0; l-- {
		nearestNeighbors = h.searchLayer(v, eps, h.efConstruction, l)
		neighbors := h.selectNeighbors(v, nearestNeighbors, h.m, l, h.extendCandidates, h.keepPrunedConnections)

		v.SetConnections(l, neighbors)

		for _, n := range neighbors {

		}
	}

	return nil
}

// TODO: test
func (h *Hnsw) searchLayer(v *Vertex, eps []element, ef int, level int64) []element {
	visited := NewSet[int64]()
	for _, e := range eps {
		visited.Add(e.id)
	}

	candidates := NewMinHeapFromSlice(eps)
	nearestNeighbors := NewMaxHeapFromSlice(eps)

	for candidates.Len() > 0 {
		c := heap.Pop(candidates).(element)
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

					heap.Push(candidates, e)
					heap.Push(nearestNeighbors, e)

					if nearestNeighbors.Len() > ef {
						heap.Pop(nearestNeighbors)
					}
				}
			}
		}

	}

	return nearestNeighbors.elements
}

func (h *Hnsw) selectNeighborsHeuristic(v *Vertex, candidates []element, m int, level int64, extendCandidates, keepPruned bool) []int64 {
	var result []int64

	workingQ := NewMinHeapFromSlice(candidates)

	set := NewSet[int64]()

	if extendCandidates {
		for _, c := range candidates {
			set.Add(c.id)
			for _, n := range h.nodes[c.id].GetConnections(level) {
				if !set.Contains(n) {
					set.Add(n)
					nVertex := h.nodes[n]

					nElem := element{
						id:       n,
						distance: h.calculateDistance(v.vector, nVertex.vector),
					}

					heap.Push(workingQ, nElem)
				}
			}
		}
	}

	var discards []element

	for workingQ.Len() > 0 && len(result) < m {
		e := heap.Pop(workingQ).(element)

		if len(result) == 0 {
			result = append(result, e.id)
		}

		for _, _ = range result { // TODO: get elem
			if e.distance < 5 { // TODO: what should it be?
				result = append(result, e.id)
			} else {
				discards = append(discards, e)
			}
		}
	}

	if keepPruned {
		discardedHeap := NewMaxHeapFromSlice(discards)
		for discardedHeap.Len() > 0 && len(result) < m {
			result = append(result, heap.Pop(discardedHeap).(element).id)
		}
	}

	return result
}

func (h *Hnsw) selectNeighborsSimple(v *Vertex, candidates []element, m int, _ int64, _, _ bool) []int64 {
	size := m
	if len(candidates) < size {
		size = len(candidates)
	}

	minHeap := NewMinHeapFromSlice(candidates)
	neighbors := make([]int64, size)

	for i := 0; i < size; i++ {
		neighbors[i] = heap.Pop(minHeap).(element).id
	}

	return neighbors
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
