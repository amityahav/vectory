package hnsw

import (
	"Vectory/db/core/objstore"
	"Vectory/entities/index"
	objstoreentities "Vectory/entities/objstore"
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func TestRestoreFromDisk(t *testing.T) {
	filesPath := "../tmp"
	defer os.RemoveAll(filesPath)

	store, err := objstore.NewStores(filesPath)
	require.NoError(t, err)

	dim := 128
	h, err := NewHnsw(index.DefaultHnswParams, filesPath, store)

	require.NoError(t, err)

	for i := 0; i < 100; i++ {
		vec := randomVector(dim)
		o := objstoreentities.Object{
			Id: uint64(i),
			Properties: map[string]interface{}{
				"title": "test",
			},
			Vector: vec,
		}

		require.NoError(t, h.Insert(vec, uint64(i)))
		require.NoError(t, store.PutObject(&o))
	}

	require.NoError(t, h.Delete(50))
	require.NoError(t, store.DeleteObject(50))

	require.NoError(t, h.Flush())

	// create new hnsw and restore from wal
	start := time.Now()
	hRestored, err := NewHnsw(index.DefaultHnswParams, filesPath, store)
	require.NoError(t, err)

	end := time.Since(start)

	fmt.Printf("reloading from disk took: %s\n", end)

	// compare old hnsw with restored hnsw
	require.Equal(t, h.entrypointID, hRestored.entrypointID)
	require.Equal(t, h.currentMaxLayer, hRestored.currentMaxLayer)
	require.Equal(t, h.nodes, hRestored.nodes)
	require.Equal(t, h.deletedNodes, hRestored.deletedNodes)
}
