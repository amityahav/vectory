package disk_ann

import (
	"Vectory/db/core/indexes"
	"sync"
)

var _ indexes.VectorIndex = &DiskAnn{}

type DiskAnn struct {
	sync.RWMutex

	rwIndex   *MemoryIndex
	roIndexes []*MemoryIndex
	ltIndex   *DiskIndex

	deletedObjIds *sync.Map

	listSize          int
	distanceThreshold int

	currId               uint32
	memoryIndexSizeLimit int
}

func NewDiskAnn() *DiskAnn {
	da := DiskAnn{
		deletedObjIds:        &sync.Map{},
		roIndexes:            nil,
		currId:               0,
		memoryIndexSizeLimit: 0,
	}

	da.rwIndex = newMemoryIndex(da.deletedObjIds)

	return &da
}

// Load loads the index from disk and recovers in-memory indexes in case of a crash
func Load() {
}

func (da *DiskAnn) Insert(vector []float32, objId uint32) error {
	da.Lock()
	if da.rwIndex.graph.size() == da.memoryIndexSizeLimit {
		da.rwIndex.ReadOnly()
		go da.rwIndex.Snapshot()
		da.roIndexes = append(da.roIndexes, da.rwIndex)
		da.rwIndex = newMemoryIndex(da.deletedObjIds)
	}
	da.Unlock()

	da.Lock()
	currId := da.currId
	da.currId += 1
	da.Unlock()

	err := da.rwIndex.Insert(vector, da.listSize, da.distanceThreshold, currId, objId)
	if err != nil {
		return err
	}

	return nil
}

func (da *DiskAnn) Delete(objId uint32) bool {
	da.deletedObjIds.Store(objId, struct{}{})

	return true
}

// Search all indexes with onlySearch = true, maintain a MinMax heap to keep only K-NN from all indexes and return them
func (da *DiskAnn) Search(vector []float32, k int) {

}
