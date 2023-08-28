package hnsw

import (
	"encoding/binary"
	"github.com/pkg/errors"
	w "github.com/tidwall/wal"
	"io"
	"sync"
)

const (
	AddVertex byte = iota
	SetEntryPointWithMaxLayer
	SetConnectionsAtLevel
	addConnectionAtLevel
	deleteVertex
)

type wal struct {
	mu     sync.RWMutex
	f      *w.Log
	batch  *w.Batch
	seqNum uint64
}

func newWal(path string) (*wal, error) {
	f, err := w.Open(path, w.DefaultOptions)
	if err != nil {
		return nil, errors.Wrapf(err, "failed opening WAL at %s", path)
	}

	n, err := f.LastIndex()
	if err != nil {
		return nil, errors.Wrapf(err, "failed retrieving WAL's last sequence number")
	}
	n++

	return &wal{f: f,
		batch:  new(w.Batch),
		seqNum: n,
	}, nil
}

func (w *wal) flush() error {
	if err := w.f.WriteBatch(w.batch); err != nil {
		return err
	}

	w.batch.Clear()

	return nil
}

func (w *wal) addVertex(v *Vertex) {
	/*
		bytes = [opcode, v.id, level], len(bytes) = 1 + 8 + 4
	*/
	level := len(v.connections) - 1
	bytes := make([]byte, 13)

	bytes[0] = AddVertex
	binary.LittleEndian.PutUint64(bytes[1:9], v.id)
	binary.LittleEndian.PutUint32(bytes[9:], uint32(level))

	w.writeBatch(bytes)
}

func (w *wal) setEntryPointWithMaxLayer(id uint64, level int) {
	/*
		bytes = [opcode, id, level], len(bytes) = 1 + 8 + 4
	*/

	bytes := make([]byte, 13)

	bytes[0] = SetEntryPointWithMaxLayer
	binary.LittleEndian.PutUint64(bytes[1:9], id)
	binary.LittleEndian.PutUint32(bytes[9:], uint32(level))

	w.writeBatch(bytes)
}

func (w *wal) setConnectionsAtLevel(id uint64, level int, neighbors []uint64) {
	/*
		bytes = [opcode, id, level, len(neighbors), neighbors], len(bytes) = 1 + 8 + 4 + 4 + 8*len(neighbors)
	*/

	bytes := make([]byte, 17+8*len(neighbors))

	bytes[0] = SetConnectionsAtLevel
	binary.LittleEndian.PutUint64(bytes[1:9], id)
	binary.LittleEndian.PutUint32(bytes[9:13], uint32(level))
	binary.LittleEndian.PutUint32(bytes[13:17], uint32(len(neighbors)))

	offset := 17
	for _, n := range neighbors {
		binary.LittleEndian.PutUint64(bytes[offset:], n)
		offset += 8
	}

	w.writeBatch(bytes)
}

func (w *wal) addConnectionAtLevel(id uint64, level int, n uint64) {
	/*
		bytes = [opcode, id, level, nid], len(bytes) = 1 + 8 + 4 + 8
	*/

	bytes := make([]byte, 21)

	bytes[0] = addConnectionAtLevel
	binary.LittleEndian.PutUint64(bytes[1:9], id)
	binary.LittleEndian.PutUint32(bytes[9:13], uint32(level))
	binary.LittleEndian.PutUint64(bytes[13:], n)

	w.writeBatch(bytes)
}

func (w *wal) deleteVertex(id uint64) {
	/*
		bytes = [opcode, id] = 1 + 8
	*/

	bytes := make([]byte, 9)

	bytes[0] = deleteVertex
	binary.LittleEndian.PutUint64(bytes[1:9], id)

	w.writeBatch(bytes)
}

func (w *wal) read(seqNum uint64) ([]byte, error) {
	return w.f.Read(seqNum)
}

func (w *wal) writeBatch(data []byte) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.batch.Write(w.seqNum, data)
	w.seqNum++
}

type walReader struct {
	wal *wal
	pos uint64
}

func (w *wal) walReader() *walReader {
	return &walReader{
		wal: w,
		pos: 1,
	}
}

func (r *walReader) Next() ([]byte, error) {
	data, err := r.wal.read(r.pos)
	if err != nil {
		if err == w.ErrNotFound {
			return nil, io.EOF
		}

		return nil, err
	}

	r.pos++

	return data, nil
}
