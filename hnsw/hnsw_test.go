package hnsw

import (
	"Vectory/hnsw/distance"
	bufio2 "bufio"
	"encoding/binary"
	"io"
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
		efConstruction:        128,
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
			Id:     int64(i),
			Vector: randomVector(dim),
		})

		if err != nil {
			t.Error(err)
		}
	}

	_ = hnsw.Search(&Vertex{
		Id:     -1,
		Vector: randomVector(dim),
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
	// Loading vectors
	start := time.Now()
	vertices := loadSiftBaseVectors("./siftsmall/siftsmall_base.fvecs")
	end := time.Since(start)
	log.Printf("Loading vectors took: %s", end.String())

	// Building index
	start = time.Now()
	hnsw := buildIndexSequential(vertices)
	end = time.Since(start)
	log.Printf("Inserting vectors took: %s", end.String())

	// Searching index
	queryVectors := loadSiftQueryVectors("./siftsmall/siftsmall_query.fvecs")
	truthNeighbors := loadSiftTruthVectors("./siftsmall/siftsmall_groundtruth.ivecs")

	var avgRecall float32
	for i, q := range queryVectors {
		var match float32
		ann := hnsw.Search(&Vertex{Vector: q}, 100, 64)
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

		log.Printf("recall for V#%d: %f", i, match/float32(len(ann)))
	}

	avgRecall /= float32(len(queryVectors) * len(queryVectors))

	log.Printf("avg recall: %f", avgRecall)
}

func buildIndexSequential(vertices []*Vertex) *Hnsw {
	hnsw := newHnsw()

	for i, v := range vertices {
		err := hnsw.Insert(v)
		if err != nil {
			log.Printf("failed inserting V#%d. err: %s", v.Id, err.Error())
		}

		log.Printf("Inserted %d vertices", i)
	}

	return hnsw
}

func buildIndexParallel(insertionChannel chan *Vertex) *Hnsw {
	hnsw := newHnsw()
	var wg sync.WaitGroup

	for i := 0; i < 1; i++ {
		wg.Add(1)
		go func(i int) {
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
		}(i)
	}

	wg.Wait()

	return hnsw
}

func loadSiftBaseVectors(path string) []*Vertex {
	return sequential(path)
}

func loadSiftQueryVectors(path string) [][]float32 {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	buf := bufio2.NewReader(f)
	b := make([]byte, 4)

	vectors := make([][]float32, 100)

	for i := 0; i < 100; i++ {
		dim, err := readUint32(buf, b)
		if err != nil {
			log.Fatal(err)
		}

		vectors[i] = make([]float32, dim)

		for j := 0; j < int(dim); j++ {
			vectors[i][j], err = readFloat32(buf, b)
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
	buf := bufio2.NewReader(f)
	b := make([]byte, 4)

	vectors := make([][]int64, 100)

	for i := 0; i < 100; i++ {
		dim, err := readUint32(buf, b)
		if err != nil {
			log.Fatal(err)
		}

		vectors[i] = make([]int64, dim)

		for j := 0; j < int(dim); j++ {
			fl32, err := readUint32(buf, b)
			if err != nil {
				log.Fatal(err)
			}

			vectors[i][j] = int64(fl32)
		}
	}

	return vectors
}

func sequential(path string) []*Vertex {
	s := make([]*Vertex, 0, 10000)

	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	buf := bufio2.NewReader(f)
	b := make([]byte, 4)

	defer f.Close()

	for i := 0; i < 10000; i++ {
		dim, err := readUint32(buf, b)
		if err != nil {
			log.Fatal(err)
		}

		v := Vertex{
			Id:     int64(i),
			Vector: make([]float32, dim),
		}

		for j := 0; j < int(dim); j++ {
			v.Vector[j], err = readFloat32(buf, b)
			if err != nil {
				log.Fatal(err)
			}

		}

		s = append(s, &v)
	}

	return s
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

			defer f.Close()

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
					Id:     int64(chunkNum*chunkSize + j),
					Vector: make([]float32, dim),
				}

				for l := 0; l < int(dim); l++ {
					v.Vector[l], err = readFloat32(buf, b)
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
		log.Printf("Loaded V#%d", v.Id)
	}
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
