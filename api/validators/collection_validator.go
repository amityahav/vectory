package validators

import (
	"Vectory/entities"
	"Vectory/gen/api/models"
	"errors"
)

func ValidateCollection(cfg *models.Collection) error {
	if cfg.Name == "" {
		return ErrCollectionNameEmpty
	}

	var err error

	switch cfg.IndexType {
	case entities.Hnsw:
		err = validateHnswParams(cfg.IndexParams)
	case entities.DiskAnn:
	default:
		return ErrIndexTypeUnsupported
	}

	if err != nil {
		return err
	}

	switch cfg.DataType {
	case entities.Text:
	default:
		return ErrDataTypeUnsupported
	}

	return nil
}

var (
	ErrCollectionNameEmpty  = errors.New("collection name field is empty")
	ErrIndexTypeUnsupported = errors.New("index_type inserted is not supported")
	ErrDataTypeUnsupported  = errors.New("data_type inserted is not supported")
)
