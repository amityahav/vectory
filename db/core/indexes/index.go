package indexes

import "Vectory/db/core/indexes/utils"

type VectorIndex interface {
	// Insert a new vector and its corresponding objId
	Insert(vector []float32, vectorId uint64) error

	// Delete vertex corresponding with objId
	Delete(vectorId int64) bool

	// Search for K-NN of vector
	Search(q []float32, k int) []utils.Element
}

func NewVectorIndex() {

}
