package db

import (
	objstoreentities "Vectory/entities/objstore"
	"context"
)

// Insert inserts one object to the collection.
func (c *Collection) Insert(ctx context.Context, obj *objstoreentities.Object) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.validateObjectsMappings([]*objstoreentities.Object{obj}); err != nil {
		return err
	}

	if err := c.embedObjectsIfNeeded(ctx, []*objstoreentities.Object{obj}); err != nil {
		return err
	}

	if err := c.insert(obj); err != nil {
		return err
	}

	if err := c.vectorIndex.Flush(); err != nil {
		return err
	}

	return nil
}

// InsertBatch inserts a batch of objects to the collection.
// it does that by splitting the batch into equally sized chunks distributed among multiple worker threads.
func (c *Collection) InsertBatch(ctx context.Context, objs []*objstoreentities.Object) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.validateObjectsMappings(objs); err != nil {
		return err
	}

	if err := c.embedObjectsIfNeeded(ctx, objs); err != nil {
		return err
	}

	workers := c.wp.MaxWorkers()
	objsInChunk := len(objs) / workers
	group, ctx := c.wp.GroupContext(ctx)
	defer ctx.Done()

	var offset int

	for i := 0; i < workers; i++ {
		end := offset + objsInChunk
		if i == workers-1 { // remainder
			end += len(objs) % workers
		}

		workerFunc := func(start, end int) func() error {
			return func() error {
				for _, obj := range objs[start:end] {
					err := c.insert(obj)
					if err != nil { // TODO: dont fail the entire batch?
						return err
					}
				}

				return nil
			}
		}(offset, end)

		group.Submit(workerFunc)

		offset = end
	}

	err := group.Wait()
	if err != nil {
		return err
	}

	return c.vectorIndex.Flush()
}

// InsertBatch2 is the same as InsertBatch but creates a channel from objs and share it among the worker threads.
func (c *Collection) InsertBatch2(ctx context.Context, objs []*objstoreentities.Object) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if err := c.validateObjectsMappings(objs); err != nil {
		return err
	}

	if err := c.embedObjectsIfNeeded(ctx, objs); err != nil {
		return err
	}

	objects := make(chan *objstoreentities.Object, len(objs))
	for _, o := range objs {
		objects <- o
	}
	close(objects)

	group, ctx := c.wp.GroupContext(ctx)
	defer ctx.Done()

	for i := 0; i < c.wp.MaxWorkers(); i++ {
		group.Submit(func() error {
			for {
				obj, ok := <-objects
				if !ok {
					break
				}

				err := c.insert(obj)
				if err != nil {
					return err
				}
			}

			return nil
		})
	}

	err := group.Wait()
	if err != nil {
		return err
	}

	return c.vectorIndex.Flush()
}

// insert handles the actual insertion of the object both to the object storage and index.
// it will also create an embedding for the data if vector is not specified.
func (c *Collection) insert(obj *objstoreentities.Object) error {
	id, err := c.idCounter.FetchAndInc()
	if err != nil {
		return err
	}

	obj.Id = id

	if err = c.objStore.Put(obj); err != nil {
		return err
	}

	return c.vectorIndex.Insert(obj.Vector, obj.Id)
}
