package disk_ann

import "sync"

// Api - DiskANN API
type Api interface {
	// Insert a new vector to the RW-Index
	Insert(vector []float32) error

	// Delete Lazy deletion of vertex v
	Delete(id uint32)

	// Search for K-NN by querying Long-term index, RW-index and all RO-indexes, aggregating and filtering results
	Search(vector []float32, k int, l int)
}

type DiskAnn struct {
	sync.RWMutex
	rwIndex              *MemoryIndex
	roIndexes            []*MemoryIndex
	deletedIds           *sync.Map
	memoryIndexSizeLimit int
	ltIndex              string
}

func NewDiskAnn() *DiskAnn {
	da := DiskAnn{
		deletedIds:           &sync.Map{},
		roIndexes:            nil,
		memoryIndexSizeLimit: 0,
		ltIndex:              "",
	}

	da.rwIndex = newMemoryIndex(0, da.deletedIds)

	return &da
}

// Load loads the index from disk and recovers in-memory indexes in case of a crash
func Load() {
}

func (da *DiskAnn) Insert(vector []float32) error {
	da.Lock()
	if da.rwIndex.graph.size() == da.memoryIndexSizeLimit {
		da.rwIndex.ReadOnly()
		go da.rwIndex.Snapshot()
		da.roIndexes = append(da.roIndexes, da.rwIndex)
		da.rwIndex = newMemoryIndex(da.rwIndex.currId, da.deletedIds)
	}
	da.Unlock()

	err := da.rwIndex.Insert(vector)
	if err != nil {
		return err
	}

	return nil
}

func (da *DiskAnn) Delete(id uint32) bool {
	da.deletedIds.Store(id, struct{}{})

	return true
}
