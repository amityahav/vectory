package embeddings

import (
	"Vectory/entities/embeddings/hugging_face"
	"Vectory/entities/embeddings/hugging_face/text2vec"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type text2vecEmbedder struct {
	client *http.Client
	config *text2vec.Config
}

func NewText2vecEmbedder(cfg *text2vec.Config) *text2vecEmbedder {
	e := text2vecEmbedder{
		client: http.DefaultClient,
		config: cfg,
	}

	return &e
}

func (e *text2vecEmbedder) Embed(ctx context.Context, inputs []string) ([][]float32, error) {
	body := text2vec.EmbeddingRequest{
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

	var er text2vec.EmbeddingResponse
	if err = json.NewDecoder(res.Body).Decode(&er); err != nil {
		return nil, err
	}

	return er, nil
}

func (e *text2vecEmbedder) GetURL() string {
	return fmt.Sprintf("%s/%s/%s", hugging_face.BaseURL, hugging_face.Path, text2vec.ModelName)
}
