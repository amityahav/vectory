package hnsw

import (
	"Vectory/db/core/objstore"
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
)

// loadFromWAL builds index from wal if exists
func (h *Hnsw) loadFromWAL() error {
	r := h.wal.walReader()
	d := deserializer{state: h}

	for {
		record, err := r.Next()
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		if err = d.restore(record); err != nil {
			return errors.Wrap(err, "failed building index from WAL")
		}
	}

	return nil
}

// populateVerticesVectors prefetches all vectors from the object storage
// TODO: optimize it by implementing simple in-house KV store?
func (h *Hnsw) populateVerticesVectors(store *objstore.Stores) error {
	if store == nil {
		return nil
	}

	s := store.GetVectorsStore()

	for k := range s.Keys() {
		id := binary.LittleEndian.Uint64(k)
		vec, _, err := store.GetVector(id)
		if err != nil {
			return err
		}

		h.nodes[id].vector = vec
	}

	return nil
}
