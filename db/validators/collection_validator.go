package validators

import (
	"Vectory/entities"
	"encoding/json"
	"errors"
)

func ValidateCollection(cfg *entities.Collection) error {
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

func validateHnswParams(params interface{}) error {
	var hnswParams entities.HnswParams

	b, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &hnswParams)
	if err != nil {
		return err
	}

	if hnswParams.M <= 0 {
		return errors.New("m must be greater than zero")
	}

	if hnswParams.MMax <= 0 {
		return errors.New("m_max must be greater than zero")
	}

	if hnswParams.Ef <= 0 {
		return errors.New("ef must be greater than zero")
	}

	if hnswParams.EfConstruction <= 0 {
		return errors.New("ef_construction must be greater than zero")
	}

	switch hnswParams.DistanceType {
	case entities.DotProduct:
	case entities.Euclidean:
	default:
		return errors.New("unsupported distance_type")
	}

	return nil
}

var (
	ErrCollectionNameEmpty  = errors.New("collection name field is empty")
	ErrIndexTypeUnsupported = errors.New("index_type inserted is not supported")
	ErrDataTypeUnsupported  = errors.New("data_type inserted is not supported")
)
