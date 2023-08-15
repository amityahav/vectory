package db

import "errors"

var (
	ErrUnknownIndexType         = errors.New("unknown index type")
	ErrCollectionAlreadyExists  = errors.New("collection with the same name already exists")
	ErrCollectionDoesntExist    = errors.New("collection does not exist")
	ErrValidationFailed         = errors.New("validation failed")
	ErrMissingVectorAndEmbedder = errors.New("can't insert an object without vector when there's no embedder")
)
