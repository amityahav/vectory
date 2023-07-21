package db

import (
	"Vectory/db/core/indexes"
	"Vectory/db/core/indexes/hnsw"
	"Vectory/entities"
	"Vectory/gen/api/models"
	"encoding/json"
)

var _ CRUD = &Collection{}

type Collection struct {
	id       int
	name     string
	dataType string
	objStore any //obj_store.db
	index    indexes.VectorIndex
	logger   any
	embedder any
	wal      any
}

func NewCollection(id int, cfg *models.Collection) (*Collection, error) {
	c := Collection{
		id:       id,
		name:     cfg.Name,
		dataType: cfg.DataType,
		embedder: cfg.Embedder,
		objStore: nil,
		logger:   nil,
		wal:      nil,
	}

	switch cfg.IndexType { // cfg.IndexParams has already been validated
	case entities.Hnsw:
		var params entities.HnswParams

		b, _ := json.Marshal(cfg.IndexParams)
		_ = json.Unmarshal(b, &params)

		c.index = hnsw.NewHnsw(params)
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
