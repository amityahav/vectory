package disk_ann

import (
	"Vectory/db/core/indexes"
	"Vectory/db/core/indexes/utils"
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
	memoryIndexSizeLimit uint32
}

func NewDiskAnn() *DiskAnn {
	da := DiskAnn{
		deletedObjIds:        &sync.Map{},
		roIndexes:            nil,
		currId:               0,
		memoryIndexSizeLimit: 0,
	}

	da.rwIndex = newMemoryIndex(da.deletedObjIds, 0)

	return &da
}

// Load loads the size from disk and recovers in-memory indexes in case of a crash
func Load() {
}

func (da *DiskAnn) Insert(vector []float32, objId uint32) error {
	da.Lock()
	if da.rwIndex.Size() == da.memoryIndexSizeLimit {
		da.rwIndex.ReadOnly()
		go da.rwIndex.Snapshot()
		da.roIndexes = append(da.roIndexes, da.rwIndex)
		da.rwIndex = newMemoryIndex(da.deletedObjIds, da.currId)
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

func (da *DiskAnn) Search(q []float32, k int) []utils.Element {
	rwResults, _ := da.rwIndex.Search(q, k, da.listSize, true)

	for _, roIndex := range da.roIndexes {
		roResults, _ := roIndex.Search(q, k, da.listSize, true)
		rwResults = append(rwResults, roResults...)
	}

	ltResults, _ := da.ltIndex.Search(q, k, da.listSize, true)

	rwResults = append(rwResults, ltResults...)

	results := utils.NewMinMaxHeapFromSlice(rwResults)

	for results.Len() > k {
		utils.PopMax(results)
	}

	return results.Elements
}
