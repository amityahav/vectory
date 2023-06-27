package disk_ann

import (
	"Vectory/indexes"
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

	currId uint32

	readOnly     bool
	snapshotPath string
}

func newMemoryIndex(currId uint32, deletedIds *sync.Map) *MemoryIndex {
	mi := MemoryIndex{
		graph:             nil,
		calculateDistance: nil,
		deletedIds:        deletedIds,
		currId:            currId,
		readOnly:          false,
		snapshotPath:      "",
	}

	return &mi
}

// TODO: beam search support as well
func (mi *MemoryIndex) Search(q []float32, k int) ([]uint32, []uint32) {
	return nil, nil
}

func (mi *MemoryIndex) Insert(v []float32) error {
	mi.Lock() // TODO: optimize locking
	defer mi.Unlock()

	if mi.readOnly {
		return ErrReadOnlyIndex
	}

	_, visited := mi.Search(v, 1)

	mi.currId += 1

	vertex := Vertex{
		Id:     mi.currId,
		Vector: v,
	}

	mi.graph.addVertex(&vertex)

	vertex.Neighbors = mi.RobustPrune(&vertex, visited)

	for _, n := range vertex.Neighbors {
		neighbor := mi.graph.vertices[n]

		neighbor.Neighbors = append(neighbor.Neighbors, vertex.Id)

		if len(neighbor.Neighbors) > mi.graph.maxDegree {
			neighbor.Neighbors = mi.RobustPrune(neighbor, neighbor.Neighbors)
		}
	}

	return nil
}

func (mi *MemoryIndex) Delete(id uint32) error {
	// TODO: delete consolidation?
	return nil
}

func (mi *MemoryIndex) RobustPrune(v *Vertex, candidates []uint32) []uint32 {
	// TODO locking
	candidateDistances := make([]indexes.Element, 0, len(candidates)+len(v.Neighbors))
	deletedCandidates := map[uint32]struct{}{}

	for _, c := range candidates {
		if c == v.Id {
			continue
		}

		cVertex := mi.graph.vertices[c]
		e := indexes.Element{
			Id:       int64(c),
			Distance: mi.calculateDistance(v.Vector, cVertex.Vector),
		}

		candidateDistances = append(candidateDistances, e)
	}

	for _, n := range v.Neighbors {
		nVertex := mi.graph.vertices[n]
		e := indexes.Element{
			Id:       int64(n),
			Distance: mi.calculateDistance(v.Vector, nVertex.Vector),
		}

		candidateDistances = append(candidateDistances, e)
	}

	candidatesHeap := indexes.NewMinHeapFromSlice(candidateDistances)
	newNeighbors := make([]uint32, 0, mi.graph.maxDegree)

	for candidatesHeap.Len() != 0 {
		min := heap.Pop(candidatesHeap).(indexes.Element)
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
