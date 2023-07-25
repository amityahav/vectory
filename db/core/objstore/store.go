package objstore

import "git.mills.io/prologic/bitcask"

type ObjectStore struct {
	db *bitcask.Bitcask
}

func NewObjectStore(filesPath string) (*ObjectStore, error) {
	db, err := bitcask.Open(filesPath)
	if err != nil {
		return nil, err
	}

	s := ObjectStore{db: db}

	return &s, nil
}

func (s *ObjectStore) Put(id uint64, o *Object) {

}
