package indexes

type VectorIndex interface {
	// Insert a new vector and its corresponding objId
	Insert(vector []float32, objId uint32) error

	// Delete vertex corresponding with objId
	Delete(objId uint32) bool

	// Search for K-NN of vector
	Search(vector []float32, k int)
}

func NewVectorIndex() {

}
