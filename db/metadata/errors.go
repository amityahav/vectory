package metadata

import "errors"

var (
	ErrPathNotDirectory      = errors.New("the path provided is not a directory")
	ErrCollectionDoesntExist = errors.New("collection does not exist")
)
