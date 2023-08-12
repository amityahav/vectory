package hnsw

import (
	"Vectory/db/core/index"
	"Vectory/db/core/index/distance"
	"Vectory/db/core/index/utils"
	"Vectory/db/core/objstore"
	"Vectory/entities"
	"Vectory/entities/collection"
	"fmt"
	"math"
	"sync"
)

var _ index.VectorIndex = &Hnsw{}

type Hnsw struct {
	sync.RWMutex
	m                int
	mMax             int
	mMax0            int
	efConstruction   int
	ef               int
	mL               float64
	entrypointID     uint64
	currentMaxLayer  int64
	nodes            map[uint64]*Vertex
	distFunc         func([]float32, []float32) float32
	selectNeighbors  func(*Vertex, []utils.Element, int) []uint64
	initialInsertion *sync.Once
	filesPath        string
	wal              *wal
}

func NewHnsw(params collection.HnswParams, filesPath string, store *objstore.ObjectStore) (*Hnsw, error) {
	h := Hnsw{
		m:                params.M,
		mMax:             params.MMax,
		ef:               params.Ef,
		efConstruction:   params.EfConstruction,
		nodes:            make(map[uint64]*Vertex), // TODO: change to an array
		initialInsertion: &sync.Once{},
		filesPath:        fmt.Sprintf("%s/%s", filesPath, "index"),
	}

	h.mMax0 = 2 * h.mMax
	h.mL = 1 / math.Log(float64(h.m))

	switch params.DistanceType {
	case entities.DotProduct:
		h.distFunc = distance.Dot
	case entities.Euclidean:
		h.distFunc = distance.EuclideanDistance
	}

	h.selectNeighbors = h.selectNeighborsSimple

	if params.Heuristic {
		h.selectNeighbors = h.selectNeighborsHeuristic
	}

	w, err := newWal(h.filesPath)
	if err != nil {
		return nil, err
	}

	h.wal = w

	if err = h.loadFromWAL(); err != nil {
		return nil, err
	}

	if err = h.populateVerticesVectors(store); err != nil {
		return nil, err
	}

	return &h, nil
}

func (h *Hnsw) Flush() error {
	return h.wal.flush()
}
