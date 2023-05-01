package hnsw

func min(a, b int64) int64 {
	m := b
	if a < b {
		m = a
	}

	return m
}
