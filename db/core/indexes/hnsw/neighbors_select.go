package hnsw

import (
	"Vectory/db/core/indexes/utils"
	"container/heap"
)

func (h *Hnsw) selectNeighborsHeuristic(v *Vertex, candidates []utils.Element, m int) []uint64 {
	result := make([]uint64, 0, m)

	workingQ := utils.NewMinHeapFromSliceDeep(candidates, cap(candidates))

	visited := newSet[uint64]()
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

func (h *Hnsw) selectNeighborsSimple(_ *Vertex, candidates []utils.Element, m int) []uint64 {
	size := m
	if len(candidates) < size {
		size = len(candidates)
	}

	minHeap := utils.NewMinHeapFromSliceDeep(candidates, cap(candidates))
	neighbors := make([]uint64, 0, size)

	for i := 0; i < size; i++ {
		neighbors = append(neighbors, heap.Pop(minHeap).(utils.Element).Id)
	}

	return neighbors
}
