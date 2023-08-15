package db

import objstoreentities "Vectory/entities/objstore"

// Update updates obj in the collection.
func (c *Collection) Update(obj *objstoreentities.Object) error {
	// TODO: handle race conditions
	return nil
}
