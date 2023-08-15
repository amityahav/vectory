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

func (c *Collection) GetConfig() collection.Collection {
	return c.config
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
