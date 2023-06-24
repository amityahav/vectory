package disk_ann

import "sync"

// Api - DiskANN API
type Api interface {
	// Insert a new vector to the RW-Index
	Insert(vector []float32) error

	// Delete Lazy deletion of vertex v
	Delete(v string)

	// Search for K-NN by querying Long-term index, RW-index and all RO-indexes, aggregating and filtering results
	Search(vector []float32, k int, l int)
}

type DiskAnn struct {
	sync.RWMutex
	rwIndex              *MemoryIndex
	roIndexes            []*MemoryIndex
	memoryIndexSizeLimit int
	ltIndex              string
}

func NewDiskAnn() {
}

// Load loads the index from disk and recovers in-memory indexes in case of a crash
func Load() {
}

func (da *DiskAnn) Insert(vector []float32) error {
	da.Lock()
	if da.rwIndex.graph.size() == da.memoryIndexSizeLimit {
		go da.rwIndex.Snapshot()
		da.rwIndex = newMemoryIndex()
	}
	da.Unlock()

	err := da.rwIndex.Insert(vector)
	if err != nil {
		return err
	}

}
