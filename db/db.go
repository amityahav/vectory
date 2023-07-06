package db

import (
	"Vectory/db/metadata"
	"context"
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

func NewDB(cfg *Config) (*DB, error) {
	db := DB{
		config:      cfg,
		collections: nil,
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

	return &db, nil
}

func (db *DB) CreateCollection(ctx context.Context, cfg *collectionConfig) error {
	err := db.metadataManager.CreateCollection(ctx)
	if err != nil {
		return err
	}

	// TODO: persist
	if _, ok := db.collections.Load(cfg.Name); ok {
		return ErrCollectionAlreadyExists
	}

	c, err := NewCollection(cfg)
	if err != nil {
		return err
	}

	db.collections.Store(cfg.Name, c)

	return nil
}

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
