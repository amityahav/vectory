package db

import (
	"Vectory/db/embeddings"
	"Vectory/entities/collection"
	"Vectory/entities/embeddings/hugging_face/text2vec"
	"Vectory/entities/index"
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

func TestCollection_Insert(t *testing.T) {
	ctx := context.Background()
	filesPath := "./tmp"

	db, err := Open(filesPath)
	require.NoError(t, err)
	defer os.RemoveAll(filesPath)

	{
		c, err := db.CreateCollection(ctx, &collection.Collection{
			Name:        "test_collection",
			IndexType:   index.Hnsw,
			DataType:    "text",
			IndexParams: index.DefaultHnswParams,
			Mappings:    []string{"title", "content"},
		})
		require.NoError(t, err)

		objects := make([]objstore.Object, 0, 100)
		ids := make([]uint64, 0, 100)
		t.Run("sequential insertion with vectors", func(t *testing.T) {
			for i := 0; i < 100; i++ {
				o := objstore.Object{
					Properties: map[string]interface{}{
						"title":   "Test",
						"content": "blah",
					},
					Vector: randomVector(128),
				}

				err = c.Insert(ctx, &o)
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

		t.Run("insert with no vector and no embedder", func(t *testing.T) {
			obj := objstore.Object{
				Properties: map[string]interface{}{
					"title":   "Test",
					"content": "blah",
				},
			}

			require.ErrorIs(t, c.Insert(ctx, &obj), ErrMissingVectorAndEmbedder)
		})

		t.Run("insert with invalid mappings", func(t *testing.T) {
			obj := objstore.Object{
				Properties: map[string]interface{}{
					"title": "Test",
					"x":     "blah",
				},
			}

			require.Error(t, c.Insert(ctx, &obj))
		})
	}

	{
		c, err := db.CreateCollection(ctx, &collection.Collection{
			Name:         "test_collection2",
			IndexType:    index.Hnsw,
			EmbedderType: embeddings.FakeEmbedder,
			DataType:     "text",
			IndexParams:  index.DefaultHnswParams,
			Mappings:     []string{"title", "content"},
		})
		require.NoError(t, err)

		t.Run("insert with vector and embedder", func(t *testing.T) {
			vec := randomVector(128)
			obj := objstore.Object{
				Properties: map[string]interface{}{
					"title":   "Test",
					"content": "blah",
				},
				Vector: vec,
			}

			err = c.Insert(ctx, &obj)
			require.NoError(t, err)
			require.Equal(t, vec, obj.Vector)

		})

		t.Run("insert without vector and embedder", func(t *testing.T) {
			obj := objstore.Object{
				Properties: map[string]interface{}{
					"title":   "Test",
					"content": "blah",
				},
			}

			err = c.Insert(ctx, &obj)
			require.NoError(t, err)
			require.NotNil(t, obj.Vector)

		})

		t.Run("InsertBatch with vectors and embedder", func(t *testing.T) {
			vec := randomVector(128)
			obj := objstore.Object{
				Properties: map[string]interface{}{
					"title":   "Test",
					"content": "blah",
				},
				Vector: vec,
			}

			obj2 := objstore.Object{
				Properties: map[string]interface{}{
					"title":   "Test",
					"content": "blah",
				},
			}

			err = c.InsertBatch(ctx, []*objstore.Object{&obj, &obj2})
			require.NoError(t, err)
			require.Equal(t, vec, obj.Vector)
			require.NotNil(t, obj2.Vector)
		})
	}

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
		Name:        "test_collection",
		IndexType:   index.Hnsw,
		DataType:    "text",
		IndexParams: index.DefaultHnswParams,
		Mappings:    []string{"title", "content"},
	})
	require.NoError(b, err)

	dim := 128
	objects := make([]*objstore.Object, 100000)
	for i := 0; i < len(objects); i++ {
		objects[i] = &objstore.Object{
			Id: 0,
			Properties: map[string]interface{}{
				"title":   "Test",
				"content": "blah",
			},
			Vector: randomVector(dim),
		}
	}

	//objects = loadObjects("./core/index/hnsw/siftsmall/siftsmall_base.fvecs")

	start := time.Now()
	err = c.InsertBatch(ctx, objects)
	require.NoError(b, err)
	end := time.Since(start)

	fmt.Printf("batch insertion took: %s\n", end)
}

func BenchmarkCollection_InsertBatch_WithEmbedder(b *testing.B) {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath("./profile")).Stop()
	b.ResetTimer()
	b.ReportAllocs()

	filesPath := "./tmp"

	ctx := context.Background()
	db, err := Open(filesPath)
	require.NoError(b, err)

	defer os.RemoveAll(filesPath)

	c, err := db.CreateCollection(ctx, &collection.Collection{
		Name:         "test_collection",
		IndexType:    index.Hnsw,
		EmbedderType: text2vec.Text2VecHuggingFace,
		EmbedderConfig: text2vec.Config{
			ApiKey: os.Getenv("api_key"),
		},
		DataType:    "text",
		IndexParams: index.DefaultHnswParams,
		Mappings:    []string{"title", "content"},
	})
	require.NoError(b, err)

	objects := make([]*objstore.Object, 1)
	for i := 0; i < len(objects); i++ {
		objects[i] = &objstore.Object{
			Properties: map[string]interface{}{
				"title":   "Test",
				"content": "blah",
			}}
	}

	//objects := loadObjects("./core/index/hnsw/siftsmall/siftsmall_base.fvecs")

	start := time.Now()
	err = c.InsertBatch(ctx, objects)
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
		Name:        "test_collection",
		IndexType:   index.Hnsw,
		DataType:    "text",
		IndexParams: index.DefaultHnswParams,
		Mappings:    []string{"title", "content"},
	})
	require.NoError(b, err)

	dim := 128
	objects := make([]*objstore.Object, 100000)
	for i := 0; i < len(objects); i++ {
		objects[i] = &objstore.Object{
			Id: 0,
			Properties: map[string]interface{}{
				"title":   "Test",
				"content": "blah",
			}, Vector: randomVector(dim),
		}
	}

	//objects = loadObjects("./core/index/hnsw/siftsmall/siftsmall_base.fvecs")

	start := time.Now()
	err = c.InsertBatch2(ctx, objects)
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
			Properties: map[string]interface{}{
				"title":   "Test",
				"content": "blah",
			},
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
