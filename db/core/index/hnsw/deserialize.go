package hnsw

import (
	"encoding/binary"
	"fmt"
)

type deserializer struct {
	state *Hnsw
}

func (d *deserializer) restore(record []byte) error {
	op := d.getOp(record)

	switch op {
	case AddVertex:
		d.addVertex(record[1:])
	case SetEntryPointWithMaxLayer:
		d.setEntryPointWithMaxLayer(record[1:])
	case SetConnectionsAtLevel:
		d.setConnectionsAtLevel(record[1:])
	case addConnectionAtLevel:
		d.addConnectionAtLevel(record[1:])
	case deleteVertex:
		d.deleteVertex(record[1:])
	default:
		return fmt.Errorf("unkown opcode %d", op)
	}

	return nil
}

func (d *deserializer) getOp(record []byte) byte {
	return record[0]
}

func (d *deserializer) addVertex(record []byte) {
	var offset int
	v := Vertex{}

	v.id = binary.LittleEndian.Uint64(record[offset:])
	offset += 8

	level := int64(binary.LittleEndian.Uint32(record[offset:]))
	offset += 4

	v.Init(level+1, d.state.mMax, d.state.mMax0)
	d.state.nodes[v.id] = &v
}

func (d *deserializer) setEntryPointWithMaxLayer(record []byte) {
	var offset int

	d.state.entrypointID = binary.LittleEndian.Uint64(record[offset:])
	offset += 8

	d.state.currentMaxLayer = int64(binary.LittleEndian.Uint32(record[offset:]))
	offset += 4

}
func (d *deserializer) addConnectionAtLevel(record []byte) {
	var offset int

	id := binary.LittleEndian.Uint64(record[offset:])
	offset += 8

	level := binary.LittleEndian.Uint32(record[offset:])
	offset += 4

	nid := binary.LittleEndian.Uint64(record[offset:])
	offset += 8

	v := d.state.nodes[id]

	v.AddConnection(int64(level), nid)
}
func (d *deserializer) setConnectionsAtLevel(record []byte) {
	var offset int

	id := binary.LittleEndian.Uint64(record[offset:])
	offset += 8

	level := binary.LittleEndian.Uint32(record[offset:])
	offset += 4

	size := binary.LittleEndian.Uint32(record[offset:])
	offset += 4

	v := d.state.nodes[id]
	neighbors := make([]uint64, int(size))

	for i := 0; i < int(size); i++ {
		neighbors[i] = binary.LittleEndian.Uint64(record[offset:])
		offset += 8
	}

	v.SetConnections(int64(level), neighbors)
}

func (d *deserializer) deleteVertex(record []byte) {
	var offset int

	id := binary.LittleEndian.Uint64(record[offset:])
	offset += 8

	d.state.deletedNodes[id] = struct{}{}
}
