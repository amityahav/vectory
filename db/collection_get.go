package db

import (
	"Vectory/entities/collection"
	objstoreentities "Vectory/entities/objstore"
	"context"
	"github.com/pkg/errors"
)

// Get returns the objects with objIds from the collection.
func (c *Collection) Get(objIds []uint64) ([]objstoreentities.Object, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, ErrCollectionClosed
	}

	objects := make([]objstoreentities.Object, 0, len(objIds))
	for _, id := range objIds {
		obj, found, err := c.stores.GetObject(id)
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
func (c *Collection) SemanticSearch(ctx context.Context, obj *objstoreentities.Object, k int) (*collection.SemanticSearchResult, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if c.closed {
		return nil, ErrCollectionClosed
	}

	if err := c.embedObjectsIfNeeded(ctx, []*objstoreentities.Object{obj}); err != nil {
		return nil, err
	}

	results := c.vectorIndex.Search(obj.Vector, k)
	ids := make([]uint64, 0, k)

	for _, e := range results {
		ids = append(ids, e.Id)
	}

	objs, err := c.stores.GetObjects(ids)
	if err != nil {
		return nil, err
	}

	res := collection.SemanticSearchResult{
		Hits: len(results),
	}

	resObjs := make([]objstoreentities.ObjectWithDistance, 0, len(results))
	for i := 0; i < len(results); i++ {
		resObjs = append(resObjs, objstoreentities.ObjectWithDistance{
			Id:         objs[i].Id,
			Properties: objs[i].Properties,
			Distance:   results[i].Distance,
		})
	}

	res.Objects = resObjs

	return &res, nil
}
