package hnsw

// Vertex struct in a multi-layer graph
type Vertex struct {
	id          int64
	connections [][]int64

	vector []float32
}

func (v *Vertex) Init(level int64, mMax, mMax0 int) {
	v.connections = make([][]int64, level)

	// initialising maximum connections to be + 1 in order to avoid allocating extra space when cap is full
	for i := level - 1; i > 0; i-- {
		v.connections[i] = make([]int64, 0, mMax+1)
	}

	v.connections[0] = make([]int64, 0, mMax0+1)
}

func (v *Vertex) GetConnections(level int64) []int64 {
	return v.connections[level]
}

func (v *Vertex) AddConnection(level, id int64) {
	v.connections[level] = append(v.connections[level], id)
}

func (v *Vertex) AddConnections(level int64, ids []int64) {
	v.connections[level] = append(v.connections[level], ids...)
}

func (v *Vertex) SetConnections(level int64, ids []int64) {
	v.connections[level] = ids
}
