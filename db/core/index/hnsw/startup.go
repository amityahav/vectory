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
func (h *Hnsw) populateVerticesVectors(store *objstore.ObjectStore) error {
	if store == nil {
		return nil
	}

	s := store.GetStore()

	for k := range s.Keys() {
		o, _, err := store.GetObject(binary.LittleEndian.Uint64(k))
		if err != nil {
			return err
		}

		h.nodes[o.Id].vector = o.Vector
	}

	return nil
}
