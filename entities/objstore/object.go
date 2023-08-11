package objstore

import (
	"encoding/binary"
	"math"
)

const (
	Text = iota
	Image
)

type Object struct {
	Id uint64
	//DataType int // TODO: currently supports only text objects
	Data   string
	Vector []float32
}

func (o *Object) Serialize() []byte {
	b := make([]byte, 4+len(o.Data)+4+4*len(o.Vector)) // Id + DataLen + Data + VecDim + Vector

	var offset int

	//binary.LittleEndian.PutUint32(b[offset:], uint32(o.DataType))
	//offset += 4

	binary.LittleEndian.PutUint32(b[offset:], uint32(len(o.Data)))
	offset += 4

	copy(b[offset:], o.Data)
	offset += len(o.Data)

	binary.LittleEndian.PutUint32(b[offset:], uint32(len(o.Vector)))
	offset += 4

	for _, f := range o.Vector {
		binary.LittleEndian.PutUint32(b[offset:], math.Float32bits(f))
		offset += 4
	}

	return b
}

func (o *Object) Deserialize(object []byte) {
	var offset int

	//o.DataType = int(binary.LittleEndian.Uint32(object[offset:]))
	//offset += 4

	dataLen := int(binary.LittleEndian.Uint32(object[offset:]))
	offset += 4

	dataBytes := make([]byte, dataLen)
	copy(dataBytes, object[offset:])
	o.Data = string(dataBytes)
	offset += dataLen

	dim := int(binary.LittleEndian.Uint32(object[offset:]))
	offset += 4

	vec := make([]float32, dim)
	for i := 0; i < dim; i++ {
		vec[i] = math.Float32frombits(binary.LittleEndian.Uint32(object[offset:]))
		offset += 4
	}

	o.Vector = vec
}
