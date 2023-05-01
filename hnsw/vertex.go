package hnsw

type Vertex struct {
	id          int64
	connections [][]int64

	vector []float32
}

func (v *Vertex) Init(level, mMax, mMax0 int64) {
	v.connections = make([][]int64, level)
	for i := level - 1; i > 0; i-- {
		v.connections[i] = make([]int64, mMax)
	}

	v.connections[0] = make([]int64, mMax0)
}
