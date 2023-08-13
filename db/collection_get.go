package db

import (
	objstoreentities "Vectory/entities/objstore"
	"github.com/pkg/errors"
)

func (c *Collection) Get(objIds []uint64) ([]objstoreentities.Object, error) {
	objects := make([]objstoreentities.Object, 0, len(objIds))
	for _, id := range objIds {
		obj, found, err := c.objStore.Get(id)
		if err != nil {
			return nil, errors.Wrapf(err, "failed getting %d from object store", id)
		}

		if !found {
			continue
		}

		objects = append(objects, *obj)
	}

	return objects, nil
}

func (c *Collection) SemanticSearch(obj *objstoreentities.Object, k int) {
	// TODO: create embeddings from obj and get K-NN and then retrieve object ids from objStore
}
