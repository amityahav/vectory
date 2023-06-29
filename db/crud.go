package db

type CRUD interface {
	Insert(obj any)
	Delete(objId uint32)
	Update(objId uint32)
	Get(objIds []uint32)
	SemanticSearch(obj any, k int)
}
