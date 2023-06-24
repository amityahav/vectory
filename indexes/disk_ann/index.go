package disk_ann

type Index interface {
	Search(s *Vertex, q []float32, k int) ([]uint32, map[uint32]struct{})
}
