package objstore

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
)

const (
	Text = iota
	Image
)

type Object struct {
	Id uint64
	//DataType int // TODO: currently supports only text objects
	Properties map[string]interface{}
	Vector     []float32
}

func (o *Object) Serialize() ([]byte, error) {
	propertiesJSON, err := json.Marshal(o.Properties)
	if err != nil {
		return nil, err
	}

	b := make([]byte, 4+len(propertiesJSON)+4+4*len(o.Vector)) // PropertiesLen + Properties + VecDim + Vector

	var offset int

	//binary.LittleEndian.PutUint32(b[offset:], uint32(o.DataType))
	//offset += 4

	binary.LittleEndian.PutUint32(b[offset:], uint32(len(propertiesJSON)))
	offset += 4

	copy(b[offset:], propertiesJSON)
	offset += len(propertiesJSON)

	binary.LittleEndian.PutUint32(b[offset:], uint32(len(o.Vector)))
	offset += 4

	for _, f := range o.Vector {
		binary.LittleEndian.PutUint32(b[offset:], math.Float32bits(f))
		offset += 4
	}

	return b, nil
}

func (o *Object) Deserialize(object []byte) error {
	var offset int

	//o.DataType = int(binary.LittleEndian.Uint32(object[offset:]))
	//offset += 4

	propertiesLen := int(binary.LittleEndian.Uint32(object[offset:]))
	offset += 4

	err := json.Unmarshal(object[offset:offset+propertiesLen], &o.Properties)
	if err != nil {
		return err
	}
	offset += propertiesLen

	dim := int(binary.LittleEndian.Uint32(object[offset:]))
	offset += 4

	vec := make([]float32, dim)
	for i := 0; i < dim; i++ {
		vec[i] = math.Float32frombits(binary.LittleEndian.Uint32(object[offset:]))
		offset += 4
	}

	o.Vector = vec

	return nil
}

func (o *Object) FlatProperties() string { // TODO: better handle different types
	var res string

	for k, v := range o.Properties {
		res += fmt.Sprintf("%s: %v,", k, v)
	}

	return res
}
