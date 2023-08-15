package db

import (
	"Vectory/entities/objstore"
	"context"
)

type CRUD interface {
	Insert(ctx context.Context, obj *objstore.Object) error
	InsertBatch(ctx context.Context, objs []*objstore.Object) error
	InsertBatch2(ctx context.Context, objs []*objstore.Object) error
	Update(obj *objstore.Object) error
	Delete(objId uint64) error
	Get(objIds []uint64) ([]objstore.Object, error)
	SemanticSearch(ctx context.Context, obj *objstore.Object, k int) ([]*objstore.Object, error)
}
