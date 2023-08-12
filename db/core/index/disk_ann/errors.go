package disk_ann

import "errors"

var (
	ErrReadOnlyIndex    = errors.New("can't insert to a read-only index")
	ErrVectorDimensions = errors.New("vector was not inserted because of wrong dimensions")
)
