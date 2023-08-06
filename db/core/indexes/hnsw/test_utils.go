package hnsw

import "math/rand"

func randomVector(dim int) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = rand.Float32()
	}

	return vec
}

func equalIndexes(h1 *Hnsw, h2 *Hnsw) bool {
	if h1.entrypointID != h2.entrypointID {
		return false
	}

	if h1.currentMaxLayer != h2.currentMaxLayer {
		return false
	}

	// check nodes equality

	return true
}
