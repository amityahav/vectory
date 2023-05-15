package hnsw

import (
	"Vectory/hnsw/distance"
	"math"
	"math/rand"
	"testing"
)

func TestHnsw(t *testing.T) {
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

	hnsw := NewHnsw(cfg)
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
	}, 10, cfg.efConstruction)

}

func randomVector(dim int) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = rand.Float32()
	}

	return vec
}
