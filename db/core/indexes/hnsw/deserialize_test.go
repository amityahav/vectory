package hnsw

import (
	"Vectory/entities/collection"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestDeserialize(t *testing.T) {
	filesPath := "./tmp"
	defer os.RemoveAll(filesPath)

	dim := 128
	h, err := NewHnsw(collection.HnswParams{
		M:              64,
		MMax:           128,
		EfConstruction: 100,
		Ef:             100,
		Heuristic:      true,
		DistanceType:   "dot_product",
	}, filesPath)

	require.NoError(t, err)

	for i := 0; i < 1000; i++ {
		err = h.Insert(randomVector(dim), uint64(i))
		require.NoError(t, err)
	}

	// create new hnsw and restore from wal
	hRestored, err := NewHnsw(collection.HnswParams{
		M:              64,
		MMax:           128,
		EfConstruction: 100,
		Ef:             100,
		Heuristic:      true,
		DistanceType:   "dot_product",
	}, filesPath)

	require.NoError(t, err)

	// compare old hnsw with restored hnsw
	require.True(t, equalIndexes(h, hRestored))
}

// benchmark the time taken for restoring hnsw from disk
