package embeddings

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	Text2VecHuggingFace = "text2vec-huggingface"
	BaseURL             = "https://api-inference.huggingface.co"
	Path                = "pipeline/feature-extraction"
)

type text2vecEmbedder struct {
	client *http.Client
	config *text2vecEmbedderConfig
}

type text2vecEmbedderConfig struct {
	ModelName string `json:"model_name"`
	ApiKey    string `json:"api_key"`
}

type EmbeddingRequest struct {
	Inputs []string `json:"inputs"`
}

type EmbeddingResponse [][]float32

func newText2vecEmbedder(cfg *text2vecEmbedderConfig) *text2vecEmbedder {
	e := text2vecEmbedder{
		client: http.DefaultClient,
		config: cfg,
	}

	return &e
}

func (e *text2vecEmbedder) Embed(ctx context.Context, inputs []string) ([][]float32, error) {
	body := EmbeddingRequest{
		Inputs: inputs,
	}

	b, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.GetURL(), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	res, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}

	var er EmbeddingResponse
	if err = json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}

	return er, nil
}

func (e *text2vecEmbedder) GetURL() string {
	return fmt.Sprintf("%s/%s/%s", BaseURL, Path, e.config.ModelName)
}
