package hnsw

func (h *Hnsw) Delete(id uint64) error {
	h.Lock()
	defer h.Unlock()

	// TODO: record in WAL
	h.deletedNodes[id] = struct{}{}

	return nil
}
