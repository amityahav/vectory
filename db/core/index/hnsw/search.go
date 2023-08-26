package hnsw

import (
	"Vectory/db/core/index/utils"
	"container/heap"
)

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

	var i int
	for minHeap.Len() > 0 {
		if i == k {
			break
		}

		e := heap.Pop(minHeap).(utils.Element)
		if _, ok := h.deletedNodes[e.Id]; ok {
			continue
		}

		res = append(res, e)
		i++
	}

	return res
}

func (h *Hnsw) searchLayer(v *Vertex, eps []utils.Element, ef int, level int64) []utils.Element {
	visited := newSet[uint64]()
	for _, e := range eps {
		visited.Add(e.Id)
	}

	candidates := utils.NewMinHeapFromSliceDeep(eps, ef+1)
	nearestNeighbors := utils.NewMaxHeapFromSliceDeep(eps, ef+1)
	connections := make([]uint64, h.mMax0) // reused for all candidates

	for candidates.Len() > 0 {
		c := heap.Pop(candidates).(utils.Element)
		f := nearestNeighbors.Peek().(utils.Element)

		if c.Distance > f.Distance {
			break
		}

		h.RLock()
		cVertex := h.nodes[c.Id]
		h.RUnlock()

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
