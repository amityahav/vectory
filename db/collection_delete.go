package db

import "github.com/pkg/errors"

// Delete deletes an object with objId from the collection.
func (c *Collection) Delete(objId uint64) error {
	// TODO: delete in objStore and in index
	_, found, err := c.objStore.GetObject(objId)
	if err != nil {
		return errors.Wrapf(err, "failed getting %d from object store", objId)
	}

	if !found {
		// nothing to do
		return nil
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	err = c.objStore.Delete(objId)
	if err != nil {
		return errors.Wrapf(err, "failed deleting %d from object store", objId)
	}

	err = c.vectorIndex.Delete(objId)
	if err != nil {
		return errors.Wrapf(err, "failed deleting %d from vector index", objId)
	}

	return c.vectorIndex.Flush()
}
