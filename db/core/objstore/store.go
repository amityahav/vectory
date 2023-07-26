package objstore

import (
	"Vectory/entities/objstore"
	"encoding/binary"
	"git.mills.io/prologic/bitcask"
)

const storeDir = "object_storage"

type ObjectStore struct {
	db *bitcask.Bitcask
}

type ObjectStatus struct {
}

func NewObjectStore(filesPath string) (*ObjectStore, error) {
	db, err := bitcask.Open(filesPath + "/" + storeDir)
	if err != nil {
		return nil, err
	}

	s := ObjectStore{db: db}

	return &s, nil
}

func (s *ObjectStore) Put(obj *objstore.Object) error {
	idBytes := make([]byte, 8) // TODO: can be reused
	binary.LittleEndian.PutUint64(idBytes, obj.Id)

	return s.db.Put(idBytes, obj.Serialize())
}

func (s *ObjectStore) Get(id uint64) (*objstore.Object, error) {
	idBytes := make([]byte, 8) // TODO: can be reused
	binary.LittleEndian.PutUint64(idBytes, id)

	res, err := s.db.Get(idBytes)
	if err != nil {
		return nil, err
	}

	obj := objstore.Object{}
	obj.Deserialize(res)
	obj.Id = id

	return &obj, nil
}
