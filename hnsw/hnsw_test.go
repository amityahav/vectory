package hnsw

import (
	"Vectory/hnsw/distance"
	"encoding/binary"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"testing"
	"time"
)

func newHnsw() *Hnsw {
	cfg := hnswConfig{
		m:                     32,
		mMax:                  32,
		efConstruction:        400,
		heuristic:             true,
		distanceType:          distance.Euclidean,
		extendCandidates:      true,
		keepPrunedConnections: true,
	}

	cfg.mMax0 = 2 * cfg.mMax
	cfg.mL = 1 / math.Log(float64(cfg.m))

	return NewHnsw(cfg)
}

func TestHnsw(t *testing.T) {
	hnsw := newHnsw()
	dim := 128

	for i := 0; i < 1000; i++ {
		err := hnsw.Insert(&Vertex{
			id:     int64(i),
			vector: randomVector(dim),
		})

		if err != nil {
			t.Error(err)
		}
	}

	_ = hnsw.KnnSearch(&Vertex{
		id:     -1,
		vector: randomVector(dim),
	}, 10, 400)

}

func randomVector(dim int) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = rand.Float32()
	}

	return vec
}

func TestSift(t *testing.T) {
	insertionChannel := make(chan *Vertex)

	// Building index
	start := time.Now()
	go loadSiftBaseVectors("./siftsmall/siftsmall_base.fvecs", insertionChannel)
	hnsw := buildIndex(insertionChannel)
	end := time.Since(start)

	log.Printf("Building index took: %s", end.String())

	// Searching index
	queryVectors := loadSiftQueryVectors("./siftsmall/siftsmall_query.fvecs")
	truthNeighbors := loadSiftTruthVectors("./siftsmall/siftsmall_groundtruth.ivecs")

	var avgRecall float32
	for i, q := range queryVectors {
		var match float32
		ann := hnsw.KnnSearch(&Vertex{vector: q}, 100, 100)
		truth := truthNeighbors[i]

		for _, n := range ann {
			for _, tn := range truth {
				if n == tn {
					match++
					break
				}
			}
		}

		avgRecall += match
	}

	avgRecall /= float32(len(queryVectors))

	log.Printf("avg recall: %f", avgRecall)
}

func buildIndex(insertionChannel chan *Vertex) *Hnsw {
	hnsw := newHnsw()
	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				v, ok := <-insertionChannel
				if !ok {
					break
				}

				err := hnsw.Insert(v)
				if err != nil {
					log.Fatal(err)
					return
				}
			}
		}()
	}

	wg.Wait()

	return hnsw
}

func loadSiftBaseVectors(path string, insertionChannel chan *Vertex) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	dim, err := readUint32(f)
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 10000; i++ {
		v := Vertex{
			id:     int64(i),
			vector: make([]float32, dim),
		}

		for j := 0; j < int(dim); j++ {
			v.vector[j], err = readFloat32(f)
			if err != nil {
				log.Fatal(err)
			}

		}

		if (i+1)%1000 == 0 {
			log.Printf("loaded %d vectors", i+1)
		}

		insertionChannel <- &v
	}

	close(insertionChannel)
}

func loadSiftQueryVectors(path string) [][]float32 {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	dim, err := readUint32(f)
	if err != nil {
		log.Fatal(err)
	}

	vectors := make([][]float32, 100)

	for i := 0; i < 100; i++ {
		vectors[i] = make([]float32, dim)

		for j := 0; j < int(dim); j++ {
			vectors[i][j], err = readFloat32(f)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	return vectors
}

func loadSiftTruthVectors(path string) [][]int64 {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	dim, err := readUint32(f)
	if err != nil {
		log.Fatal(err)
	}

	vectors := make([][]int64, 100)

	for i := 0; i < 100; i++ {
		vectors[i] = make([]int64, dim)

		for j := 0; j < int(dim); j++ {
			fl32, err := readUint32(f)
			if err != nil {
				log.Fatal(err)
			}

			vectors[i][j] = int64(fl32)
		}
	}

	return vectors
}

func readUint32(f *os.File) (uint32, error) {
	b := make([]byte, 4)
	_, err := f.Read(b)
	if err != nil {
		return 0, err
	}

	return binary.LittleEndian.Uint32(b), nil
}

func readFloat32(f *os.File) (float32, error) {
	b := make([]byte, 4)
	_, err := f.Read(b)
	if err != nil {
		return 0, err
	}

	return math.Float32frombits(binary.LittleEndian.Uint32(b)), nil
}
