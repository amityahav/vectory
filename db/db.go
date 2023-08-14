package db

import (
	"Vectory/db/metadata"
	"Vectory/entities/collection"
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
)

type DB struct {
	mu              sync.RWMutex
	metadataManager *metadata.MetaManager
	collections     *sync.Map
	logger          *logrus.Logger
	filesPath       string
}

// Open initialises Vectory and init collections and additional metadata if exists.
func Open(filesPath string) (*DB, error) {
	db := DB{
		logger:    logrus.New(),
		filesPath: filesPath,
	}

	err := db.init()
	if err != nil {
		return nil, errors.Wrap(err, "failed opening vectory")
	}

	return &db, nil
}

// CreateCollection creates a new collection in the database and cache it in memory.
func (db *DB) CreateCollection(ctx context.Context, cfg *collection.Collection) (*Collection, error) {
	err := collection.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}

	if _, ok := db.collections.Load(cfg.Name); ok { // for safety
		return nil, ErrCollectionAlreadyExists
	}

	collectionID, err := db.metadataManager.CreateCollection(ctx, cfg)
	if err != nil {
		return nil, err
	}

	c, err := newCollection(collectionID, cfg, db.filesPath)
	if err != nil {
		return nil, err
	}

	db.collections.Store(c.name, c)

	return c, nil
}

// DeleteCollection deletes collection both on disk and memory.
func (db *DB) DeleteCollection(ctx context.Context, name string) error {
	// TODO: handle case where deleting and another user has ref to the collection trying accessing removed files
	if _, ok := db.collections.Load(name); !ok {
		return fmt.Errorf("%w: %s", ErrValidationFailed, ErrCollectionDoesntExist)
	}

	err := db.metadataManager.DeleteCollection(ctx, name)
	if err != nil {
		return err
	}

	db.collections.Delete(name)

	return nil
}

// GetCollection returns the collection with name.
func (db *DB) GetCollection(ctx context.Context, name string) (*Collection, error) {
	c, ok := db.collections.Load(name)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrValidationFailed, ErrCollectionDoesntExist)
	}

	return c.(*Collection), nil
}

// init collections and metadata to memory.
func (db *DB) init() error {
	ctx := context.Background()

	db.logger.SetFormatter(&logrus.JSONFormatter{})

	mm, err := metadata.NewMetaManager(db.filesPath)
	if err != nil {
		return err
	}

	db.metadataManager = mm

	db.collections = &sync.Map{}

	cols, err := db.metadataManager.GetCollections(ctx)
	if err != nil {
		return err
	}

	for _, col := range cols {
		c, err := newCollection(col.ID, &collection.Collection{
			Name:           col.Name,
			IndexType:      col.IndexType,
			EmbedderType:   col.EmbedderType,
			IndexParams:    col.IndexParams,
			EmbedderConfig: col.EmbedderConfig,
			DataType:       col.DataType,
		}, db.filesPath)

		if err != nil {
			return err
		}

		db.collections.Store(c.name, c)
	}

	return nil
}
