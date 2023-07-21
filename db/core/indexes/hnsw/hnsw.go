package hnsw

import (
	"Vectory/db/core/indexes"
	"Vectory/db/core/indexes/distance"
	"Vectory/db/core/indexes/utils"
	"Vectory/entities"
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"sync"
)

var _ indexes.VectorIndex = &Hnsw{}

type Hnsw struct {
	sync.RWMutex
	m                int
	mMax             int
	mMax0            int
	efConstruction   int
	ef               int
	mL               float64
	entrypointID     int64
	currentMaxLayer  int64
	nodes            map[int64]*Vertex
	distFunc         func([]float32, []float32) float32
	selectNeighbors  func(*Vertex, []utils.Element, int) []int64
	initialInsertion *sync.Once
	curId            int64
}

func NewHnsw(params entities.HnswParams) *Hnsw {
	h := Hnsw{
		m:                params.M,
		mMax:             params.MMax,
		ef:               params.Ef,
		efConstruction:   params.EfConstruction,
		nodes:            make(map[int64]*Vertex), // TODO: change to an array
		initialInsertion: &sync.Once{},
	}

	h.mMax0 = 2 * h.mMax
	h.mL = 1 / math.Log(float64(h.m))

	switch params.DistanceType {
	case entities.DotProduct:
		h.distFunc = distance.Dot
	case entities.Euclidean:
		h.distFunc = distance.EuclideanDistance
	}

	h.selectNeighbors = h.selectNeighborsSimple

	if params.Heuristic {
		h.selectNeighbors = h.selectNeighborsHeuristic
	}

	return &h
}

func (h *Hnsw) Search(q []float32, k int) []utils.Element {
	var currentNearestElements []utils.Element

	dummy := Vertex{
		vector: q,
	}

	res := make([]utils.Element, 0, k)

	h.RLock()
	entrypointID := h.entrypointID
	epVertex := h.nodes[entrypointID]
	currentMaxLayer := h.currentMaxLayer
	h.RUnlock()

	dist := h.calculateDistance(epVertex.vector, q)

	eps := make([]utils.Element, 0, 1)
	eps = append(eps, utils.Element{Id: entrypointID, Distance: dist})

	for l := currentMaxLayer; l > 0; l-- {
		currentNearestElements = h.searchLayer(&dummy, eps, 1, l)
		eps[0] = currentNearestElements[0]
	}

	currentNearestElements = h.searchLayer(&dummy, eps, h.ef, 0)

	minHeap := utils.NewMinHeapFromSlice(currentNearestElements)

	for i := 0; i < k && minHeap.Len() > 0; i++ {
		res = append(res, heap.Pop(minHeap).(utils.Element))
	}

	return res
}

func (h *Hnsw) Delete(objId int64) bool {
	panic("implement me")
}

func (h *Hnsw) Insert(vector []float32, vectorId int64, objId uint64) error {
	var (
		first bool
		err   error
	)

	//nextId := atomic.AddInt64(&h.curId, 1)
	//currId := nextId - 1

	v := Vertex{
		id:     vectorId,
		vector: vector,
		objId:  objId,
	}

	h.initialInsertion.Do(func() {
		if h.isEmpty() {
			err = h.insertFirstVertex(&v)
			if err != nil {
				return
			}

			first = true
		}

	})

	if err != nil {
		return fmt.Errorf("initialInsertion: %s", err.Error())
	}

	h.Lock()
	h.nodes[v.id] = &v
	h.Unlock()

	if first {
		return nil
	}

	h.RLock()
	entrypointID := h.entrypointID
	epVertex := h.nodes[entrypointID]
	currentMaxLayer := h.currentMaxLayer
	h.RUnlock()

	vertexLayer := h.calculateLevelForVertex()
	dist := h.calculateDistance(epVertex.vector, v.vector)

	var nearestNeighbors []utils.Element

	eps := make([]utils.Element, 0, 1)
	eps = append(eps, utils.Element{Id: entrypointID, Distance: dist})

	v.Init(vertexLayer+1, h.mMax, h.mMax0)

	// Lookup Phase
	for l := currentMaxLayer; l > vertexLayer; l-- {
		nearestNeighbors = h.searchLayer(&v, eps, 1, l)
		eps[0] = nearestNeighbors[0]
	}

	// Construction Phase
	maxConn := h.mMax
	for l := min(currentMaxLayer, vertexLayer); l >= 0; l-- {
		nearestNeighbors = h.searchLayer(&v, eps, h.efConstruction, l)
		neighbors := h.selectNeighbors(&v, nearestNeighbors, h.m)

		v.SetConnections(l, neighbors)

		if l == 0 {
			maxConn = h.mMax0
		}

		for _, n := range neighbors {
			h.RLock()
			nVertex := h.nodes[n]
			h.RUnlock()

			nVertex.Lock()
			connections := nVertex.GetConnections(l)

			if len(connections) < maxConn {
				nVertex.AddConnection(l, v.id)
			} else { // pruning
				elems := make([]utils.Element, 0, len(connections)+1)

				elems = append(elems, utils.Element{
					Id:       v.id,
					Distance: h.calculateDistance(nVertex.vector, v.vector),
				})

				for _, nn := range connections {
					h.RLock()
					nnVertex := h.nodes[nn]
					h.RUnlock()

					elems = append(elems, utils.Element{Id: nn, Distance: h.calculateDistance(nVertex.vector, nnVertex.vector)})
				}

				newNeighbors := h.selectNeighbors(nVertex, elems, maxConn)

				nVertex.SetConnections(l, newNeighbors)
			}
			nVertex.Unlock()
		}

		eps = nearestNeighbors
	}

	h.Lock()
	if vertexLayer > currentMaxLayer {
		h.entrypointID = v.id
		h.currentMaxLayer = vertexLayer
	}
	h.Unlock()

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

		h.RLock()
		cVertex := h.nodes[c.Id]
		h.RUnlock()

		connections := make([]int64, h.mMax0) // reused for all candidates

		cVertex.Lock()
		connections = connections[:len(cVertex.GetConnections(level))]
		copy(connections, cVertex.GetConnections(level))
		cVertex.Unlock()

		for _, nid := range connections {
			if visited.Contains(nid) {
				continue
			}

			visited.Add(nid)

			f = nearestNeighbors.Peek().(utils.Element)

			h.RLock()
			neighbour := h.nodes[nid]
			h.RUnlock()

			dist := h.calculateDistance(neighbour.vector, v.vector)
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

	return nearestNeighbors.Elements
}

func (h *Hnsw) selectNeighborsHeuristic(v *Vertex, candidates []utils.Element, m int) []int64 {
	result := make([]int64, 0, m)

	workingQ := utils.NewMinHeapFromSliceDeep(candidates, cap(candidates))

	visited := NewSet[int64]()
	for _, c := range candidates {
		visited.Add(c.Id)
	}

	//connections := make([]int64, h.mMax0) // reused for all neighbors

	//if extendCandidates {
	//	for _, c := range candidates {
	//		h.RLock()
	//		cVertex := h.nodes[c.id]
	//		h.RUnlock()
	//
	//		//cVertex.Lock()
	//		connections = connections[:len(cVertex.GetConnections(level))]
	//		copy(connections, cVertex.GetConnections(level))
	//		//cVertex.Unlock()
	//
	//		for _, n := range connections {
	//			if !visited.Contains(n) && v.id != n {
	//				visited.Add(n)
	//
	//				h.RLock()
	//				nVertex := h.nodes[n]
	//				h.RUnlock()
	//
	//				nElem := utils.Element{
	//					id:       n,
	//					Distance: h.calculateDistance(v.vector, nVertex.vector),
	//				}
	//
	//				heap.Push(workingQ, nElem)
	//			}
	//		}
	//	}
	//}

	discards := make([]utils.Element, 0, workingQ.Len())

	for workingQ.Len() > 0 && len(result) < m {
		e := heap.Pop(workingQ).(utils.Element)

		flag := true
		for _, r := range result {
			h.RLock()
			eVertex := h.nodes[e.Id]
			rVertex := h.nodes[r]
			h.RUnlock()

			if h.distFunc(eVertex.vector, rVertex.vector) < e.Distance {
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

	//if keepPruned {
	//	discardedHeap := utils.NewMinHeapFromSlice(discards)
	//	for discardedHeap.Len() > 0 && len(result) < m {
	//		result = append(result, heap.Pop(discardedHeap).(utils.Element).id)
	//	}
	//}

	return result
}

func (h *Hnsw) selectNeighborsSimple(_ *Vertex, candidates []utils.Element, m int) []int64 {
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
	h.Lock()
	defer h.Unlock()

	v.Init(1, -1, h.mMax0)
	h.entrypointID = v.id
	h.currentMaxLayer = 0

	return nil
}

func (h *Hnsw) isEmpty() bool {
	return len(h.nodes) == 0
}

func (h *Hnsw) calculateLevelForVertex() int64 {
	return int64(math.Floor(-math.Log(rand.Float64()) * h.mL))
}
