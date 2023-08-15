package embeddings

import (
	"Vectory/entities/objstore"
	"context"
	"math/rand"
)

const FakeEmbedder = "fake_embedder"

// Fake used for tests purposes
type fake struct {
}

func NewFakeEmbedder() *fake {
	return &fake{}
}

func (e *fake) Embed(_ context.Context, objects []*objstore.Object) error {
	for i := 0; i < len(objects); i++ {
		objects[i].Vector = randomVector(128)
	}

	return nil
}

func randomVector(dim int) []float32 {
	vec := make([]float32, dim)

	for i := range vec {
		vec[i] = rand.Float32()
	}

	return vec
}
