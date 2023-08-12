package hnsw

import (
	"Vectory/db/core/objstore"
	"Vectory/entities/collection"
	objstoreentities "Vectory/entities/objstore"
	bufio2 "bufio"
	"fmt"
	"github.com/pkg/profile"
	"github.com/stretchr/testify/require"
	"log"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

func BenchmarkHnsw_Insert(b *testing.B) {
	b.ResetTimer()
	b.ReportAllocs()
	//defer profile.Start(profile.CPUProfile, profile.ProfilePath("./profile")).Stop()
	filesPath := "./tmp"

	hnsw, _ := newHnsw(filesPath)
	size := 1000
	dim := 128
	ch := make(chan job, size)
	for i := 0; i < size; i++ {
		j := job{
			id:     uint64(i),
			vector: randomVector(dim),
		}

		ch <- j
	}
	close(ch)

	start := time.Now()

	wg := sync.WaitGroup{}
	go insertInParallel(ch, hnsw, &wg)
	wg.Wait()

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
	h, err := NewHnsw(collection.HnswParams{
		M:              64,
		MMax:           128,
		EfConstruction: 100,
		Ef:             100,
		Heuristic:      true,
		DistanceType:   "dot_product",
	}, filesPath, store)

	require.NoError(t, err)

	for i := 0; i < 2000; i++ {
		vec := randomVector(dim)
		o := objstoreentities.Object{
			Id:     uint64(i),
			Data:   "",
			Vector: vec,
		}

		require.NoError(t, h.Insert(vec, uint64(i)))
		require.NoError(t, store.Put(&o))
	}

	// create new hnsw and restore from wal
	start := time.Now()
	hRestored, err := NewHnsw(collection.HnswParams{
		M:              64,
		MMax:           128,
		EfConstruction: 100,
		Ef:             100,
		Heuristic:      true,
		DistanceType:   "dot_product",
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
	ch := make(chan job)
	go loadSiftBaseVectors("./siftsmall/siftsmall_base.fvecs", ch)

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

func sequential(path string, ch chan job) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

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

		ch <- job{
			id:     uint64(i),
			vector: vector,
		}
	}

	close(ch)
}

func concurrent() {
	insertionChannel := make(chan *Vertex)
	cpus := runtime.NumCPU()
	vectorsCount := 10000
	chunkSize := vectorsCount / cpus
	remainder := vectorsCount % cpus

	chunkSizes := make([]int, cpus)
	for i := 0; i < len(chunkSizes); i++ {
		chunkSizes[i] = chunkSize
	}

	chunkSizes[len(chunkSizes)-1] += remainder

	var wg sync.WaitGroup
	wg.Add(cpus)
	for i := 0; i < cpus; i++ {
		go func(chunkNum int) {
			defer wg.Done()

			f, err := os.Open("hnsw/siftsmall/siftsmall_base.fvecs")
			if err != nil {
				log.Fatal(err)
			}

			buf := bufio2.NewReader(f)
			b := make([]byte, 4)

			offset := chunkNum * chunkSize * (128*4 + 4)
			_, err = f.Seek(int64(offset), 0)
			if err != nil {
				log.Fatal(err)
			}

			for j := 0; j < chunkSizes[chunkNum]; j++ {
				dim, err := readUint32(buf, b)
				if err != nil {
					log.Fatal(err)
				}

				v := Vertex{
					id:     uint64(chunkNum*chunkSize + j),
					vector: make([]float32, dim),
				}

				for l := 0; l < int(dim); l++ {
					v.vector[l], err = readFloat32(buf, b)
					if err != nil {
						log.Fatal(err)
					}
				}

				insertionChannel <- &v
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(insertionChannel)
	}()

	for v := range insertionChannel {
		log.Printf("Loaded V#%d", v.id)
	}
}
