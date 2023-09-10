package db

import (
	"Vectory/db/core/index"
	"Vectory/db/core/index/hnsw"
	"Vectory/db/core/objstore"
	"Vectory/db/embeddings"
	"Vectory/entities/collection"
	"Vectory/entities/embeddings/hugging_face/text2vec"
	indexentities "Vectory/entities/index"
	objstoreentities "Vectory/entities/objstore"
	"context"
	"encoding/json"
	"fmt"
	"github.com/alitto/pond"
	"runtime"
	"sync"
)

var _ CRUD = &Collection{}

type Collection struct {
	mu          sync.RWMutex
	id          int
	name        string
	dataType    string
	stores      *objstore.Stores
	vectorIndex index.VectorIndex
	idCounter   *IdCounter
	logger      any
	embedder    embeddings.Embedder
	wp          *pond.WorkerPool
	filesPath   string
	config      collection.Collection
	closed      bool
}

func newCollection(id int, cfg *collection.Collection, filesPath string) (*Collection, error) {
	c := Collection{
		id:       id,
		name:     cfg.Name,
		dataType: cfg.DataType,
		wp:       pond.New(runtime.NumCPU(), 1000), // TODO: should be configurable
		config:   *cfg,
	}

	c.filesPath = fmt.Sprintf("%s/%s", filesPath, c.name)

	os, err := objstore.NewStores(c.filesPath)
	if err != nil {
		return nil, err
	}

	c.stores = os

	counter, err := newIdCounter(c.filesPath)
	if err != nil {
		return nil, err
	}

	c.idCounter = counter

	switch cfg.IndexType {
	case indexentities.Hnsw:
		var params indexentities.HnswParams

		b, _ := json.Marshal(cfg.IndexParams) // validated in wrapper function
		_ = json.Unmarshal(b, &params)

		idx, err := hnsw.NewHnsw(params, c.filesPath, os)
		if err != nil {
			return nil, err
		}

		c.vectorIndex = idx
	default:
		return nil, ErrUnknownIndexType
	}

	switch cfg.EmbedderType {
	case text2vec.Text2VecHuggingFace:
		var config text2vec.Config

		b, _ := json.Marshal(cfg.EmbedderConfig) // validated in wrapper function
		_ = json.Unmarshal(b, &config)

		c.embedder = embeddings.NewText2vecEmbedder(&config)
	case embeddings.FakeEmbedder: // for test purposes
		c.embedder = embeddings.NewFakeEmbedder()
	}

	return &c, nil
}

// GetConfig returns collection's configurations.
func (c *Collection) GetConfig() (*collection.Collection, error) {
	if c.closed {
		return nil, ErrCollectionClosed
	}

	return &c.config, nil
}

// GetSize returns the number of objects in the collection.
func (c *Collection) GetSize() (int, error) {
	if c.closed {
		return 0, ErrCollectionClosed
	}

	return c.stores.Size(), nil
}

// Close closes the collection.
func (c *Collection) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.idCounter.close(); err != nil {
		return err
	}

	if err := c.stores.Close(); err != nil {
		return err
	}

	c.closed = true

	return nil
}

// IsClosed indicates whether the collection is closed or not.
func (c *Collection) IsClosed() bool {
	return c.closed
}

// TODO: currently checking naively the mapping keys but in future check types as well
func (c *Collection) validateObjectsMappings(objs []*objstoreentities.Object) error {
	for i, obj := range objs {
		if len(obj.Properties) > len(c.config.Mappings) {
			return fmt.Errorf("length mismatch of object number %d properties and collection's mappings", i)
		}

		for _, m := range c.config.Mappings {
			if _, ok := obj.Properties[m]; !ok {
				return fmt.Errorf("object number %d does not have property %s", i, m)
			}
		}
	}

	return nil
}

func (c *Collection) embedObjectsIfNeeded(ctx context.Context, objs []*objstoreentities.Object) error {
	if c.embedder == nil {
		for _, o := range objs {
			if o.Vector == nil {
				return ErrMissingVectorAndEmbedder
			}

		}
	} else {
		filtered := make([]*objstoreentities.Object, 0, len(objs))
		for _, o := range objs {
			if o.Vector != nil {
				continue
			}

			filtered = append(filtered, o)
		}

		if err := c.embedder.Embed(ctx, filtered); err != nil {
			return err
		}
	}

	return nil
}
