package main

import (
	"Vectory/db"
	"Vectory/entities/collection"
	"Vectory/entities/embeddings/hugging_face/text2vec"
	"Vectory/entities/index"
	"Vectory/entities/objstore"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	ctx := context.Background()
	vectory, err := db.Open("./data")
	if err != nil {
		log.Fatal(err)
	}

	c, err := collectionLoadOrCreate(ctx, vectory, &collection.Collection{
		Name:         "wine reviews",
		IndexType:    index.Hnsw,
		DataType:     collection.TextDataType,
		EmbedderType: text2vec.Text2VecHuggingFace,
		EmbedderConfig: text2vec.Config{
			ApiKey: os.Getenv("API_KEY"),
		},
		IndexParams: &index.DefaultHnswParams,
		Mappings: []string{"points",
			"title",
			"description",
			"taster_name",
			"taster_twitter_handle",
			"price",
			"designation",
			"variety",
			"region_1",
			"region_2",
			"province",
			"country",
			"winery"},
	})
	if err != nil {
		log.Fatal(err)
	}

	objects := loadReviews(500)
	err = c.InsertBatch(ctx, objects)
	if err != nil {
		log.Fatal(err)
	}

	res, err := c.SemanticSearch(ctx, &objstore.Object{
		Properties: map[string]interface{}{
			"question": "whats the best red wine in italy?",
		},
	}, 5)

	fmt.Println(res)

}

func collectionLoadOrCreate(ctx context.Context, vectory *db.DB, config *collection.Collection) (*db.Collection, error) {
	c, err := vectory.GetCollection(ctx, config.Name)
	if err != nil && errors.Is(err, db.ErrValidationFailed) {
		return vectory.CreateCollection(ctx, config)
	}

	return c, nil
}

func loadReviews(size int) []*objstore.Object {
	f, err := os.OpenFile("./examples/winemag-data-130k-v2.json", os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal(err)
	}

	b, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}

	var wineReviews []map[string]interface{}
	err = json.Unmarshal(b, &wineReviews)
	if err != nil {
		log.Fatal(err)
	}

	objects := make([]*objstore.Object, 0, size)
	for _, r := range wineReviews[0:size] {
		objects = append(objects, &objstore.Object{
			Properties: r,
		})
	}

	return objects
}
