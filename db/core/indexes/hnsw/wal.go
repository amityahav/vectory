package hnsw

import (
	"encoding/binary"
	"github.com/pkg/errors"
	w "github.com/rosedblabs/wal"
)

const (
	AddVertex byte = iota
	SetEntryPointWithMaxLayer
	SetConnectionsAtLevel
	addConnectionAtLevel
)

type wal struct {
	f *w.WAL // wal operations are guarded internally by a lock
}

func newWal(path string) (*wal, error) {
	f, err := w.Open(w.Options{
		DirPath:       path,
		SegmentSize:   w.GB,
		SementFileExt: ".SEG",
		BlockCache:    32 * w.KB * 10,
		Sync:          false,
		BytesPerSync:  0,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed opening WAL at %s", path)
	}

	return &wal{f: f}, nil
}

func (w *wal) walReader() *w.Reader {
	return w.f.NewReader()
}

func (w *wal) addVertex(v *Vertex) error {
	/*
		bytes = [opcode, v.id, level], len(bytes) = 1 + 8 + 4
	*/
	level := len(v.connections) - 1
	bytes := make([]byte, 13)

	bytes[0] = AddVertex
	binary.LittleEndian.PutUint64(bytes[1:9], v.id)
	binary.LittleEndian.PutUint32(bytes[9:13], uint32(level))

	_, err := w.f.Write(bytes)

	return err
}

func (w *wal) setEntryPointWithMaxLayer(id uint64, level int) error {
	/*
		bytes = [opcode, id, level], len(bytes) = 1 + 8 + 4
	*/

	bytes := make([]byte, 13)

	bytes[0] = SetEntryPointWithMaxLayer
	binary.LittleEndian.PutUint64(bytes[1:9], id)
	binary.LittleEndian.PutUint32(bytes[9:13], uint32(level))

	_, err := w.f.Write(bytes)

	return err
}

func (w *wal) setConnectionsAtLevel(id uint64, level int, neighbors []uint64) error {
	/*
		bytes = [opcode, id, level, len(neighbors), neighbors], len(bytes) = 1 + 8 + 4 + 4 + 4*len(neighbors)
	*/

	bytes := make([]byte, 17+4*len(neighbors))

	bytes[0] = SetConnectionsAtLevel
	binary.LittleEndian.PutUint64(bytes[1:9], id)
	binary.LittleEndian.PutUint32(bytes[9:13], uint32(level))
	binary.LittleEndian.PutUint32(bytes[13:17], uint32(len(neighbors)))

	offset := 17
	for _, n := range neighbors {
		binary.LittleEndian.PutUint64(bytes[offset:], n)
		offset += 4
	}

	_, err := w.f.Write(bytes)

	return err
}

func (w *wal) addConnectionAtLevel(id uint64, level int, n uint64) error {
	/*
		bytes = [opcode, id, level, nid], len(bytes) = 1 + 8 + 4 + 8
	*/

	bytes := make([]byte, 21)

	bytes[0] = addConnectionAtLevel
	binary.LittleEndian.PutUint64(bytes[1:9], id)
	binary.LittleEndian.PutUint32(bytes[9:13], uint32(level))
	binary.LittleEndian.PutUint64(bytes[13:18], n)

	_, err := w.f.Write(bytes)

	return err
}
