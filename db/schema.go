package db

import "Vectory/db/core/indexes"

var _ CRUD = &Schema{}

type Schema struct {
	id       any
	name     any
	objStore any
	index    indexes.VectorIndex
	logger   any
	embbeder any
	wal      any
}

func (s *Schema) Insert(obj any) {
	// TODO: store obj, create embedding and store objId and vector in index
}

func (s *Schema) Delete(objId uint32) {
	// TODO: delete in objStore and in index
}

func (s *Schema) Update(objId uint32) {
	// TODO: delete both in index and objStore and create again
}

func (s *Schema) Get(objIds []uint32) {
	// TODO: get objects with objIds from objStore
}

func (s *Schema) SemanticSearch(obj any, k int) {
	// TODO: create embeddings from obj and get K-NN and then retrieve object ids from objStore
}
