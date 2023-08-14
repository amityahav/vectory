package embeddings

import (
	"context"
	"encoding/json"
	"fmt"
)

type EmbedderConfig struct {
	Name   string
	Config interface{}
}

type Embedder interface {
	Embed(ctx context.Context, inputs []string) ([][]float32, error)
}

func NewEmbedder(config *EmbedderConfig) (Embedder, error) {
	switch config.Name {
	case Text2VecHuggingFace:
		var cfg text2vecEmbedderConfig

		b, _ := json.Marshal(cfg)
		_ = json.Unmarshal(b, &cfg) // TODO: validate

		return newText2vecEmbedder(&cfg), nil
	default:
		return nil, fmt.Errorf("unsupported embedder with name: %s", config.Name)
	}
}
