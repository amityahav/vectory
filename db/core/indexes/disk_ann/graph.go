package disk_ann

import (
	"encoding/binary"
	"math"
)

type Graph struct {
	maxDegree uint32

	vertices map[uint32]*Vertex
}

func newGraph(maxDegree uint32) *Graph {
	g := Graph{
		maxDegree: maxDegree, // TODO
		vertices:  map[uint32]*Vertex{},
	}

	return &g
}

func (g *Graph) addVertex(v *Vertex) {
	g.vertices[v.id] = v
}

func (g *Graph) serializeVertices(buff []byte, ids []uint32) int {
	var offset int

	for _, id := range ids {
		v := g.vertices[id]
		n := v.serialize(buff[offset:], g.maxDegree)
		offset += n
	}

	return offset
}

func (g *Graph) deserializeVertices(buff []byte, firstId uint32, numberOfVertices int) {
	var offset int

	for i := 0; i < numberOfVertices; i++ {
		v := Vertex{}
		v.deserialize(buff[offset:])

	}
}

type Vertex struct {
	id        uint32
	objId     uint32
	neighbors []uint32
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
		binary.LittleEndian.PutUint32(buff[offset:], n)
		offset += 4
	}

	// TODO: zero padding (needed?)
	for i := len(v.neighbors); i < int(maxDegree); i++ {
		binary.LittleEndian.PutUint32(buff[offset:], 0)
		offset += 4
	}

	return offset
}

func (v *Vertex) deserialize(buff []byte, dim uint32) {
	// TODO: calculate v id based on position on disk
	var offset int

	v.objId = binary.LittleEndian.Uint32(buff[offset:])
	offset += 4

	v.vector = make([]float32, 0, int(dim))

	for i := 0; i < int(dim); i++ {
		floatAsUint32 := binary.LittleEndian.Uint32(buff[offset:])
		offset += 4

		v.vector = append(v.vector, math.Float32frombits(floatAsUint32))
	}

	// TODO: extract floats until the first 0

}
