package db

import (
	"Vectory/db/core/indexes"
	"Vectory/db/core/indexes/disk_ann"
)

var _ CRUD = &Collection{}

type collectionConfig struct {
	Name      string `json:"name"`
	IndexType string `json:"index_type"`
}

type Collection struct {
	id       any
	name     string
	objStore any //obj_store.db
	index    indexes.VectorIndex
	logger   any
	embedder any
	wal      any
}

func NewCollection(cfg *collectionConfig) (*Collection, error) {
	// TODO: persist
	c := Collection{
		id:       nil,
		name:     cfg.Name,
		objStore: nil,
		logger:   nil,
		embedder: nil,
		wal:      nil,
	}

	switch cfg.IndexType {
	case "disk_ann":
		c.index = disk_ann.NewDiskAnn()
	default:
		return nil, ErrUnknownIndexType
	}

	return &c, nil
}

func (c *Collection) Insert(obj any) {
	// TODO: store obj, create embedding and store objId and vector in index
}

func (c *Collection) InsertWithVector(obj any, vector []float32) {

}

func (c *Collection) Delete(objId uint32) {
	// TODO: delete in objStore and in index
}

func (c *Collection) Update(objId uint32) {
	// TODO: delete both in index and objStore and create again
}

func (c *Collection) Get(objIds []uint32) {
	// TODO: get objects with objIds from objStore
}

func (c *Collection) SemanticSearch(obj any, k int) {
	// TODO: create embeddings from obj and get K-NN and then retrieve object ids from objStore
}

func (c *Collection) DeleteMyself() error {
	return nil
}
