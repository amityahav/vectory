package text2vec

const (
	Text2VecHuggingFace = "text2vec-huggingface"
	ModelName           = "msmarco-bert-base-dot-v5"
)

type Config struct {
	ApiKey string `json:"api_key"`
}

type EmbeddingRequest struct {
	Inputs []string `json:"inputs"`
}

type EmbeddingResponse [][]float32
