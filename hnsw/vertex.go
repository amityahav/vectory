package hnsw

type Vertex struct {
	id          int64
	level       int32
	connections [][]int64

	vector []float32
}
