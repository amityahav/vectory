package hnsw

import (
	"Vectory/entities/collection"
	bufio2 "bufio"
	"encoding/binary"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime"
	"sync"
)

type job struct {
	id     uint64
	vector []float32
}

func newHnsw(filesPath string) (*Hnsw, error) {
	return NewHnsw(collection.HnswParams{
		M:              64,
		MMax:           128,
		EfConstruction: 100,
		Ef:             100,
		Heuristic:      true,
		DistanceType:   "dot_product",
	}, filesPath, nil)
}

func randomVector(dim int) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = rand.Float32()
	}

	return vec
}

func insertInParallel(insertionChannel chan job, h *Hnsw, wg *sync.WaitGroup) {
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				job, ok := <-insertionChannel
				if !ok {
					break
				}

				err := h.Insert(job.vector, job.id)
				if err != nil {
					log.Fatal(err)
					return
				}
			}
		}()
	}
}

func buildIndexParallel(insertionChannel chan job) *Hnsw {
	hnsw, _ := newHnsw("")
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

				err := hnsw.Insert(job.vector, job.id)
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
			id:     uint64(i),
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

func loadSiftTruthVectors(path string) [][]uint64 {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	buf := bufio2.NewReader(f)
	b := make([]byte, 4)

	vectors := make([][]uint64, 100)

	for i := 0; i < 100; i++ {
		dim, err := readUint32(buf, b)
		if err != nil {
			log.Fatal(err)
		}

		vectors[i] = make([]uint64, dim)

		for j := 0; j < int(dim); j++ {
			fl32, err := readUint32(buf, b)
			if err != nil {
				log.Fatal(err)
			}

			vectors[i][j] = uint64(fl32)
		}
	}

	return vectors
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
