package db

import (
	"Vectory/db/core/indexes"
	"Vectory/db/core/indexes/hnsw"
	"Vectory/db/core/objstore"
	"Vectory/entities"
	"Vectory/gen/ent"
	"encoding/json"
	"fmt"
)

var _ CRUD = &Collection{}

type Collection struct {
	id        int
	name      string
	dataType  string
	objStore  *objstore.ObjectStore //obj_store.db
	index     indexes.VectorIndex
	idCounter *IdCounter
	logger    any
	embedder  any
	filesPath string
	wal       any
	config    entities.Collection
}

func NewCollection(id int, cfg *entities.Collection, filesPath string) (*Collection, error) {
	c := Collection{
		id:       id,
		name:     cfg.Name,
		dataType: cfg.DataType,
		embedder: nil,
		objStore: nil,
		logger:   nil,
		wal:      nil,
		config:   *cfg,
	}

	c.filesPath = fmt.Sprintf("%s/%s", filesPath, c.name)

	switch cfg.IndexType {
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

func (c *Collection) Insert(obj *objstore.Object) error {
	// TODO: validate obj data type is the same as collection's
	id, err := c.idCounter.FetchAndInc()
	if err != nil {
		return err // TODO better error wrappings
	}

	// insert into object store
	c.objStore.Put(id, obj)

	// embed

	// insert into vector index

	return nil
}

func (c *Collection) InsertWithVector(obj *objstore.Object, vector []float32) error {
	return nil
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

func (c *Collection) restore(col *ent.Collection) error {
	c.id = col.ID
	c.name = col.Name
	c.dataType = col.DataType
	c.config = entities.Collection{
		Name:        col.Name,
		IndexType:   col.IndexType,
		Embedder:    col.Embedder,
		DataType:    col.DataType,
		IndexParams: col.IndexParams,
	}

	// restore index

	// restore objStore

	// restore embedder

	// restore wal

	return nil
}

func (c *Collection) GetConfig() entities.Collection {
	return c.config
}
