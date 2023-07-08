package db

import (
	"Vectory/db/metadata"
	"Vectory/gen/api/models"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
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

// Init initialises Vectory and load collections and additional metadata if exists
func Init(cfg *Config) (*DB, error) {
	db := DB{
		config:      cfg,
		collections: &sync.Map{},
		logger:      logrus.New(),
		wal:         nil,
	}

	db.logger.SetFormatter(&logrus.JSONFormatter{})

	if stat, err := os.Stat(cfg.FilesPath); err != nil || !stat.IsDir() {
		return nil, ErrPathNotDirectory
	}

	mm, err := metadata.NewMetaManager(db.config.FilesPath)
	if err != nil {
		return nil, err
	}

	db.metadataManager = mm

	err = db.load()
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

	// creating collection data's dir
	err = os.Mkdir(fmt.Sprintf("%s/%s", db.config.FilesPath, cfg.Name), 0750)
	if err != nil && !os.IsExist(err) {
		return 0, err
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
	// TODO: persist
	c, ok := db.collections.Load(name)

	if !ok {
		return ErrCollectionDoesntExist
	}

	err := c.(*Collection).DeleteMyself()
	if err != nil {
		return err
	}

	db.collections.Delete(name)

	return nil
}

// load collections and metadata to memory
func (db *DB) load() error {
	ctx := context.Background()

	_, err := db.metadataManager.GetCollections(ctx)
	if err != nil {
		return err
	}

	return nil
}
