package entities

type Collection struct {
	// name
	Name string `json:"name"`

	// index type
	IndexType string `json:"index_type"`

	// embedder
	Embedder string `json:"embedder"`

	// data type
	DataType string `json:"data_type"`

	// index params
	IndexParams interface{} `json:"index_params"`
}

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
