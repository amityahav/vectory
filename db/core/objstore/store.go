package objstore

import (
	"Vectory/entities/objstore"
	"encoding/binary"
	"errors"
	"git.mills.io/prologic/bitcask"
)

const (
	objectsDir = "object_storage"
	vectorsDir = "vectors_storage"
)

type Stores struct {
	// objects is a persistent storage for all objects in a collection
	objects *bitcask.Bitcask

	// vectors is a persistent storage for all vectors in a collection
	vectors *bitcask.Bitcask
}

func NewStores(filesPath string) (*Stores, error) {
	objects, err := bitcask.Open(filesPath + "/" + objectsDir)
	if err != nil {
		return nil, err
	}

	vectors, err := bitcask.Open(filesPath + "/" + vectorsDir)

	s := Stores{objects: objects, vectors: vectors}

	return &s, nil
}

func (s *Stores) PutObject(obj *objstore.Object) error {
	idBytes := make([]byte, 8) // TODO: can be reused
	binary.LittleEndian.PutUint64(idBytes, obj.Id)

	objBytes, err := obj.SerializeProperties()
	if err != nil {
		return err
	}

	vecBytes, err := obj.SerializeVector()
	if err != nil {
		return err
	}

	err = s.objects.Put(idBytes, objBytes)
	if err != nil {
		return err
	}

	return s.vectors.Put(idBytes, vecBytes)
}

func (s *Stores) GetObject(id uint64) (*objstore.Object, bool, error) {
	idBytes := make([]byte, 8) // TODO: can be reused
	binary.LittleEndian.PutUint64(idBytes, id)

	object, err := s.objects.Get(idBytes)
	if err != nil {
		if errors.Is(err, bitcask.ErrKeyNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	vector, err := s.vectors.Get(idBytes)
	if err != nil {
		if errors.Is(err, bitcask.ErrKeyNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	obj := objstore.Object{}
	err = obj.DeserializeProperties(object)
	if err != nil {
		return nil, false, err
	}

	err = obj.DeserializeVector(vector)
	if err != nil {
		return nil, false, err
	}

	obj.Id = id

	return &obj, true, nil
}

func (s *Stores) GetObjects(ids []uint64) ([]*objstore.Object, error) {
	objects := make([]*objstore.Object, 0, len(ids))
	for _, id := range ids {
		object, found, err := s.GetObject(id)
		if err != nil {
			return nil, err
		}

		if !found { // should not happen?
			continue
		}

		objects = append(objects, object)
	}

	return objects, nil
}

func (s *Stores) DeleteObject(id uint64) error {
	idBytes := make([]byte, 8) // TODO: can be reused
	binary.LittleEndian.PutUint64(idBytes, id)

	// TODO: currently delete only the actual object but keep its vector in the vectors store for index recovery and traversal
	return s.objects.Delete(idBytes)
}

func (s *Stores) GetVector(id uint64) ([]float32, bool, error) {
	idBytes := make([]byte, 8) // TODO: can be reused
	binary.LittleEndian.PutUint64(idBytes, id)

	vector, err := s.vectors.Get(idBytes)
	if err != nil {
		if errors.Is(err, bitcask.ErrKeyNotFound) {
			return nil, false, nil
		}

		return nil, false, err
	}

	obj := objstore.Object{}
	err = obj.DeserializeVector(vector)
	if err != nil {
		return nil, false, err
	}

	return obj.Vector, true, nil
}

// TODO: we can do better
func (s *Stores) GetVectorsStore() *bitcask.Bitcask {
	return s.vectors
}

func (s *Stores) Size() int {
	return s.objects.Len()
}

func (s *Stores) Close() error {
	if err := s.objects.Close(); err != nil {
		return err
	}

	if err := s.vectors.Close(); err != nil {
		return err
	}

	return nil
}
