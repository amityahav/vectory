package db

import (
	"Vectory/db/core/index"
	"Vectory/db/core/index/hnsw"
	"Vectory/db/core/objstore"
	"Vectory/db/embeddings"
	"Vectory/entities"
	"Vectory/entities/collection"
	"encoding/json"
	"fmt"
	"github.com/alitto/pond"
	"runtime"
)

var _ CRUD = &Collection{}

type Collection struct {
	id          int
	name        string
	dataType    string
	objStore    *objstore.ObjectStore
	vectorIndex index.VectorIndex
	idCounter   *IdCounter
	logger      any
	embedder    embeddings.Embedder
	wp          *pond.WorkerPool
	filesPath   string
	config      collection.Collection
}

func newCollection(id int, cfg *collection.Collection, filesPath string) (*Collection, error) {
	c := Collection{
		id:       id,
		name:     cfg.Name,
		dataType: cfg.DataType,
		embedder: nil,
		wp:       pond.New(runtime.NumCPU(), 1000), // TODO: should be configurable
		config:   *cfg,
	}

	c.filesPath = fmt.Sprintf("%s/%s", filesPath, c.name)

	os, err := objstore.NewObjectStore(c.filesPath)
	if err != nil {
		return nil, err
	}

	c.objStore = os

	counter, err := newIdCounter(c.filesPath)
	if err != nil {
		return nil, err
	}

	c.idCounter = counter

	switch cfg.IndexType {
	case entities.Hnsw:
		var params collection.HnswParams

		b, _ := json.Marshal(cfg.IndexParams) // validated in wrapper functions
		_ = json.Unmarshal(b, &params)

		idx, err := hnsw.NewHnsw(params, c.filesPath, os)
		if err != nil {
			return nil, err
		}

		c.vectorIndex = idx
	default:
		return nil, ErrUnknownIndexType
	}

	return &c, nil
}

func (c *Collection) GetConfig() collection.Collection {
	return c.config
}
