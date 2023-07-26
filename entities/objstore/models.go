package objstore

import "encoding/binary"

const (
	Text = iota
	Image
)

type Object struct {
	Id uint64
	//DataType int // TODO: currently supports only text objects
	Data string
}

func (o *Object) Serialize() []byte {
	b := make([]byte, 4+len(o.Data)) // Id + DataLen + Data

	var offset int

	//binary.LittleEndian.PutUint32(b[offset:], uint32(o.DataType))
	//offset += 4

	binary.LittleEndian.PutUint32(b[offset:], uint32(len(o.Data)))
	offset += 4

	copy(b[offset:], o.Data)
	offset += len(o.Data)

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

}
