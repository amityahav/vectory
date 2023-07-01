package db

import "sync"

type DB struct {
	sync.RWMutex
	collections *sync.Map
	logger      any
	config      any
	wal         any
}

func NewDB(path string) (*DB, error) {
	return &DB{}, nil
}

func (db *DB) CreateCollection(cfg *collectionConfig) error {
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

func (db *DB) DeleteCollection(name string) error {
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
