package db

import (
	"Vectory/db/core/indexes"
	"Vectory/db/core/indexes/hnsw"
	"Vectory/db/core/objstore"
	"Vectory/entities"
	"Vectory/entities/collection"
	objstoreentities "Vectory/entities/objstore"
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
	config    collection.Collection
}

func newCollection(id int, cfg *collection.Collection, filesPath string) (*Collection, error) {
	c := Collection{
		id:       id,
		name:     cfg.Name,
		dataType: cfg.DataType,
		embedder: nil,
		logger:   nil,
		wal:      nil,
		config:   *cfg,
	}

	c.filesPath = fmt.Sprintf("%s/%s", filesPath, c.name)

	os, err := objstore.NewObjectStore(c.filesPath)
	if err != nil {
		return nil, err
	}

	c.objStore = os

	counter, err := NewIdCounter(c.filesPath)
	if err != nil {
		return nil, err
	}

	c.idCounter = counter

	switch cfg.IndexType {
	case entities.Hnsw:
		var params collection.HnswParams

		b, _ := json.Marshal(cfg.IndexParams)
		_ = json.Unmarshal(b, &params)

		c.index = hnsw.NewHnsw(params)
	default:
		return nil, ErrUnknownIndexType
	}

	return &c, nil
}

func (c *Collection) Insert(obj *objstoreentities.Object) error {
	// TODO: validate obj data type is the same as collection's

	return nil
}

func (c *Collection) Update(obj *objstoreentities.Object) error {
	// TODO: handle race conditions
	return nil
}

func (c *Collection) InsertWithVector(obj *objstoreentities.Object, vector []float32) error {
	id, err := c.idCounter.FetchAndInc()
	if err != nil {
		return err
	}

	obj.Id = id

	if err = c.objStore.Put(obj); err != nil {
		return err
	}

	if err = c.index.Insert(vector, obj.Id); err != nil {
		return err
	}

	// TODO: flush WALS

	return nil
}

func (c *Collection) Delete(objId uint64) {
	// TODO: delete in objStore and in index
}

func (c *Collection) Get(objIds []uint64) ([]objstoreentities.Object, error) {
	// TODO: get objects with objIds from objStore
	objects := make([]objstoreentities.Object, 0, len(objIds))
	for _, id := range objIds {
		obj, err := c.objStore.Get(id)
		if err != nil {
			return nil, err // TODO: continue instead?
		}

		objects = append(objects, *obj)
	}

	return objects, nil
}

func (c *Collection) SemanticSearch(obj any, k int) {
	// TODO: create embeddings from obj and get K-NN and then retrieve object ids from objStore
}

func (c *Collection) restore(col *ent.Collection) error {
	c.id = col.ID
	c.name = col.Name
	c.dataType = col.DataType
	c.config = collection.Collection{
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

func (c *Collection) GetConfig() collection.Collection {
	return c.config
}
