package db

import (
	"Vectory/entities"
	"Vectory/entities/collection"
	"Vectory/entities/objstore"
	"context"
	"fmt"
	"github.com/pkg/profile"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"testing"
	"time"
)

func TestCollection(t *testing.T) {
	filesPath := "./tmp"

	ctx := context.Background()
	db, err := Open(filesPath)
	require.NoError(t, err)
	defer os.RemoveAll(filesPath)

	c, err := db.CreateCollection(ctx, &collection.Collection{
		Name:      "test_collection",
		IndexType: "hnsw",
		DataType:  "text",
		IndexParams: collection.HnswParams{
			M:              64,
			MMax:           128,
			EfConstruction: 100,
			Ef:             100,
			Heuristic:      true,
			DistanceType:   entities.Euclidean,
		},
	})
	require.NoError(t, err)

	objects := make([]objstore.Object, 0, 1000)
	ids := make([]uint64, 0, 1000)
	t.Run("insertion with vectors", func(t *testing.T) {
		for i := 0; i < 1000; i++ {
			o := objstore.Object{
				Data:   fmt.Sprintf("%d", i),
				Vector: randomVector(128),
			}

			err = c.Insert(&o)
			require.NoError(t, err)

			o.Id = uint64(i)
			objects = append(objects, o)
			ids = append(ids, o.Id)
		}
	})

	t.Run("get inserted objects", func(t *testing.T) {
		objs, err := c.Get(ids)
		require.NoError(t, err)

		for i, o := range objects {
			require.Equal(t, o, objs[i])
		}
	})
}

func BenchmarkCollection_Insert(b *testing.B) {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath("./profile")).Stop()
	b.ResetTimer()
	b.ReportAllocs()

	filesPath := "./tmp"

	ctx := context.Background()
	db, err := Open(filesPath)
	require.NoError(b, err)

	defer os.RemoveAll(filesPath)

	c, err := db.CreateCollection(ctx, &collection.Collection{
		Name:      "test_collection",
		IndexType: "hnsw",
		DataType:  "text",
		IndexParams: collection.HnswParams{
			M:              64,
			MMax:           128,
			EfConstruction: 100,
			Ef:             100,
			Heuristic:      true,
			DistanceType:   entities.Euclidean,
		},
	})
	require.NoError(b, err)

	dim := 128
	objects := make([]*objstore.Object, 10000)
	for i := 0; i < len(objects); i++ {
		objects[i] = &objstore.Object{
			Id:     0,
			Data:   "Hello world",
			Vector: randomVector(dim),
		}
	}

	start := time.Now()
	err = c.InsertBatch(objects)
	require.NoError(b, err)
	end := time.Since(start)

	fmt.Printf("batch insertion took: %s\n", end)
}

func randomVector(dim int) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = rand.Float32()
	}

	return vec
}
