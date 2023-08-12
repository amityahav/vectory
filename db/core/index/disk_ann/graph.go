package disk_ann

import (
	"encoding/binary"
	"math"
)

type Graph struct {
	vertices map[uint64]*Vertex
}

func newGraph() *Graph {
	g := Graph{
		vertices: map[uint64]*Vertex{},
	}

	return &g
}

func (g *Graph) addVertex(v *Vertex) {
	g.vertices[v.id] = v
}

func (g *Graph) serializeVertices(buff []byte, ids []uint64, maxDegree uint32) int {
	var offset int

	for _, id := range ids {
		v := g.vertices[id]
		n := v.serialize(buff[offset:], maxDegree)
		offset += n
	}

	return offset
}

func (g *Graph) deserializeVertices(buff []byte, dim uint32, maxDegree uint32, numberOfVertices uint64, currId uint64) int {
	var offset int

	for i := 0; i < int(numberOfVertices); i++ {
		v := Vertex{}
		v.id = currId + uint64(i)
		n := v.deserialize(buff[offset:], dim, maxDegree)
		g.vertices[v.id] = &v

		offset += n
	}

	return offset
}

type Vertex struct {
	id        uint64
	objId     uint32
	neighbors []uint64
	vector    []float32
}

func (v *Vertex) serialize(buff []byte, maxDegree uint32) int {
	var offset int

	binary.LittleEndian.PutUint32(buff[offset:], v.objId)
	offset += 4

	for _, float := range v.vector {
		floatAsUint32 := math.Float32bits(float)
		binary.LittleEndian.PutUint32(buff[offset:], floatAsUint32)
		offset += 4
	}

	for _, n := range v.neighbors {
		binary.LittleEndian.PutUint64(buff[offset:], n)
		offset += 8
	}

	// TODO: zero padding (needed?)
	for i := len(v.neighbors); i < int(maxDegree); i++ {
		binary.LittleEndian.PutUint32(buff[offset:], 0)
		offset += 4
	}

	return offset
}

func (v *Vertex) deserialize(buff []byte, dim uint32, maxDegree uint32) int {
	var offset int

	v.objId = binary.LittleEndian.Uint32(buff[offset:])
	offset += 4

	v.vector = make([]float32, 0, int(dim))

	for i := 0; i < int(dim); i++ {
		floatAsUint32 := binary.LittleEndian.Uint32(buff[offset:])
		offset += 4

		v.vector = append(v.vector, math.Float32frombits(floatAsUint32))
	}

	// TODO: can be optimized for allocations if we store v's degree on disk and then allocate its neighbors size accordingly
	var remainder int

	for i := 0; i < int(maxDegree); i++ {
		n := binary.LittleEndian.Uint64(buff[offset:])
		offset += 4

		if n == 0 { // no more neighbors
			remainder = int(maxDegree) - i - 1 // paddings remainder
			break
		}

		v.neighbors = append(v.neighbors, n)
	}

	return offset + remainder*4
}
