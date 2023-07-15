package disk_ann

import (
	"Vectory/db/core/indexes/distance"
	"github.com/pkg/profile"
	"github.com/stretchr/testify/require"
	"log"
	"math/rand"
	"sync"
	"testing"
)

func TestMemoryIndex(t *testing.T) {
	defer profile.Start(profile.CPUProfile, profile.ProfilePath("./profile")).Stop()
	mi := newMemoryIndex(distance.Dot, &sync.Map{}, 1, 128, 128)
	listSize := 100
	a := float32(1.3)

	t.Run("insertion", func(t *testing.T) {
		for i := mi.firstId; i <= 134; i++ {
			if i == 134 {
				log.Printf("hi")
			}
			err := mi.insert(randomVector(mi.dim), listSize, a, i, i)
			require.NoError(t, err)
			if i%100 == 0 {
				log.Printf("inserted %d", i)
			}

		}
	})

	//t.Run("snapshot", func(t *testing.T) {
	//	path := "./ro_0.vctry"
	//
	//	err := mi.snapshot(path)
	//	require.NoError(t, err)
	//
	//	d, err := newDal(path)
	//	require.NoError(t, err)
	//
	//	restoredIndex, err := d.readIndex()
	//	require.NoError(t, err)
	//
	//	require.Equal(t, mi.graph, restoredIndex.graph) // TODO: should check for the whole graphs
	//
	//	err = os.Remove(path)
	//	require.NoError(t, err)
	//})
}

func randomVector(dim uint32) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = rand.Float32()
	}

	return vec
}
