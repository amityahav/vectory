package disk_ann

import (
	"Vectory/db/core/indexes/distance"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"os"
	"sync"
	"testing"
)

func TestMemoryIndex(t *testing.T) {
	mi := newMemoryIndex(distance.Dot, &sync.Map{}, 1, 1, 1)
	listSize := 2
	a := float32(1.3)

	t.Run("insertion", func(t *testing.T) {
		for i := mi.firstId; i <= 100; i++ {
			if i == 134 {
				log.Printf("hi")
			}
			err := mi.insert(randomVector(mi.dim), listSize, a, i, i)
			require.NoError(t, err)
			log.Printf("inserted %d", i)

		}
	})

	t.Run("snapshot", func(t *testing.T) {
		path := "./ro_0.vctry"

		err := mi.snapshot(path)
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
