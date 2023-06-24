package disk_ann

import (
	"sync"
)

type MemoryIndex struct {
	sync.RWMutex
	graph             *Graph
	calculateDistance func([]float32, float32) float32

	// search list size
	l int

	// distance threshold
	a int

	// starting point
	s uint32

	lastInsertedId uint32

	readOnly     bool
	snapshotPath string
}

func newMemoryIndex() *MemoryIndex {
	mi := MemoryIndex{
		graph:             nil,
		calculateDistance: nil,
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

	mi.lastInsertedId += 1

	vertex := Vertex{
		Id:     mi.lastInsertedId,
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

func (mi *MemoryIndex) Delete(deleteList []uint32) error {
	return nil
}

func (mi *MemoryIndex) RobustPrune(v *Vertex, candidates []uint32) []uint32 {
	return nil
}

func (mi *MemoryIndex) Snapshot() {

}
