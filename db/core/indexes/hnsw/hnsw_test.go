package hnsw

import (
	bufio2 "bufio"
	"encoding/binary"
	"github.com/pkg/profile"
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

type job struct {
	id     int64
	vector []float32
}

func newHnsw() *Hnsw {
	return nil
}

func TestHnsw(t *testing.T) {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath("./profile")).Stop()

	hnsw := newHnsw()
	dim := 128

	for i := 0; i < 10000; i++ {
		if i%1000 == 0 {
			log.Printf("%d", i)
		}

		err := hnsw.Insert(randomVector(dim), int64(i), 1)

		if err != nil {
			t.Error(err)
		}
	}

	_ = hnsw.Search(randomVector(dim), 10)

}

func randomVector(dim int) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = rand.Float32()
	}

	return vec
}

func TestSift(t *testing.T) {
	defer profile.Start(profile.MemProfile, profile.ProfilePath("./profile")).Stop()

	// Loading vectors
	ch := make(chan job)
	go loadSiftBaseVectors("./siftsmall/siftsmall_base.fvecs", ch)

	// Building index
	start := time.Now()
	hnsw := buildIndexParallel(ch)
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

//
//func buildIndexSequential(vertices []*Vertex) *Hnsw {
//	hnsw := newHnsw()
//
//	for i, v := range vertices {
//		err := hnsw.Insert(v)
//		if err != nil {
//			log.Printf("failed inserting V#%d. err: %s", v.id, err.Error())
//		}
//
//		if i%1000 == 0 {
//			log.Printf("Inserted %d vertices", i)
//		}
//	}
//
//	return hnsw
//}

func buildIndexParallel(insertionChannel chan job) *Hnsw {
	hnsw := newHnsw()
	var wg sync.WaitGroup

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				job, ok := <-insertionChannel
				if !ok {
					break
				}

				err := hnsw.Insert(job.vector, job.id, 1)
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

func loadSiftBaseVectors(path string, ch chan job) {
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
			id:     int64(i),
			vector: vector,
		}

		if i%1000 == 0 {
			log.Printf("inserted %d", i)
		}
	}

	close(ch)
}

func loadSiftQueryVectors(path string) [][]float32 {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

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
			id:     int64(i),
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
					id:     int64(chunkNum*chunkSize + j),
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
