package disk_ann

type Graph struct {
	maxDegree int

	vertices map[uint32]*Vertex
}

func newGraph(maxDegree int) *Graph {
	g := Graph{
		maxDegree: maxDegree, // TODO
		vertices:  map[uint32]*Vertex{},
	}

	return &g
}

func (g *Graph) size() int {
	return len(g.vertices)
}

func (g *Graph) addVertex(v *Vertex) {
	g.vertices[v.Id] = v
}

type Vertex struct {
	Id        uint32
	DataId    uint32
	Neighbors []uint32
	Vector    []float32
}
