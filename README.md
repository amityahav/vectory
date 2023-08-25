# Vectory
![](docs/logo.png )

### What is Vectory
Vectory is an embeddable vector database built for fun. the storage engine is based on the [HNSW index](https://arxiv.org/abs/1603.09320).
it was inspired by [Weaviate](https://github.com/weaviate/weaviate).

### Status 
Vectory is not production ready and there's still work to do to make it stable.
I've worked on this project alone and now im looking for people to contribute, so if you are passionate about 
vector search, storage engines and Go, you can check out the open issues, submit a PR and i will be happy to review it.

### Design overview
![design](./docs/imgs/design.jpg)

the database object is composed of multiple components:

1. `Metadata manager` - is responsible for all collections' metadata such as name, index/embedder parameters and documents mappings in a persisted manner. 
2. `API` - currently there is support for REST API for creating/deleting collections when deploying Vectory on the cloud.
3. `Collection`:
   1. `Vector Index` - in-memory index for all the objects vectors.
   2. `Object store` - on-disk KV store for storing all objects.
   3. `Embbeder` - optional component which take as input a list of objects and returns their embeddings.


### How to use

```go
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
```