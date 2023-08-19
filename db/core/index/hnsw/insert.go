package hnsw

import (
	"Vectory/db/core/index/utils"
	"fmt"
)

func (h *Hnsw) Insert(vector []float32, vectorId uint64) error {
	var (
		first bool
		err   error
	)

	v := Vertex{
		id:     vectorId,
		vector: vector,
	}

	h.initialInsertion.Do(func() {
		if h.isEmpty() {
			err = h.insertFirstVertex(&v)
			if err != nil {
				return
			}

			first = true
		}

	})

	if err != nil {
		return fmt.Errorf("initialInsertion: %s", err.Error())
	}

	if first {
		return nil
	}

	vertexLayer := h.calculateLevelForVertex()
	v.Init(vertexLayer+1, h.mMax, h.mMax0)

	h.wal.addVertex(&v)

	h.Lock()
	h.nodes[v.id] = &v
	h.Unlock()

	h.RLock()
	entrypointID := h.entrypointID
	epVertex := h.nodes[entrypointID]
	currentMaxLayer := h.currentMaxLayer
	h.RUnlock()

	dist := h.calculateDistance(epVertex.vector, v.vector)

	var nearestNeighbors []utils.Element

	eps := make([]utils.Element, 0, 1)
	eps = append(eps, utils.Element{Id: entrypointID, Distance: dist})

	// Lookup Phase
	for l := currentMaxLayer; l > vertexLayer; l-- {
		nearestNeighbors = h.searchLayer(&v, eps, 1, l)
		eps[0] = nearestNeighbors[0]
	}

	// Construction Phase
	maxConn := h.mMax
	for l := min(currentMaxLayer, vertexLayer); l >= 0; l-- {
		nearestNeighbors = h.searchLayer(&v, eps, h.efConstruction, l)
		neighbors := h.selectNeighbors(&v, nearestNeighbors, h.m)

		h.wal.setConnectionsAtLevel(v.id, int(l), neighbors)

		v.SetConnections(l, neighbors)

		if l == 0 {
			maxConn = h.mMax0
		}

		for _, n := range neighbors {
			h.RLock()
			nVertex := h.nodes[n]
			h.RUnlock()

			nVertex.Lock()
			connections := nVertex.GetConnections(l)

			if len(connections) < maxConn {
				h.wal.addConnectionAtLevel(n, int(l), v.id)
				nVertex.AddConnection(l, v.id)
			} else { // pruning
				elems := make([]utils.Element, 0, len(connections)+1)

				elems = append(elems, utils.Element{
					Id:       v.id,
					Distance: h.calculateDistance(nVertex.vector, v.vector),
				})

				for _, nn := range connections {
					h.RLock()
					nnVertex := h.nodes[nn]
					h.RUnlock()

					elems = append(elems, utils.Element{Id: nn, Distance: h.calculateDistance(nVertex.vector, nnVertex.vector)})
				}

				newNeighbors := h.selectNeighbors(nVertex, elems, maxConn)
				h.wal.setConnectionsAtLevel(n, int(l), newNeighbors)
				nVertex.SetConnections(l, newNeighbors)
			}
			nVertex.Unlock()
		}

		eps = nearestNeighbors
	}

	h.Lock()
	if vertexLayer > currentMaxLayer {
		h.wal.setEntryPointWithMaxLayer(v.id, int(vertexLayer))
		h.entrypointID = v.id
		h.currentMaxLayer = vertexLayer
	}
	h.Unlock()

	return nil
}

func (h *Hnsw) insertFirstVertex(v *Vertex) error {
	h.Lock()
	defer h.Unlock()

	v.Init(1, -1, h.mMax0)

	h.wal.setEntryPointWithMaxLayer(v.id, 0)

	h.entrypointID = v.id
	h.currentMaxLayer = 0

	h.wal.addVertex(v)
	h.nodes[v.id] = v

	return nil
}