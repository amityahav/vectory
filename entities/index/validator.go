package index

import (
	"Vectory/entities/distance"
	"encoding/json"
	"errors"
)

func ValidateHnswParams(params interface{}) error {
	var hnswParams HnswParams

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
	case distance.DotProduct:
	case distance.Euclidean:
	default:
		return errors.New("unsupported distance_type")
	}

	return nil
}
