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
	collections     map[string]*Collection
	logger          *logrus.Logger
	filesPath       string
	closed          bool
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
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.closed {
		return nil, ErrDatabaseClosed
	}

	err := collection.Validate(cfg)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidationFailed, err)
	}

	if _, ok := db.collections[cfg.Name]; ok { // for safety
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

	db.collections[c.name] = c

	return c, nil
}

// DeleteCollection deletes collection both on disk and memory.
func (db *DB) DeleteCollection(ctx context.Context, name string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if db.closed {
		return ErrDatabaseClosed
	}

	if _, ok := db.collections[name]; !ok {
		return fmt.Errorf("%w: %s", ErrValidationFailed, ErrCollectionDoesntExist)
	}

	err := db.metadataManager.DeleteCollection(ctx, name)
	if err != nil {
		return err
	}

	delete(db.collections, name)

	return nil
}

// GetCollection returns the collection with name.
func (db *DB) GetCollection(_ context.Context, name string) (*Collection, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	if db.closed {
		return nil, ErrDatabaseClosed
	}

	c, ok := db.collections[name]
	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrValidationFailed, ErrCollectionDoesntExist)
	}

	return c, nil
}

// Close closes the database
func (db *DB) Close() error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if err := db.metadataManager.Close(); err != nil {
		return err
	}

	for _, c := range db.collections {
		if c.IsClosed() {
			continue
		}

		if err := c.Close(); err != nil {
			return err
		}
	}

	db.closed = true

	return nil
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

	db.collections = map[string]*Collection{}

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

		db.collections[c.name] = c
	}

	return nil
}
