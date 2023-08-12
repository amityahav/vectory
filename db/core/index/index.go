package index

import "Vectory/db/core/index/utils"

type VectorIndex interface {
	// Insert a new vector and its corresponding objId
	Insert(vector []float32, vectorId uint64) error

	// Delete vertex corresponding with objId
	Delete(id uint64) error

	// Search for K-NN of vector
	Search(q []float32, k int) []utils.Element

	// Flush WAL to disk
	Flush() error
}
