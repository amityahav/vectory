package index

import "Vectory/entities/distance"

const (
	Hnsw    = "hnsw"
	DiskAnn = "disk_ann"
)

type HnswParams struct {
	// Number of established connections
	M int `json:"m"`

	// Maximum number of connections for each element per layer
	MMax int `json:"m_max"`

	// size of the dynamic candidate list
	EfConstruction int `json:"ef_construction"`

	Ef int `json:"ef"`

	Heuristic bool `json:"heuristic"`

	DistanceType string `json:"distance_type"`
}

var DefaultHnswParams = HnswParams{
	M:              64,
	MMax:           128,
	EfConstruction: 100,
	Ef:             100,
	Heuristic:      true,
	DistanceType:   distance.Euclidean,
}
