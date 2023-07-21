package db

import "errors"

var (
	ErrUnknownIndexType        = errors.New("unknown index type")
	ErrCollectionAlreadyExists = errors.New("collection with the same name already exists")
)
