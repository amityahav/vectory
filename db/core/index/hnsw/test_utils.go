package hnsw

import (
	"Vectory/entities"
	"Vectory/entities/collection"
	bufio2 "bufio"
	"encoding/binary"
	"io"
	"log"
	"math"
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
		M:              8,
		MMax:           16,
		EfConstruction: 400,
		Ef:             100,
		Heuristic:      true,
		DistanceType:   entities.Euclidean,
	}, filesPath, nil)
}

func randomVector(dim int) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = 23
	}

	return vec
}

func insertInParallel(insertionChannel chan job, h *Hnsw, wg *sync.WaitGroup) {
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				j, ok := <-insertionChannel
				if !ok {
					break
				}

				err := h.Insert(j.vector, j.id)
				if err != nil {
					log.Fatal(err)
					return
				}

				//err = h.Flush()
				//if err != nil {
				//	log.Fatal(err)
				//	return
				//}
			}
		}()
	}
	wg.Wait()
}

func buildIndexParallel(insertionChannel chan job, filesPath string) *Hnsw {
	hnsw, _ := newHnsw(filesPath)
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

func loadRandomVectors(ch chan job) {
	vec := randomVector(128)
	for i := 0; i < 10000; i++ {
		ch <- job{
			id:     uint64(i),
			vector: vec,
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
