package handlers

import "errors"

var (
	ErrCollectionNameEmpty  = errors.New("collection name field is empty")
	ErrIndexTypeUnsupported = errors.New("index_type inserted is not supported")
	ErrDataTypeUnsupported  = errors.New("data_type inserted is not supported")
)
