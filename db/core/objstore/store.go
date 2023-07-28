package objstore

import (
	"Vectory/entities/objstore"
	"encoding/binary"
	"errors"
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

func (s *ObjectStore) Get(id uint64) (*objstore.Object, bool, error) {
	idBytes := make([]byte, 8) // TODO: can be reused
	binary.LittleEndian.PutUint64(idBytes, id)

	res, err := s.db.Get(idBytes)
	if err != nil {
		if errors.Is(err, bitcask.ErrKeyNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	obj := objstore.Object{}
	obj.Deserialize(res)
	obj.Id = id

	return &obj, true, nil
}

func (s *ObjectStore) Delete(id uint64) error {
	idBytes := make([]byte, 8) // TODO: can be reused
	binary.LittleEndian.PutUint64(idBytes, id)

	return s.db.Delete(idBytes)
}
