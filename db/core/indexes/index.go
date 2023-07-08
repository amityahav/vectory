package indexes

import "Vectory/db/core/indexes/utils"

var SupportedIndexes = map[string]struct{}{
	"disk_ann": {},
	"hnsw":     {},
}

type VectorIndex interface {
	// Insert a new vector and its corresponding objId
	Insert(vector []float32, objId uint32) error

	// Delete vertex corresponding with objId
	Delete(objId uint32) bool

	// Search for K-NN of vector
	Search(q []float32, k int) []utils.Element
}

func NewVectorIndex() {

}
