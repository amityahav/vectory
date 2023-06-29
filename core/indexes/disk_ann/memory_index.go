package disk_ann

import (
	"Vectory/core/indexes/utils"
	"container/heap"
	"sync"
)

type MemoryIndex struct {
	sync.RWMutex
	graph             *Graph
	calculateDistance func([]float32, []float32) float32
	deletedDataIds    *sync.Map

	// search list size
	l int

	// distance threshold
	a int

	// starting point
	s uint32

	readOnly     bool
	snapshotPath string
}

func newMemoryIndex(deletedDataIds *sync.Map) *MemoryIndex {
	mi := MemoryIndex{
		graph:             nil,
		calculateDistance: nil,
		deletedDataIds:    deletedDataIds,
		readOnly:          false,
		snapshotPath:      "",
	}

	return &mi
}

// TODO: beam search support as well
func (mi *MemoryIndex) Search(q []float32, k int, searchApi bool) ([]utils.Element, []utils.Element) {
	// TODO: locking
	sVertex := mi.graph.vertices[mi.s]
	e := utils.Element{
		Id:       int64(mi.s),
		Distance: mi.calculateDistance(sVertex.Vector, q),
	}

	if searchApi {
		e.DataId = sVertex.DataId
	}

	resultSet := utils.NewMinMaxHeapFromSlice([]utils.Element{e})
	visited := map[uint32]struct{}{}

	var candidatesVisited []utils.Element

	for resultSet.Len() != 0 {
		min := utils.Pop(resultSet).(utils.Element)

		if _, ok := visited[uint32(min.Id)]; ok {
			continue
		}

		visited[uint32(min.Id)] = struct{}{}

		minVertex := mi.graph.vertices[uint32(min.Id)]

		// filter deleted vertices from the result
		if _, ok := mi.deletedDataIds.Load(minVertex.DataId); !ok {
			candidatesVisited = append(candidatesVisited, min)
		}

		for _, n := range minVertex.Neighbors {
			nVertex := mi.graph.vertices[n]

			e = utils.Element{
				Id:       int64(n),
				Distance: mi.calculateDistance(nVertex.Vector, q),
			}

			if searchApi {
				e.DataId = nVertex.DataId
			}

			utils.Push(resultSet, e)

			for resultSet.Len() > mi.l {
				utils.PopMax(resultSet)
			}
		}
	}

	// K-NN is returned when search is performed
	if searchApi {
		candidatesHeap := utils.NewMinHeapFromSliceDeep(candidatesVisited, len(candidatesVisited))

		knn := make([]utils.Element, 0, k)
		for i := 0; i < k && candidatesHeap.Len() > 0; i++ {
			knn = append(knn, heap.Pop(candidatesHeap).(utils.Element))
		}

		return knn, nil
	}

	return nil, candidatesVisited
}

func (mi *MemoryIndex) Insert(v []float32, currId, dataId uint32) error {
	mi.Lock() // TODO: optimize locking
	defer mi.Unlock()

	if mi.readOnly {
		return ErrReadOnlyIndex
	}
	vertex := Vertex{
		Id:     currId,
		DataId: dataId,
		Vector: v,
	}

	mi.graph.addVertex(&vertex)

	_, visited := mi.Search(v, 1, false)

	vertex.Neighbors = mi.RobustPrune(&vertex, visited)

	for _, n := range vertex.Neighbors {
		neighbor := mi.graph.vertices[n]

		neighbor.Neighbors = append(neighbor.Neighbors, vertex.Id)

		if len(neighbor.Neighbors) > mi.graph.maxDegree {
			distances := make([]utils.Element, 0, len(neighbor.Neighbors))

			for _, nn := range neighbor.Neighbors {
				nnVertex := mi.graph.vertices[nn]
				distances = append(distances, utils.Element{
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

func (mi *MemoryIndex) RobustPrune(v *Vertex, candidates []utils.Element) []uint32 {
	// TODO locking
	deletedCandidates := map[uint32]struct{}{v.Id: {}} // excluding vertex v

	for _, n := range v.Neighbors {
		nVertex := mi.graph.vertices[n]
		e := utils.Element{
			Id:       int64(n),
			Distance: mi.calculateDistance(v.Vector, nVertex.Vector),
		}

		candidates = append(candidates, e)
	}

	candidatesHeap := utils.NewMinHeapFromSlice(candidates)
	newNeighbors := make([]uint32, 0, mi.graph.maxDegree)

	for candidatesHeap.Len() != 0 {
		min := heap.Pop(candidatesHeap).(utils.Element)
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
