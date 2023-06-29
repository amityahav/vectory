package disk_ann

import "sync"

// Api - DiskANN API
type Api interface {
	// Insert a new vector and its corresponding data to the RW-Index
	Insert(vector []float32, data any) error

	// Delete Lazy deletion of vertex v by its dataId
	Delete(dataId uint32)

	// Search for K-NN by querying Long-term index, RW-index and all RO-indexes, aggregating and filtering results
	Search(vector []float32, k int, l int)
}

type DiskAnn struct {
	sync.RWMutex

	rwIndex   *MemoryIndex
	roIndexes []*MemoryIndex
	ltIndex   *DiskIndex

	deletedDataIds *sync.Map

	currId               uint32
	memoryIndexSizeLimit int
}

func NewDiskAnn() *DiskAnn {
	da := DiskAnn{
		deletedDataIds:       &sync.Map{},
		roIndexes:            nil,
		currId:               0,
		memoryIndexSizeLimit: 0,
	}

	da.rwIndex = newMemoryIndex(da.deletedDataIds)

	return &da
}

// Load loads the index from disk and recovers in-memory indexes in case of a crash
func Load() {
}

func (da *DiskAnn) Insert(vector []float32, data any) error {
	da.Lock()
	if da.rwIndex.graph.size() == da.memoryIndexSizeLimit {
		da.rwIndex.ReadOnly()
		go da.rwIndex.Snapshot()
		da.roIndexes = append(da.roIndexes, da.rwIndex)
		da.rwIndex = newMemoryIndex(da.deletedDataIds)
	}
	da.Unlock()

	// TODO: insert data to some datastore and keep its identifier
	var dataId uint32

	da.Lock()
	currId := da.currId
	da.currId += 1
	da.Unlock()

	err := da.rwIndex.Insert(vector, currId, dataId)
	if err != nil {
		return err
	}

	return nil
}

func (da *DiskAnn) Delete(id uint32) bool {
	da.deletedDataIds.Store(id, struct{}{})

	// TODO: remove data from the datastore

	return true
}

// Search all indexes with searchApi = true, maintain a MinMax heap to keep only K-NN from all indexes and return them
func (da *DiskAnn) Search(vector []float32, k int, l int) {

}
