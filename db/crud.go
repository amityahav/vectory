package db

import "Vectory/entities/objstore"

type CRUD interface {
	Insert(obj *objstore.Object) error
	InsertWithVector(obj *objstore.Object, vector []float32) error
	Update(obj *objstore.Object) error
	Delete(objId uint64) error
	Get(objIds []uint64) ([]objstore.Object, error)
	SemanticSearch(obj any, k int)
}
