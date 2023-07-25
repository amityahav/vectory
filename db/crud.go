package db

import "Vectory/db/core/objstore"

type CRUD interface {
	Insert(obj *objstore.Object) error
	InsertWithVector(obj *objstore.Object, vector []float32) error
	Delete(objId uint32)
	Update(objId uint32)
	Get(objIds []uint32)
	SemanticSearch(obj any, k int)
}
