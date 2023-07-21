package validators

import (
	"Vectory/entities"
	"encoding/json"
	"errors"
)

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
