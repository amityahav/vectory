package examples

import (
	"Vectory/db"
	"Vectory/entities/collection"
	"Vectory/entities/index"
	"Vectory/entities/objstore"
	"context"
	"log"
)

func main() {
	ctx := context.Background()
	vectory, err := db.Open("./data")
	if err != nil {
		log.Fatal(err)
	}

	collection, err := vectory.CreateCollection(ctx, &collection.Collection{
		Name:      "wikipedia",
		IndexType: index.Hnsw,
		DataType:  collection.TextDataType,
		IndexParams: &index.HnswParams{
			M:              0,
			MMax:           0,
			EfConstruction: 0,
			Ef:             0,
			Heuristic:      false,
			DistanceType:   "",
		},
		Mappings: []string{"id", "title", "text", "url", "wiki_id", "views", "paragraph_id", "langs"},
	})
	if err != nil {
		log.Fatal(err)
	}

	objects := loadDataset()

	err = collection.InsertBatch(ctx, objects)
	if err != nil {
		log.Fatal(err)
	}

}

func loadDataset() []*objstore.Object {
	return nil
}
