package collection

const (
	TextDataType = "text"
)

type Collection struct {

	// name
	Name string `json:"name,omitempty"`

	// index type
	IndexType string `json:"index_type,omitempty"`

	// embedder type
	EmbedderType string `json:"embedder_type,omitempty"`

	// data type
	DataType string `json:"data_type,omitempty"`

	// index params
	IndexParams interface{} `json:"index_params,omitempty"`

	// embedder config
	EmbedderConfig interface{} `json:"embedder_config,omitempty"`

	// schema
	Schema interface{} `json:"schema,omitempty"`
}
