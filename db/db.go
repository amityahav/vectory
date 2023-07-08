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

func NewDB(cfg *Config) (*DB, error) {
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

	return &db, nil
}

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
