package objstore

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestObject(t *testing.T) {
	t.Run("serialization-deserialization", func(t *testing.T) {
		obj := Object{
			Data:   "This is an object",
			Vector: []float32{0.123, 0.12, 12, 12.3},
		}

		b := obj.Serialize()

		obj2 := Object{}
		obj2.Deserialize(b)

		require.Equal(t, obj, obj2)
	})
}
