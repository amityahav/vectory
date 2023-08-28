package hnsw

import (
	"fmt"
	"github.com/pkg/profile"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

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
