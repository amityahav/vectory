package disk_ann

import (
	"Vectory/db/core/indexes/utils"
	"container/heap"
	"encoding/binary"
	"sync"
)

type MemoryIndex struct {
	sync.RWMutex
	graph             *Graph
	calculateDistance func([]float32, []float32) float32
	deletedObjIds     *sync.Map
	initialInsertion  *sync.Once

	// starting point
	s uint32

	// vectors dimension
	dim uint32

	// max vertex degree
	maxDegree uint32

	// the first vertex id that was inserted into the size
	firstId uint32

	readOnly     bool
	size         uint32 // immutable size once the index is snapshot
	snapshotPath string
}

func newMemoryIndex(distFunc func([]float32, []float32) float32,
	deletedObjIds *sync.Map, firstId uint32, maxDegree uint32, dim uint32) *MemoryIndex {
	mi := MemoryIndex{
		graph:             newGraph(),
		calculateDistance: distFunc,
		deletedObjIds:     deletedObjIds,
		initialInsertion:  &sync.Once{},
		firstId:           firstId,
		maxDegree:         maxDegree,
		dim:               dim,
		readOnly:          false,
	}

	return &mi
}

// TODO: beam search support as well
func (mi *MemoryIndex) search(q []float32, k int, listSize int, onlySearch bool) ([]utils.Element, []utils.Element) {
	// TODO: locking
	sVertex := mi.graph.vertices[mi.s]
	e := utils.Element{
		Id:       mi.s,
		Distance: mi.calculateDistance(sVertex.vector, q),
	}

	if onlySearch {
		e.ObjId = sVertex.objId
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

		// add non-deleted vertices to the result
		if _, ok := mi.deletedObjIds.Load(minVertex.objId); !ok {
			candidatesVisited = append(candidatesVisited, min)
		}

		for _, n := range minVertex.neighbors {
			nVertex := mi.graph.vertices[n]

			e = utils.Element{
				Id:       n,
				Distance: mi.calculateDistance(nVertex.vector, q),
			}

			if onlySearch {
				e.ObjId = nVertex.objId
			}

			utils.Push(resultSet, e)

			for resultSet.Len() > listSize {
				utils.PopMax(resultSet)
			}
		}
	}

	// K-NN is returned when search is performed
	if onlySearch {
		candidatesHeap := utils.NewMinHeapFromSliceDeep(candidatesVisited, len(candidatesVisited))

		knn := make([]utils.Element, 0, k)
		for i := 0; i < k && candidatesHeap.Len() > 0; i++ {
			knn = append(knn, heap.Pop(candidatesHeap).(utils.Element))
		}

		return knn, nil
	}

	return nil, candidatesVisited
}

func (mi *MemoryIndex) insert(v []float32, listSize int, distanceThreshold float32, currId, dataId uint32) error {
	mi.Lock() // TODO: optimize locking
	defer mi.Unlock()

	if mi.readOnly {
		return ErrReadOnlyIndex
	}

	if len(v) != int(mi.dim) {
		return ErrVectorDimensions
	}

	vertex := Vertex{
		id:     currId,
		objId:  dataId,
		vector: v,
	}

	var first bool
	mi.initialInsertion.Do(func() {
		if mi.s == 0 {
			mi.s = vertex.id
		}

		if mi.Size() == 0 { // first insertion
			mi.graph.addVertex(&vertex)
			first = true
		}
	})

	if first {
		return nil
	}

	mi.graph.addVertex(&vertex)

	_, visited := mi.search(v, 1, listSize, false)

	vertex.neighbors = mi.robustPrune(&vertex, visited, distanceThreshold)

	for _, n := range vertex.neighbors {
		neighbor := mi.graph.vertices[n]

		neighbor.neighbors = append(neighbor.neighbors, vertex.id)

		if len(neighbor.neighbors) > int(mi.maxDegree) {
			distances := make([]utils.Element, 0, mi.maxDegree)

			for _, nn := range neighbor.neighbors {
				nnVertex := mi.graph.vertices[nn]
				distances = append(distances, utils.Element{
					Id:       nn,
					Distance: mi.calculateDistance(nnVertex.vector, neighbor.vector),
				})
			}
			neighbor.neighbors = mi.robustPrune(neighbor, distances, distanceThreshold)
		}
	}

	return nil
}

func (mi *MemoryIndex) delete(id uint32) error {
	// TODO: delete consolidation?
	return nil
}

func (mi *MemoryIndex) robustPrune(v *Vertex, candidates []utils.Element, distanceThreshold float32) []uint32 {
	// TODO locking
	deletedCandidates := map[uint32]struct{}{v.id: {}} // excluding vertex v

	for _, n := range v.neighbors {
		nVertex := mi.graph.vertices[n]
		e := utils.Element{
			Id:       n,
			Distance: mi.calculateDistance(v.vector, nVertex.vector),
		}

		candidates = append(candidates, e)
	}

	candidatesHeap := utils.NewMinHeapFromSlice(candidates)
	newNeighbors := make([]uint32, 0, mi.maxDegree)

	for candidatesHeap.Len() != 0 {
		min := heap.Pop(candidatesHeap).(utils.Element)
		if _, ok := deletedCandidates[uint32(min.Id)]; ok {
			continue
		}

		newNeighbors = append(newNeighbors, uint32(min.Id))

		if len(v.neighbors) == int(mi.maxDegree) {
			break
		}

		for _, c := range candidatesHeap.Elements {
			minVertex := mi.graph.vertices[uint32(min.Id)]
			cVertex := mi.graph.vertices[uint32(c.Id)]

			if distanceThreshold*mi.calculateDistance(minVertex.vector, cVertex.vector) <= c.Distance {
				deletedCandidates[cVertex.id] = struct{}{}
			}
		}
	}

	return newNeighbors
}
func (mi *MemoryIndex) ReadOnly() {
	mi.readOnly = true
}

func (mi *MemoryIndex) Size() uint32 {
	return uint32(len(mi.graph.vertices))
}

func (mi *MemoryIndex) snapshot(path string) error {
	d, err := newDal(path)
	if err != nil {
		return err
	}

	err = d.writeIndex(mi)
	if err != nil {
		return err
	}

	mi.snapshotPath = path

	return nil
}

func (mi *MemoryIndex) serializeMetadata(buff []byte) int {
	var offset int

	binary.LittleEndian.PutUint32(buff[offset:], mi.dim)
	offset += 4

	binary.LittleEndian.PutUint32(buff[offset:], mi.maxDegree)
	offset += 4

	binary.LittleEndian.PutUint32(buff[offset:], mi.firstId)
	offset += 4

	binary.LittleEndian.PutUint32(buff[offset:], mi.Size())
	offset += 4

	binary.LittleEndian.PutUint32(buff[offset:], mi.s)
	offset += 4

	return offset
}

func (mi *MemoryIndex) deserializeMetadata(buff []byte) int {
	var offset int

	mi.dim = binary.LittleEndian.Uint32(buff[offset:])
	offset += 4

	mi.maxDegree = binary.LittleEndian.Uint32(buff[offset:])
	offset += 4

	mi.firstId = binary.LittleEndian.Uint32(buff[offset:])
	offset += 4

	mi.size = binary.LittleEndian.Uint32(buff[offset:])
	offset += 4

	mi.s = binary.LittleEndian.Uint32(buff[offset:])
	offset += 4

	return offset
}

func loadMemoryIndex(path string) (*MemoryIndex, error) {
	d, err := newDal(path)
	if err != nil {
		return nil, err
	}

	mi, err := d.readIndex()
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return mi, nil
}
