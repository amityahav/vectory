package main

import (
	"Vectory/db"
	"Vectory/entities/collection"
	"Vectory/entities/embeddings/hugging_face/text2vec"
	"Vectory/entities/index"
	"Vectory/entities/objstore"
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	vectory, _ := db.Open("./data")

	c, _ := vectory.CreateCollection(ctx, &collection.Collection{
		Name:           "movie reviews",
		IndexType:      index.Hnsw,
		EmbedderType:   text2vec.Text2VecHuggingFace,
		DataType:       collection.TextDataType,
		IndexParams:    index.DefaultHnswParams,
		EmbedderConfig: text2vec.Config{ApiKey: os.Getenv("API_KEY")},
		Mappings: []string{
			"title",
			"review",
		},
	})

	// insert a single object with vector.
	_ = c.Insert(ctx, &objstore.Object{
		Properties: map[string]interface{}{
			"title":  "movie-1",
			"review": "great movie..",
		},
		Vector: []float32{1, 2, 3, 4},
	})

	// insert a single object and let the embedder create vector for you.
	_ = c.Insert(ctx, &objstore.Object{
		Properties: map[string]interface{}{
			"title":  "movie-2",
			"review": "bad movie..",
		},
	})

	// insert a batch of objects and let the embedder create vectors for you.
	_ = c.InsertBatch(ctx, []*objstore.Object{{
		Properties: map[string]interface{}{
			"title":  "movie-3",
			"review": "bad movie..",
		},
	}, {
		Properties: map[string]interface{}{
			"title":  "movie-4",
			"review": "bad movie.."},
	}, {
		Properties: map[string]interface{}{
			"title":  "movie-5",
			"review": "bad movie.."},
	}})

	// perform a semantic search over the inserted objects.
	res, _ := c.SemanticSearch(ctx, &objstore.Object{
		Properties: map[string]interface{}{
			"question": "whats the best movie to watch?",
		},
	}, 5)

	fmt.Println(res)
}
