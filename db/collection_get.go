package db

import (
	objstoreentities "Vectory/entities/objstore"
	"context"
	"github.com/pkg/errors"
)

// Get returns the objects with objIds from the collection.
func (c *Collection) Get(objIds []uint64) ([]objstoreentities.Object, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	objects := make([]objstoreentities.Object, 0, len(objIds))
	for _, id := range objIds {
		obj, found, err := c.objStore.GetObject(id)
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

// SemanticSearch returns the approximate k-nn of obj.
func (c *Collection) SemanticSearch(ctx context.Context, obj *objstoreentities.Object, k int) ([]*objstoreentities.Object, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if err := c.embedObjectsIfNeeded(ctx, []*objstoreentities.Object{obj}); err != nil {
		return nil, err
	}

	results := c.vectorIndex.Search(obj.Vector, k)
	ids := make([]uint64, 0, k)

	for _, e := range results {
		ids = append(ids, e.Id)
	}

	return c.objStore.GetObjects(ids) // TODO: support returning distances, vectors of objects as well
}
