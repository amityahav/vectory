package hnsw

import (
	"Vectory/hnsw/distance"
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"sync"
)

type Hnsw struct {
	hnswConfig

	entrypointID    int64
	currentMaxLayer int64

	nodes           map[int64]*Vertex
	distFunc        func([]float32, []float32) float32
	selectNeighbors func(*Vertex, []element, int, int64, bool, bool) []int64

	initialInsertion *sync.Once
}

type hnswConfig struct {
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

	heuristic bool

	distanceType string

	// Flags for the heuristic neighbors selection
	extendCandidates      bool
	keepPrunedConnections bool
}

func NewHnsw(config hnswConfig) *Hnsw {
	hnsw := Hnsw{
		hnswConfig:       config,
		entrypointID:     0,
		currentMaxLayer:  0,
		nodes:            make(map[int64]*Vertex),
		initialInsertion: &sync.Once{},
	}

	hnsw.setSelectNeighborsFunc(hnsw.heuristic)
	hnsw.setDistanceFunction(hnsw.distanceType)

	return &hnsw
}

func (h *Hnsw) setDistanceFunction(distanceType string) {
	switch distanceType {
	case distance.DotProduct:
		h.distFunc = distance.Dot
	}
}

func (h *Hnsw) setSelectNeighborsFunc(heuristic bool) {
	h.selectNeighbors = h.selectNeighborsSimple

	if heuristic {
		h.selectNeighbors = h.selectNeighborsHeuristic
	}
}

func (h *Hnsw) KnnSearch(v *Vertex, k, ef int) []int64 {
	var currentNearestElements []element

	res := make([]int64, 0, k)

	entrypointID := h.entrypointID
	currentMaxLayer := h.currentMaxLayer

	dist := h.calculateDistance(h.nodes[entrypointID].vector, v.vector)

	eps := make([]element, 0, 1)
	eps = append(eps, element{id: entrypointID, distance: dist})

	for l := currentMaxLayer; l > 0; l-- {
		currentNearestElements = h.searchLayer(v, eps, 1, l)
		eps[0] = currentNearestElements[0]
	}

	currentNearestElements = h.searchLayer(v, eps, ef, 0)

	minHeap := NewMinHeapFromSlice(currentNearestElements)

	for i := 0; i < k; i++ {
		res = append(res, heap.Pop(minHeap).(element).id)
	}

	return res
}

func (h *Hnsw) Insert(v *Vertex) error {
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

	h.nodes[v.id] = v

	if first {
		return nil
	}

	var nearestNeighbors []element

	dist := h.calculateDistance(h.nodes[h.entrypointID].vector, v.vector)

	eps := make([]element, 0, 1)
	eps = append(eps, element{id: h.entrypointID, distance: dist})

	currentMaxLayer := h.currentMaxLayer
	vertexLayer := h.calculateLevelForVertex()

	v.Init(vertexLayer+1, h.mMax, h.mMax0)

	// Lookup Phase
	for l := currentMaxLayer; l > vertexLayer; l-- {
		nearestNeighbors = h.searchLayer(v, eps, 1, l)
		eps[0] = nearestNeighbors[0]
	}

	// Construction Phase
	maxConn := h.mMax

	for l := min(currentMaxLayer, vertexLayer); l >= 0; l-- {
		nearestNeighbors = h.searchLayer(v, eps, h.efConstruction, l)
		neighbors := h.selectNeighbors(v, nearestNeighbors, h.m, l, h.extendCandidates, h.keepPrunedConnections)

		v.AddConnections(l, neighbors)

		if l == 0 {
			maxConn = h.mMax0
		}

		for _, n := range neighbors {
			nVertex := h.nodes[n]
			nVertex.AddConnection(l, v.id)

			if len(nVertex.GetConnections(l)) > maxConn {
				elems := make([]element, 0, len(nVertex.GetConnections(l))) // TODO: can be optimized size if we can estimate the num of connections total when extendingNeighbors

				for _, nn := range nVertex.GetConnections(l) {
					elems = append(elems, element{id: nn, distance: h.calculateDistance(v.vector, h.nodes[nn].vector)})
				}

				newNeighbors := h.selectNeighbors(v, elems, maxConn, l, h.extendCandidates, h.keepPrunedConnections)

				nVertex.SetConnections(l, newNeighbors)
			}

		}

		eps = nearestNeighbors
	}

	if vertexLayer > currentMaxLayer {
		h.entrypointID = v.id
		h.currentMaxLayer = vertexLayer
	}

	return nil
}

func (h *Hnsw) searchLayer(v *Vertex, eps []element, ef int, level int64) []element {
	visited := NewSet[int64]()
	for _, e := range eps {
		visited.Add(e.id)
	}

	candidates := NewMinHeapFromSliceDeep(eps, ef+1)
	nearestNeighbors := NewMaxHeapFromSliceDeep(eps, ef+1)

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
	result := make([]int64, 0, m)

	workingQ := NewMinHeapFromSliceDeep(candidates, cap(candidates))

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

					heap.Push(workingQ, nElem) // TODO: estimate fixed slice size?
				}
			}
		}
	}

	var discards []element // TODO: figure out optimal fixed size

	for workingQ.Len() > 0 && len(result) < m {
		e := heap.Pop(workingQ).(element)

		if len(result) == 0 {
			result = append(result, e.id)
		}

		flag := true
		for _, r := range result {
			if h.distFunc(h.nodes[e.id].vector, h.nodes[r].vector) < e.distance {
				flag = false
				break
			}
		}

		if flag {
			result = append(result, e.id)
		} else {
			discards = append(discards, e)
		}
	}

	if keepPruned {
		discardedHeap := NewMinHeapFromSlice(discards)
		for discardedHeap.Len() > 0 && len(result) < m {
			result = append(result, heap.Pop(discardedHeap).(element).id)
		}
	}

	return result
}

func (h *Hnsw) selectNeighborsSimple(_ *Vertex, candidates []element, m int, _ int64, _, _ bool) []int64 {
	size := m
	if len(candidates) < size {
		size = len(candidates)
	}

	minHeap := NewMinHeapFromSliceDeep(candidates, cap(candidates))
	neighbors := make([]int64, 0, size)

	for i := 0; i < size; i++ {
		neighbors = append(neighbors, heap.Pop(minHeap).(element).id)
	}

	return neighbors
}

func (h *Hnsw) calculateDistance(v1, v2 []float32) float32 {
	return h.distFunc(v1, v2)
}

func (h *Hnsw) insertFirstVertex(v *Vertex) error {
	v.Init(1, -1, h.mMax0)
	h.entrypointID = v.id

	return nil
}

func (h *Hnsw) isEmpty() bool {
	return len(h.nodes) == 0
}

func (h *Hnsw) calculateLevelForVertex() int64 {
	return int64(math.Floor(-math.Log(rand.Float64()) * h.mL))
}
