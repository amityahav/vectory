package db

import (
	"Vectory/db/core/indexes"
	"Vectory/db/core/indexes/hnsw"
	"Vectory/db/core/objstore"
	"Vectory/entities"
	"Vectory/entities/collection"
	objstoreentities "Vectory/entities/objstore"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
)

var _ CRUD = &Collection{}

type Collection struct {
	id          int
	name        string
	dataType    string
	objStore    *objstore.ObjectStore
	vectorIndex indexes.VectorIndex
	idCounter   *IdCounter
	logger      any
	embedder    any
	filesPath   string
	wal         any
	config      collection.Collection
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

		b, _ := json.Marshal(cfg.IndexParams) // validated in wrapper functions
		_ = json.Unmarshal(b, &params)

		index, err := hnsw.NewHnsw(params, c.filesPath)
		if err != nil {
			return nil, err
		}

		c.vectorIndex = index
	default:
		return nil, ErrUnknownIndexType
	}

	return &c, nil
}

func (c *Collection) Insert(obj *objstoreentities.Object) error {
	// TODO: validate obj data type is the same as collection's
	if obj.Vector == nil {
		// TODO: create embedding and store in obj
	}

	id, err := c.idCounter.FetchAndInc()
	if err != nil {
		return err
	}

	obj.Id = id

	if err = c.objStore.Put(obj); err != nil {
		return err
	}

	if err = c.vectorIndex.Insert(obj.Vector, obj.Id); err != nil {
		return err
	}

	// TODO: flush hnsw wal

	return nil
}

func (c *Collection) Update(obj *objstoreentities.Object) error {
	// TODO: handle race conditions
	return nil
}

func (c *Collection) Delete(objId uint64) error {
	// TODO: delete in objStore and in index
	_, found, err := c.objStore.Get(objId)
	if err != nil {
		return errors.Wrapf(err, "failed getting %d from object store", objId)
	}

	if !found {
		// nothing to do
		return nil
	}

	err = c.objStore.Delete(objId)
	if err != nil {
		return errors.Wrapf(err, "failed deleting %d from object store", objId)
	}

	err = c.vectorIndex.Delete(objId)
	if err != nil {
		return errors.Wrapf(err, "failed deleting %d from vector index", objId)
	}

	// TODO: flush hnsw wal
	return nil
}

func (c *Collection) Get(objIds []uint64) ([]objstoreentities.Object, error) {
	objects := make([]objstoreentities.Object, 0, len(objIds))
	for _, id := range objIds {
		obj, found, err := c.objStore.Get(id)
		if err != nil {
			return nil, errors.Wrapf(err, "failed getting %d from object store", id)
		}

		if !found {
			continue
		}

		objects = append(objects, *obj)
	}

	return objects, nil
}

func (c *Collection) SemanticSearch(obj any, k int) {
	// TODO: create embeddings from obj and get K-NN and then retrieve object ids from objStore
}

func (c *Collection) GetConfig() collection.Collection {
	return c.config
}
