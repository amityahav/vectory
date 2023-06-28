package disk_ann

import (
	utils2 "Vectory/core/indexes/utils"
	"container/heap"
	"sync"
)

type MemoryIndex struct {
	sync.RWMutex
	graph             *Graph
	calculateDistance func([]float32, []float32) float32
	deletedIds        *sync.Map

	// search list size
	l int

	// distance threshold
	a int

	// starting point
	s uint32

	readOnly     bool
	snapshotPath string
}

func newMemoryIndex(deletedIds *sync.Map) *MemoryIndex {
	mi := MemoryIndex{
		graph:             nil,
		calculateDistance: nil,
		deletedIds:        deletedIds,
		readOnly:          false,
		snapshotPath:      "",
	}

	return &mi
}

// TODO: beam search support as well
func (mi *MemoryIndex) Search(q []float32, k int) ([]uint32, []utils2.Element) {
	// TODO: locking
	sVertex := mi.graph.vertices[mi.s]
	resultSet := utils2.NewMinMaxHeapFromSlice([]utils2.Element{{
		Id:       int64(mi.s),
		Distance: mi.calculateDistance(sVertex.Vector, q),
	}})
	visited := map[uint32]struct{}{}

	var candidatesVisited []utils2.Element

	for resultSet.Len() != 0 {
		min := utils2.Pop(resultSet).(utils2.Element)

		if _, ok := visited[uint32(min.Id)]; ok {
			continue
		}

		visited[uint32(min.Id)] = struct{}{}

		// filter deleted vertices from the result
		if _, ok := mi.deletedIds.Load(uint32(min.Id)); !ok {
			candidatesVisited = append(candidatesVisited, min)
		}

		minVertex := mi.graph.vertices[uint32(min.Id)]

		for _, n := range minVertex.Neighbors {
			nVertex := mi.graph.vertices[n]

			utils2.Push(resultSet, utils2.Element{
				Id:       int64(n),
				Distance: mi.calculateDistance(nVertex.Vector, q),
			})

			for resultSet.Len() > mi.l {
				utils2.PopMax(resultSet)
			}
		}
	}

	candidatesHeap := utils2.NewMinHeapFromSliceDeep(candidatesVisited, len(candidatesVisited))

	// TODO: this is unused in the paper do we need it?
	knn := make([]uint32, 0, k)
	for i := 0; i < k && candidatesHeap.Len() > 0; i++ {
		knn = append(knn, uint32(heap.Pop(candidatesHeap).(utils2.Element).Id))
	}

	return knn, candidatesVisited
}

func (mi *MemoryIndex) Insert(v []float32, currId, dataId uint32) error {
	mi.Lock() // TODO: optimize locking
	defer mi.Unlock()

	if mi.readOnly {
		return ErrReadOnlyIndex
	}

	_, visited := mi.Search(v, 1)

	vertex := Vertex{
		Id:     currId,
		Vector: v,
	}

	mi.graph.addVertex(&vertex)

	vertex.Neighbors = mi.RobustPrune(&vertex, visited)

	for _, n := range vertex.Neighbors {
		neighbor := mi.graph.vertices[n]

		neighbor.Neighbors = append(neighbor.Neighbors, vertex.Id)

		if len(neighbor.Neighbors) > mi.graph.maxDegree {
			distances := make([]utils2.Element, 0, len(neighbor.Neighbors))

			for _, nn := range neighbor.Neighbors {
				nnVertex := mi.graph.vertices[nn]
				distances = append(distances, utils2.Element{
					Id:       int64(nn),
					Distance: mi.calculateDistance(nnVertex.Vector, neighbor.Vector),
				})
			}
			neighbor.Neighbors = mi.RobustPrune(neighbor, distances)
		}
	}

	return nil
}

func (mi *MemoryIndex) Delete(id uint32) error {
	// TODO: delete consolidation?
	return nil
}

func (mi *MemoryIndex) RobustPrune(v *Vertex, candidates []utils2.Element) []uint32 {
	// TODO locking
	deletedCandidates := map[uint32]struct{}{v.Id: {}} // excluding vertex v

	for _, n := range v.Neighbors {
		nVertex := mi.graph.vertices[n]
		e := utils2.Element{
			Id:       int64(n),
			Distance: mi.calculateDistance(v.Vector, nVertex.Vector),
		}

		candidates = append(candidates, e)
	}

	candidatesHeap := utils2.NewMinHeapFromSlice(candidates)
	newNeighbors := make([]uint32, 0, mi.graph.maxDegree)

	for candidatesHeap.Len() != 0 {
		min := heap.Pop(candidatesHeap).(utils2.Element)
		if _, ok := deletedCandidates[uint32(min.Id)]; ok {
			continue
		}

		newNeighbors = append(newNeighbors, uint32(min.Id))

		if len(v.Neighbors) == mi.graph.maxDegree {
			break
		}

		for _, c := range candidatesHeap.Elements {
			minVertex := mi.graph.vertices[uint32(min.Id)]
			cVertex := mi.graph.vertices[uint32(c.Id)]

			if float32(mi.a)*mi.calculateDistance(minVertex.Vector, cVertex.Vector) <= c.Distance {
				deletedCandidates[cVertex.Id] = struct{}{}
			}
		}
	}

	return newNeighbors
}

func (mi *MemoryIndex) Snapshot() {

}

func (mi *MemoryIndex) ReadOnly() {
	mi.readOnly = true
}
