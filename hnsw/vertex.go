package hnsw

// Vertex struct in a multi-layer graph
type Vertex struct {
	id          int64
	connections [][]int64

	vector []float32
}

func (v *Vertex) Init(level int64, mMax, mMax0 int) {
	v.connections = make([][]int64, level)
	for i := level - 1; i > 0; i-- {
		v.connections[i] = make([]int64, mMax)
	}

	v.connections[0] = make([]int64, mMax0)
}

func (v *Vertex) GetConnections(level int64) []int64 {
	return v.connections[level]
}
