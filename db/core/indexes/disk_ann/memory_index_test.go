package disk_ann

import (
	"Vectory/db/core/indexes/distance"
	"github.com/stretchr/testify/require"
	"math/rand"
	"os"
	"sync"
	"testing"
)

func TestName(t *testing.T) {
	mi := newMemoryIndex(distance.Dot, &sync.Map{}, 1, 128, 128)
	listSize := 100
	a := float32(1.3)

	t.Run("insertion", func(t *testing.T) {
		for i := mi.firstId; i <= 100; i++ {
			err := mi.Insert(randomVector(mi.dim), listSize, a, i, i)
			require.NoError(t, err)

		}
	})

	t.Run("snapshot", func(t *testing.T) {
		path := "./ro_0.vctry"

		err := mi.Snapshot(path)
		require.NoError(t, err)

		d, err := newDal(path)
		require.NoError(t, err)

		restoredIndex, err := d.readIndex()
		require.NoError(t, err)

		require.Equal(t, mi.graph, restoredIndex.graph) // TODO: should check for the whole graphs

		err = os.Remove(path)
		require.NoError(t, err)
	})
}

func randomVector(dim uint32) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = rand.Float32()
	}

	return vec
}
