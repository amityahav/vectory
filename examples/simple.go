package examples

import (
	"Vectory/db"
	"Vectory/entities/collection"
	"Vectory/entities/distance"
	"Vectory/entities/embeddings/hugging_face/text2vec"
	"Vectory/entities/index"
	"Vectory/entities/objstore"
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	vectory, err := db.Open("./data")
	if err != nil {
		log.Fatal(err)
	}

	c, err := vectory.CreateCollection(ctx, &collection.Collection{
		Name:      "Movie Reviews",
		IndexType: index.Hnsw,
		IndexParams: index.HnswParams{
			M:              0,
			MMax:           0,
			EfConstruction: 0,
			Ef:             0,
			Heuristic:      false,
			DistanceType:   distance.Euclidean,
		},
		EmbedderType: text2vec.Text2VecHuggingFace,
		EmbedderConfig: text2vec.Config{
			ApiKey: os.Getenv("API_KEY"),
		},
		DataType: collection.TextDataType,
	})
	if err != nil {
		log.Fatal(err)
	}

	objects := make([]*objstore.Object, 0, 100)
	for i := 0; i < 100; i++ {
		o := objstore.Object{
			Data: fmt.Sprintf("review-%d", i),
		}

		objects = append(objects, &o)
	}

	if err = c.InsertBatch(ctx, objects); err != nil {
		log.Fatal(err)
	}

	ann, err := c.SemanticSearch(ctx, &objstore.Object{
		Data: "barbie",
	}, 10)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ann)
}
