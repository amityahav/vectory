package embeddings

import (
	"Vectory/entities/objstore"
	"context"
)

type Embedder interface {
	Embed(ctx context.Context, objects []*objstore.Object) error
}
