package db

import (
	"Vectory/db/metadata"
	"Vectory/db/validators"
	"Vectory/entities"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"sync"
)

type DB struct {
	sync.RWMutex
	metadataManager *metadata.MetaManager
	collections     *sync.Map
	logger          *logrus.Logger
	filesPath       string
	wal             any
}

// Open initialises Vectory and restore collections and additional metadata if exists
func Open(filesPath string) (*DB, error) {
	db := DB{
		logger:    logrus.New(),
		filesPath: filesPath,
		wal:       nil,
	}

	db.logger.SetFormatter(&logrus.JSONFormatter{})

	mm, err := metadata.NewMetaManager(filesPath)
	if err != nil {
		return nil, err
	}

	db.metadataManager = mm

	err = db.restore()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

// CreateCollection creates a new collection in the database and cache it in memory
func (db *DB) CreateCollection(ctx context.Context, cfg *entities.Collection) (*Collection, error) {
	err := validators.ValidateCollection(cfg)
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

	c, err := NewCollection(collectionID, cfg, db.filesPath)
	if err != nil {
		return nil, err
	}

	db.collections.Store(c.name, c)

	return c, nil
}

// DeleteCollection deletes collection both on disk and memory
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

// GetCollection returns the collection with name
func (db *DB) GetCollection(ctx context.Context, name string) (*Collection, error) {
	c, ok := db.collections.Load(name)
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrValidationFailed, ErrCollectionDoesntExist)
	}

	return c.(*Collection), nil
}

// restore collections and metadata to memory
func (db *DB) restore() error {
	ctx := context.Background()

	db.collections = &sync.Map{}

	cols, err := db.metadataManager.GetCollections(ctx)
	if err != nil {
		return err
	}

	for _, col := range cols {
		nc := Collection{}

		err = nc.restore(col) // TODO: can be parallelized
		if err != nil {
			return err
		}

		db.collections.Store(nc.name, nc)
	}

	return nil
}
