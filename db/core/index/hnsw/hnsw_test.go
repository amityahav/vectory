package hnsw

import (
	"Vectory/db/core/objstore"
	"Vectory/entities/distance"
	"Vectory/entities/index"
	objstoreentities "Vectory/entities/objstore"
	"fmt"
	"github.com/pkg/profile"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func TestHnsw(t *testing.T) {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath("./profile")).Stop()

	filesPath := "./tmp"
	defer os.RemoveAll(filesPath)

	hnsw, _ := newHnsw(filesPath)
	dim := 128

	for i := 0; i < 10000; i++ {
		if i%1000 == 0 {
			log.Printf("%d", i)
		}

		err := hnsw.Insert(randomVector(dim), uint64(i))

		if err != nil {
			t.Error(err)
		}
	}

	_ = hnsw.Search(randomVector(dim), 10)
}

func BenchmarkHnsw_Insert(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	defer profile.Start(profile.CPUProfile, profile.ProfilePath("./profile")).Stop()
	filesPath := "./tmp"

	hnsw, _ := newHnsw(filesPath)
	size := 10000
	dim := 128
	ch := make(chan job, 10000)
	vec := randomVector(dim)
	for i := 0; i < size; i++ {
		ch <- job{
			id:     uint64(i),
			vector: vec,
		}

		if i%1000 == 0 {
			log.Printf("inserted %d", i)
		}

	}
	close(ch)

	start := time.Now()
	var wg sync.WaitGroup
	insertInParallel(ch, hnsw, &wg)
	end := time.Since(start)

	fmt.Printf("insertion took: %s\n", end)
	os.RemoveAll(filesPath)
	//_ = hnsw.Search(randomVector(dim), 10)
}

func TestRestoreFromDisk(t *testing.T) {
	filesPath := "../tmp"
	defer os.RemoveAll(filesPath)

	store, err := objstore.NewObjectStore(filesPath)
	require.NoError(t, err)

	dim := 128
	h, err := NewHnsw(index.HnswParams{
		M:              64,
		MMax:           128,
		EfConstruction: 100,
		Ef:             100,
		Heuristic:      true,
		DistanceType:   distance.Euclidean,
	}, filesPath, store)

	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		vec := randomVector(dim)
		o := objstoreentities.Object{
			Id: uint64(i),
			Properties: map[string]interface{}{
				"title": "test",
			},
			Vector: vec,
		}

		require.NoError(t, h.Insert(vec, uint64(i)))
		require.NoError(t, store.Put(&o))
	}

	require.NoError(t, h.Flush())

	// create new hnsw and restore from wal
	start := time.Now()
	hRestored, err := NewHnsw(index.HnswParams{
		M:              64,
		MMax:           128,
		EfConstruction: 100,
		Ef:             100,
		Heuristic:      true,
		DistanceType:   distance.Euclidean,
	}, filesPath, store)
	require.NoError(t, err)

	end := time.Since(start)

	fmt.Printf("reloading from disk took: %s\n", end)

	// compare old hnsw with restored hnsw
	require.Equal(t, h.entrypointID, hRestored.entrypointID)
	require.Equal(t, h.currentMaxLayer, hRestored.currentMaxLayer)
	require.Equal(t, h.nodes, hRestored.nodes)
}

func TestSift(t *testing.T) {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath("./profile")).Stop()

	// Loading vectors
	ch := make(chan job, 10000)
	loadSiftBaseVectors("./siftsmall/siftsmall_base.fvecs", ch)
	//loadRandomVectors(ch)

	// Building index
	filesPath := "./tmp"
	start := time.Now()
	hnsw := buildIndexParallel(ch, filesPath)
	defer os.RemoveAll(filesPath)

	duration := time.Since(start)

	log.Printf("Building index took: %v", duration)

	// Searching index
	queryVectors := loadSiftQueryVectors("./siftsmall/siftsmall_query.fvecs")
	truthNeighbors := loadSiftTruthVectors("./siftsmall/siftsmall_groundtruth.ivecs")

	k := 10

	var avgRecall float32
	for i, q := range queryVectors {
		var match float32
		ann := hnsw.Search(q, k)
		truth := truthNeighbors[i]

		for _, n := range ann {
			for _, tn := range truth {
				if n.Id == tn {
					match++
					break
				}
			}
		}

		avgRecall += match

		log.Printf("recall for V#%d: %f", i, match/float32(len(ann)))
	}

	avgRecall /= float32(k * len(queryVectors))

	log.Printf("avg recall: %f", avgRecall)
}
