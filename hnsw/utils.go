package hnsw

func min(a, b int64) int64 {
	m := b
	if a < b {
		m = a
	}

	return m
}

type Set[T comparable] map[T]struct{}

func NewSet[T comparable]() Set[T] {
	return Set[T]{}
}

func (s Set[T]) Add(elem T) {
	s[elem] = struct{}{}
	return
}

func (s Set[T]) Contains(elem T) bool {
	_, ok := s[elem]
	return ok
}
