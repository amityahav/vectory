package hnsw

import (
	"github.com/pkg/errors"
	"io"
)

// init builds index from wal if exists
func (h *Hnsw) init() error {
	walReader := h.wal.walReader()
	d := deserializer{state: h}

	for {
		record, _, err := walReader.Next()
		if err == io.EOF {
			break
		}

		if err = d.restore(record); err != nil {
			return errors.Wrap(err, "failed building index from WAL")
		}
	}

	return nil
}
