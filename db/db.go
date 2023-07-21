package db

import (
	"Vectory/db/metadata"
	"Vectory/gen/api/models"
	"Vectory/gen/ent"
	"context"
	"github.com/sirupsen/logrus"
	"sync"
)

type Config struct {
	FilesPath  string `yaml:"files_path"`
	ListenPort int    `yaml:"listen_port"`
}

type DB struct {
	sync.RWMutex
	config          *Config
	metadataManager *metadata.MetaManager
	collections     *sync.Map
	logger          *logrus.Logger
	wal             any
}

// Init initialises Vectory and restore collections and additional metadata if exists
func Init(cfg *Config) (*DB, error) {
	db := DB{
		config: cfg,
		logger: logrus.New(),
		wal:    nil,
	}

	db.logger.SetFormatter(&logrus.JSONFormatter{})

	mm, err := metadata.NewMetaManager(db.config.FilesPath)
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
func (db *DB) CreateCollection(ctx context.Context, cfg *models.Collection) (int, error) {
	collectionID, err := db.metadataManager.CreateCollection(ctx, cfg)
	if err != nil {
		return 0, err
	}

	if _, ok := db.collections.Load(cfg.Name); ok { // for safety
		return 0, ErrCollectionAlreadyExists
	}

	c, err := NewCollection(collectionID, cfg)
	if err != nil {
		return 0, err
	}

	db.collections.Store(c.name, c)

	return collectionID, nil
}

// DeleteCollection deletes collection both on disk and memory
func (db *DB) DeleteCollection(ctx context.Context, name string) error {
	err := db.metadataManager.DeleteCollection(ctx, name)
	if err != nil {
		return err
	}

	db.collections.Delete(name)

	return nil
}

// GetCollection returns the collection with name
func (db *DB) GetCollection(ctx context.Context, name string) (*ent.Collection, error) {
	return db.metadataManager.GetCollection(ctx, name)
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

		err = nc.restore(col, db.config.FilesPath)
		if err != nil {
			return err
		}

		db.collections.Store(nc.name, nc)
	}

	return nil
}
