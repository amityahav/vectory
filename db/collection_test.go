package db

import (
	"Vectory/entities"
	"Vectory/entities/collection"
	"Vectory/entities/objstore"
	bufio2 "bufio"
	"context"
	"encoding/binary"
	"fmt"
	"github.com/pkg/profile"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"math"
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

func loadObjects(path string) []*objstore.Object {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	objects := make([]*objstore.Object, 10000)

	buf := bufio2.NewReader(f)
	b := make([]byte, 4)
	for i := 0; i < 10000; i++ {
		dim, err := readUint32(buf, b)
		if err != nil {
			log.Fatal(err)
		}

		vector := make([]float32, dim)

		for j := 0; j < int(dim); j++ {
			vector[j], err = readFloat32(buf, b)
			if err != nil {
				log.Fatal(err)
			}
		}

		objects[i] = &objstore.Object{
			Data:   "",
			Vector: vector,
		}
	}

	return objects
}

func readUint32(f io.Reader, b []byte) (uint32, error) {
	_, err := f.Read(b)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(b), nil
}

func readFloat32(f io.Reader, b []byte) (float32, error) {
	_, err := f.Read(b)
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(binary.LittleEndian.Uint32(b)), nil
}

func BenchmarkCollection_InsertBatch(b *testing.B) {
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
			EfConstruction: 400,
			Ef:             100,
			Heuristic:      true,
			DistanceType:   entities.Euclidean,
		},
	})
	require.NoError(b, err)

	//dim := 128
	//objects := make([]*objstore.Object, 10000)
	//for i := 0; i < len(objects); i++ {
	//	objects[i] = &objstore.Object{
	//		Id:     0,
	//		Data:   "Hello world",
	//		Vector: randomVector(dim),
	//	}
	//}

	objects := loadObjects("./core/index/hnsw/siftsmall/siftsmall_base.fvecs")

	start := time.Now()
	err = c.InsertBatch(objects)
	require.NoError(b, err)
	end := time.Since(start)

	fmt.Printf("batch insertion took: %s\n", end)
}

func BenchmarkCollection_InsertBatch2(b *testing.B) {
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
			EfConstruction: 400,
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

	objects = loadObjects("./core/index/hnsw/siftsmall/siftsmall_base.fvecs")

	ch := make(chan *objstore.Object, 10000)
	for i := 0; i < 10000; i++ {
		ch <- objects[i]
	}
	close(ch)

	start := time.Now()
	err = c.InsertBatch2(ch)
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
