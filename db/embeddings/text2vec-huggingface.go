package embeddings

import (
	"Vectory/entities/embeddings/hugging_face"
	"Vectory/entities/embeddings/hugging_face/text2vec"
	"Vectory/entities/objstore"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
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

func (e *text2vecEmbedder) Embed(ctx context.Context, objects []*objstore.Object) error {
	ch := make(chan []float32, len(objects))
	var wg sync.WaitGroup
	for _, o := range objects {
		input := o.FlatProperties()
		body := text2vec.EmbeddingRequest{
			Inputs: []string{input},
		}

		b, err := json.Marshal(body)
		if err != nil {
			return err
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.GetURL(), bytes.NewReader(b))
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.config.ApiKey))

		wg.Add(1)
		go func(req *http.Request) {
			defer wg.Done()

			res, err := e.client.Do(req)
			if err != nil {
				fmt.Printf("failed create embeddings")
			}

			if res.StatusCode != http.StatusOK {
				fmt.Printf("failed create embeddings")
			}

			er := make(text2vec.EmbeddingResponse, len(objects))
			if err = json.NewDecoder(res.Body).Decode(&er); err != nil {
				fmt.Printf("failed create embeddings")
			}

			ch <- er[0]
		}(req)

	}

	wg.Wait()
	print(len(ch))

	return nil
}

//func (e *text2vecEmbedder) Embed(ctx context.Context, objects []*objstore.Object) error {
//	inputs := make([]string, 0, len(objects))
//	for _, o := range objects {
//		inputs = append(inputs, o.FlatProperties())
//	}
//
//	body := text2vec.EmbeddingRequest{
//		Inputs: inputs,
//	}
//
//	b, err := json.Marshal(body)
//	if err != nil {
//		return err
//	}
//
//	req, err := http.NewRequestWithContext(ctx, http.MethodPost, e.GetURL(), bytes.NewReader(b))
//	if err != nil {
//		return err
//	}
//
//	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", e.config.ApiKey))
//
//	res, err := e.client.Do(req)
//	if err != nil {
//		return err
//	}
//
//	if res.StatusCode != http.StatusOK {
//		return fmt.Errorf("failed create embeddings")
//	}
//
//	er := make(text2vec.EmbeddingResponse, len(objects))
//	if err = json.NewDecoder(res.Body).Decode(&er); err != nil {
//		return err
//	}
//
//	for i := 0; i < len(objects); i++ {
//		objects[i].Vector = er[i]
//	}
//
//	return nil
//}

func (e *text2vecEmbedder) GetURL() string {
	return fmt.Sprintf("%s/%s/%s", hugging_face.BaseURL, hugging_face.Path, text2vec.ModelName)
}
