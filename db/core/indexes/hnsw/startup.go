package hnsw

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
)

type hnswState struct {
	nodes           map[uint64]*Vertex
	entrypointID    uint64
	currentMaxLayer int64
}

type deserializer struct {
	state *hnswState
}

// init builds index from wal if exists
func (h *Hnsw) init() error {
	walReader := h.wal.walReader()
	state := hnswState{
		nodes:           make(map[uint64]*Vertex),
		entrypointID:    0,
		currentMaxLayer: 0,
	}

	d := deserializer{state: &state}

	for {
		record, _, err := walReader.Next()
		if err == io.EOF {
			break
		}

		if err = d.restore(record); err != nil {
			return errors.Wrap(err, "failed building index from WAL")
		}
	}

	h.nodes = state.nodes
	h.entrypointID = state.entrypointID
	h.currentMaxLayer = state.currentMaxLayer

	return nil
}

func (d *deserializer) restore(record []byte) error {
	op := d.getOp(record)

	switch op {
	case AddVertex:
		d.addVertex(record[1:])
	case SetEntryPointWithMaxLayer:
		d.setEntryPointWithMaxLayer(record[1:])
	case SetConnectionAtLevel:
		d.setConnectionAtLevel(record[1:])
	case SetConnectionsAtLevel:
		d.setConnectionsAtLevel(record[1:])
	default:
		return fmt.Errorf("unkown opcode %d", op)
	}

	return nil
}

func (d *deserializer) getOp(record []byte) byte {
	return record[0]
}

func (d *deserializer) addVertex(record []byte) {

}

func (d *deserializer) setEntryPointWithMaxLayer(record []byte) {

}
func (d *deserializer) setConnectionAtLevel(record []byte) {

}
func (d *deserializer) setConnectionsAtLevel(record []byte) {

}
