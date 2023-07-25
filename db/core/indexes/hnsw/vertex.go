package hnsw

import "sync"

// Vertex struct in a multi-layer graph
type Vertex struct {
	sync.Mutex

	id          uint64
	connections [][]uint64
	vector      []float32
}

func (v *Vertex) Init(level int64, mMax, mMax0 int) {
	v.connections = make([][]uint64, level)

	for i := level - 1; i > 0; i-- {
		v.connections[i] = make([]uint64, 0, mMax)
	}

	v.connections[0] = make([]uint64, 0, mMax0)
}

func (v *Vertex) GetConnections(level int64) []uint64 {
	return v.connections[level]
}

func (v *Vertex) AddConnection(level int64, id uint64) {
	v.connections[level] = append(v.connections[level], id)
}

func (v *Vertex) AddConnections(level int64, ids []uint64) {
	v.connections[level] = append(v.connections[level], ids...)
}

func (v *Vertex) SetConnections(level int64, ids []uint64) {
	v.connections[level] = ids
}

func (v *Vertex) ResetConnections(level int64) {
	v.connections[level] = v.connections[level][:0]
}
