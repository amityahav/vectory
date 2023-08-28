package hnsw

// TODO: delete currently only soft deletes the desired vertex. a mechanism for hard deletes ned to be implemented
func (h *Hnsw) Delete(id uint64) error {
	h.Lock()
	defer h.Unlock()

	h.wal.deleteVertex(id)
	h.deletedNodes[id] = struct{}{}

	return nil
}
