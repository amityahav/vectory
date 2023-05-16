package hnsw

import (
	"Vectory/hnsw/distance"
	"encoding/binary"
	"log"
	"math"
	"math/rand"
	"os"
	"sync"
	"testing"
)

func newHnsw() *Hnsw {
	cfg := hnswConfig{
		m:                     32,
		mMax:                  32,
		efConstruction:        400,
		heuristic:             true,
		distanceType:          distance.DotProduct,
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

	go loadSiftBaseVectors("./siftsmall/siftsmall_base.fvecs", insertionChannel)
	hnsw := buildIndex(insertionChannel)
	print(hnsw)

}

func buildIndex(insertionChannel chan *Vertex) *Hnsw {
	hnsw := newHnsw()
	var wg sync.WaitGroup

	for i := 0; i < 1; i++ { // TODO: change to runtime.NumCPU() and handle concurrent insertions
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
			log.Printf("loaded %d vectors", i)
		}

		insertionChannel <- &v
	}
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
