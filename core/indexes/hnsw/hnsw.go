package hnsw

import (
	distance2 "Vectory/core/indexes/distance"
	"Vectory/core/indexes/utils"
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
	selectNeighbors func(*Vertex, []utils.Element, int, int64, bool, bool) []int64

	initialInsertion *sync.Once
	mu               sync.RWMutex
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
	case distance2.DotProduct:
		h.distFunc = distance2.Dot
	case distance2.Euclidean:
		h.distFunc = distance2.EuclideanDistance
	}
}

func (h *Hnsw) setSelectNeighborsFunc(heuristic bool) {
	h.selectNeighbors = h.selectNeighborsSimple

	if heuristic {
		h.selectNeighbors = h.selectNeighborsHeuristic
	}
}

func (h *Hnsw) Search(v *Vertex, k, ef int) []int64 {
	h.mu.RLock()
	defer h.mu.RUnlock()

	var currentNearestElements []utils.Element

	res := make([]int64, 0, k)

	entrypointID := h.entrypointID
	epVertex := h.nodes[entrypointID]
	currentMaxLayer := h.currentMaxLayer

	dist := h.calculateDistance(epVertex.Vector, v.Vector)

	eps := make([]utils.Element, 0, 1)
	eps = append(eps, utils.Element{Id: entrypointID, Distance: dist})

	for l := currentMaxLayer; l > 0; l-- {
		currentNearestElements = h.searchLayer(v, eps, 1, l)
		eps[0] = currentNearestElements[0]
	}

	currentNearestElements = h.searchLayer(v, eps, ef, 0)

	minHeap := utils.NewMinHeapFromSlice(currentNearestElements)

	for i := 0; i < k && minHeap.Len() > 0; i++ {
		res = append(res, heap.Pop(minHeap).(utils.Element).Id)
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

	h.mu.Lock()
	defer h.mu.Unlock()

	h.nodes[v.Id] = v

	if first {
		return nil
	}

	entrypointID := h.entrypointID
	epVertex := h.nodes[entrypointID]
	currentMaxLayer := h.currentMaxLayer

	vertexLayer := h.calculateLevelForVertex()
	dist := h.calculateDistance(epVertex.Vector, v.Vector)

	var nearestNeighbors []utils.Element

	eps := make([]utils.Element, 0, 1)
	eps = append(eps, utils.Element{Id: entrypointID, Distance: dist})

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
			nVertex.AddConnection(l, v.Id)
			connections := nVertex.GetConnections(l)

			// pruning connections of neighbour if necessary
			if len(connections) > maxConn {
				elems := make([]utils.Element, 0, len(connections)) // TODO: can be optimized size if we can estimate the num of connections total when extendingNeighbors

				for _, nn := range connections {
					nnVertex := h.nodes[nn]

					elems = append(elems, utils.Element{Id: nn, Distance: h.calculateDistance(nVertex.Vector, nnVertex.Vector)})
				}

				newNeighbors := h.selectNeighbors(nVertex, elems, maxConn, l, h.extendCandidates, h.keepPrunedConnections)

				nVertex.SetConnections(l, newNeighbors)
			}
		}

		eps = nearestNeighbors
	}

	if vertexLayer > currentMaxLayer {
		h.entrypointID = v.Id
		h.currentMaxLayer = vertexLayer
	}

	return nil
}

func (h *Hnsw) searchLayer(v *Vertex, eps []utils.Element, ef int, level int64) []utils.Element {
	visited := NewSet[int64]()
	for _, e := range eps {
		visited.Add(e.Id)
	}

	candidates := utils.NewMinHeapFromSliceDeep(eps, ef+1)
	nearestNeighbors := utils.NewMaxHeapFromSliceDeep(eps, ef+1)

	for candidates.Len() > 0 {
		c := heap.Pop(candidates).(utils.Element)
		f := nearestNeighbors.Peek().(utils.Element)

		if c.Distance > f.Distance {
			break
		}

		cVertex := h.nodes[c.Id]
		for _, nid := range cVertex.GetConnections(level) {
			if !visited.Contains(nid) {
				visited.Add(nid)

				f = nearestNeighbors.Peek().(utils.Element)
				neighbour := h.nodes[nid]

				dist := h.calculateDistance(neighbour.Vector, v.Vector)
				if dist < f.Distance || nearestNeighbors.Len() < ef {
					e := utils.Element{Id: nid, Distance: dist}

					heap.Push(candidates, e)
					heap.Push(nearestNeighbors, e)

					if nearestNeighbors.Len() > ef {
						heap.Pop(nearestNeighbors)
					}
				}
			}
		}
	}

	return nearestNeighbors.Elements
}

func (h *Hnsw) selectNeighborsHeuristic(v *Vertex, candidates []utils.Element, m int, level int64, extendCandidates, keepPruned bool) []int64 {
	result := make([]int64, 0, m)

	workingQ := utils.NewMinHeapFromSliceDeep(candidates, cap(candidates))

	visited := NewSet[int64]()
	for _, c := range candidates {
		visited.Add(c.Id)
	}

	if extendCandidates {
		for _, c := range candidates {
			cVertex := h.nodes[c.Id]

			for _, n := range cVertex.GetConnections(level) {
				if !visited.Contains(n) && v.Id != n {
					visited.Add(n)
					nVertex := h.nodes[n]

					nElem := utils.Element{
						Id:       n,
						Distance: h.calculateDistance(v.Vector, nVertex.Vector),
					}

					heap.Push(workingQ, nElem) // TODO: estimate fixed slice size?
				}
			}
		}
	}

	var discards []utils.Element // TODO: figure out optimal fixed size

	for workingQ.Len() > 0 && len(result) < m {
		e := heap.Pop(workingQ).(utils.Element)

		flag := true
		for _, r := range result {
			eVertex := h.nodes[e.Id]
			rVertex := h.nodes[r]

			if h.distFunc(eVertex.Vector, rVertex.Vector) < e.Distance {
				flag = false
				break
			}
		}

		if flag {
			result = append(result, e.Id)
		} else {
			discards = append(discards, e)
		}
	}

	if keepPruned {
		discardedHeap := utils.NewMinHeapFromSlice(discards)
		for discardedHeap.Len() > 0 && len(result) < m {
			result = append(result, heap.Pop(discardedHeap).(utils.Element).Id)
		}
	}

	return result
}

func (h *Hnsw) selectNeighborsSimple(_ *Vertex, candidates []utils.Element, m int, _ int64, _, _ bool) []int64 {
	size := m
	if len(candidates) < size {
		size = len(candidates)
	}

	minHeap := utils.NewMinHeapFromSliceDeep(candidates, cap(candidates))
	neighbors := make([]int64, 0, size)

	for i := 0; i < size; i++ {
		neighbors = append(neighbors, heap.Pop(minHeap).(utils.Element).Id)
	}

	return neighbors
}

func (h *Hnsw) calculateDistance(v1, v2 []float32) float32 {
	return h.distFunc(v1, v2)
}

func (h *Hnsw) insertFirstVertex(v *Vertex) error {
	v.Init(1, -1, h.mMax0)
	h.entrypointID = v.Id

	return nil
}

func (h *Hnsw) isEmpty() bool {
	return len(h.nodes) == 0
}

func (h *Hnsw) calculateLevelForVertex() int64 {
	return int64(math.Floor(-math.Log(rand.Float64()) * h.mL))
}
