package db

import "Vectory/db/core/indexes"

var _ CRUD = &DB{}

// TODO: act as coordinator?
type DB struct {
	vectorIndex indexes.VectorIndex
	objStore    any
	logger      any
	config      any
}

func NewDB() {

}

func (db *DB) Insert(obj any) {
	// TODO: store obj, create embedding and store objId and vector in index
}

func (db *DB) Delete(objId uint32) {
	// TODO: delete in objStore and in index
}

func (db *DB) Update(objId uint32) {
	// TODO: delete both in index and objStore and create again
}

func (db *DB) Get(objIds []uint32) {
	// TODO: get objects with objIds from objStore
}

func (db *DB) SemanticSearch(obj any, k int) {
	// TODO: create embeddings from obj and get K-NN and then retrieve object ids from objStore
}
