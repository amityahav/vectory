package main

import (
	"Vectory/hnsw"
	bufio2 "bufio"
	"encoding/binary"
	"io"
	"log"
	"math"
	"os"
	"runtime"
	"sync"
)

func main() {
	//defer profile.Start(profile.TraceProfile, profile.ProfilePath(".")).Stop()
	//concurrent()
	sequential()

}

func sequential() {
	s := make([]hnsw.Vertex, 0, 10000)

	f, err := os.Open("hnsw/siftsmall/siftsmall_base.fvecs")
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

		v := hnsw.Vertex{
			Id:     int64(i),
			Vector: make([]float32, dim),
		}

		for j := 0; j < int(dim); j++ {
			v.Vector[j], err = readFloat32(buf, b)
			if err != nil {
				log.Fatal(err)
			}

		}

		s = append(s, v)

		if (i+1)%1 == 0 {
			log.Printf("loaded %d vectors", i+1)
		}
	}
}
func concurrent() {
	insertionChannel := make(chan *hnsw.Vertex)
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

				v := hnsw.Vertex{
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
