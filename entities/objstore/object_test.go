package objstore

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestObject(t *testing.T) {
	t.Run("serialization-deserialization", func(t *testing.T) {
		obj := Object{
			Properties: map[string]interface{}{
				"field1": "test1",
				"field2": "test2",
			},
			Vector: []float32{0.123, 0.12, 12, 12.3},
		}

		b, err := obj.SerializeProperties()
		require.NoError(t, err)

		obj2 := Object{}
		err = obj2.DeserializeProperties(b)
		require.NoError(t, err)

		require.Equal(t, obj, obj2)
	})
}
