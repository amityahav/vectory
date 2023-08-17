// Code generated by ent, DO NOT EDIT.

package ent

import (
	"Vectory/gen/ent/collection"
	"Vectory/gen/ent/predicate"
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/dialect/sql/sqljson"
	"entgo.io/ent/schema/field"
)

// CollectionUpdate is the builder for updating Collection entities.
type CollectionUpdate struct {
	config
	hooks    []Hook
	mutation *CollectionMutation
}

// Where appends a list predicates to the CollectionUpdate builder.
func (cu *CollectionUpdate) Where(ps ...predicate.Collection) *CollectionUpdate {
	cu.mutation.Where(ps...)
	return cu
}

// SetName sets the "name" field.
func (cu *CollectionUpdate) SetName(s string) *CollectionUpdate {
	cu.mutation.SetName(s)
	return cu
}

// SetIndexType sets the "index_type" field.
func (cu *CollectionUpdate) SetIndexType(s string) *CollectionUpdate {
	cu.mutation.SetIndexType(s)
	return cu
}

// SetDataType sets the "data_type" field.
func (cu *CollectionUpdate) SetDataType(s string) *CollectionUpdate {
	cu.mutation.SetDataType(s)
	return cu
}

// SetEmbedderType sets the "embedder_type" field.
func (cu *CollectionUpdate) SetEmbedderType(s string) *CollectionUpdate {
	cu.mutation.SetEmbedderType(s)
	return cu
}

// SetIndexParams sets the "index_params" field.
func (cu *CollectionUpdate) SetIndexParams(m map[string]interface{}) *CollectionUpdate {
	cu.mutation.SetIndexParams(m)
	return cu
}

// SetEmbedderConfig sets the "embedder_config" field.
func (cu *CollectionUpdate) SetEmbedderConfig(m map[string]interface{}) *CollectionUpdate {
	cu.mutation.SetEmbedderConfig(m)
	return cu
}

// SetMappings sets the "mappings" field.
func (cu *CollectionUpdate) SetMappings(s []string) *CollectionUpdate {
	cu.mutation.SetMappings(s)
	return cu
}

// AppendMappings appends s to the "mappings" field.
func (cu *CollectionUpdate) AppendMappings(s []string) *CollectionUpdate {
	cu.mutation.AppendMappings(s)
	return cu
}

// Mutation returns the CollectionMutation object of the builder.
func (cu *CollectionUpdate) Mutation() *CollectionMutation {
	return cu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (cu *CollectionUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, cu.sqlSave, cu.mutation, cu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cu *CollectionUpdate) SaveX(ctx context.Context) int {
	affected, err := cu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (cu *CollectionUpdate) Exec(ctx context.Context) error {
	_, err := cu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cu *CollectionUpdate) ExecX(ctx context.Context) {
	if err := cu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cu *CollectionUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(collection.Table, collection.Columns, sqlgraph.NewFieldSpec(collection.FieldID, field.TypeInt))
	if ps := cu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cu.mutation.Name(); ok {
		_spec.SetField(collection.FieldName, field.TypeString, value)
	}
	if value, ok := cu.mutation.IndexType(); ok {
		_spec.SetField(collection.FieldIndexType, field.TypeString, value)
	}
	if value, ok := cu.mutation.DataType(); ok {
		_spec.SetField(collection.FieldDataType, field.TypeString, value)
	}
	if value, ok := cu.mutation.EmbedderType(); ok {
		_spec.SetField(collection.FieldEmbedderType, field.TypeString, value)
	}
	if value, ok := cu.mutation.IndexParams(); ok {
		_spec.SetField(collection.FieldIndexParams, field.TypeJSON, value)
	}
	if value, ok := cu.mutation.EmbedderConfig(); ok {
		_spec.SetField(collection.FieldEmbedderConfig, field.TypeJSON, value)
	}
	if value, ok := cu.mutation.Mappings(); ok {
		_spec.SetField(collection.FieldMappings, field.TypeJSON, value)
	}
	if value, ok := cu.mutation.AppendedMappings(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, collection.FieldMappings, value)
		})
	}
	if n, err = sqlgraph.UpdateNodes(ctx, cu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{collection.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	cu.mutation.done = true
	return n, nil
}

// CollectionUpdateOne is the builder for updating a single Collection entity.
type CollectionUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *CollectionMutation
}

// SetName sets the "name" field.
func (cuo *CollectionUpdateOne) SetName(s string) *CollectionUpdateOne {
	cuo.mutation.SetName(s)
	return cuo
}

// SetIndexType sets the "index_type" field.
func (cuo *CollectionUpdateOne) SetIndexType(s string) *CollectionUpdateOne {
	cuo.mutation.SetIndexType(s)
	return cuo
}

// SetDataType sets the "data_type" field.
func (cuo *CollectionUpdateOne) SetDataType(s string) *CollectionUpdateOne {
	cuo.mutation.SetDataType(s)
	return cuo
}

// SetEmbedderType sets the "embedder_type" field.
func (cuo *CollectionUpdateOne) SetEmbedderType(s string) *CollectionUpdateOne {
	cuo.mutation.SetEmbedderType(s)
	return cuo
}

// SetIndexParams sets the "index_params" field.
func (cuo *CollectionUpdateOne) SetIndexParams(m map[string]interface{}) *CollectionUpdateOne {
	cuo.mutation.SetIndexParams(m)
	return cuo
}

// SetEmbedderConfig sets the "embedder_config" field.
func (cuo *CollectionUpdateOne) SetEmbedderConfig(m map[string]interface{}) *CollectionUpdateOne {
	cuo.mutation.SetEmbedderConfig(m)
	return cuo
}

// SetMappings sets the "mappings" field.
func (cuo *CollectionUpdateOne) SetMappings(s []string) *CollectionUpdateOne {
	cuo.mutation.SetMappings(s)
	return cuo
}

// AppendMappings appends s to the "mappings" field.
func (cuo *CollectionUpdateOne) AppendMappings(s []string) *CollectionUpdateOne {
	cuo.mutation.AppendMappings(s)
	return cuo
}

// Mutation returns the CollectionMutation object of the builder.
func (cuo *CollectionUpdateOne) Mutation() *CollectionMutation {
	return cuo.mutation
}

// Where appends a list predicates to the CollectionUpdate builder.
func (cuo *CollectionUpdateOne) Where(ps ...predicate.Collection) *CollectionUpdateOne {
	cuo.mutation.Where(ps...)
	return cuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (cuo *CollectionUpdateOne) Select(field string, fields ...string) *CollectionUpdateOne {
	cuo.fields = append([]string{field}, fields...)
	return cuo
}

// Save executes the query and returns the updated Collection entity.
func (cuo *CollectionUpdateOne) Save(ctx context.Context) (*Collection, error) {
	return withHooks(ctx, cuo.sqlSave, cuo.mutation, cuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (cuo *CollectionUpdateOne) SaveX(ctx context.Context) *Collection {
	node, err := cuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (cuo *CollectionUpdateOne) Exec(ctx context.Context) error {
	_, err := cuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (cuo *CollectionUpdateOne) ExecX(ctx context.Context) {
	if err := cuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (cuo *CollectionUpdateOne) sqlSave(ctx context.Context) (_node *Collection, err error) {
	_spec := sqlgraph.NewUpdateSpec(collection.Table, collection.Columns, sqlgraph.NewFieldSpec(collection.FieldID, field.TypeInt))
	id, ok := cuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Collection.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := cuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, collection.FieldID)
		for _, f := range fields {
			if !collection.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != collection.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := cuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := cuo.mutation.Name(); ok {
		_spec.SetField(collection.FieldName, field.TypeString, value)
	}
	if value, ok := cuo.mutation.IndexType(); ok {
		_spec.SetField(collection.FieldIndexType, field.TypeString, value)
	}
	if value, ok := cuo.mutation.DataType(); ok {
		_spec.SetField(collection.FieldDataType, field.TypeString, value)
	}
	if value, ok := cuo.mutation.EmbedderType(); ok {
		_spec.SetField(collection.FieldEmbedderType, field.TypeString, value)
	}
	if value, ok := cuo.mutation.IndexParams(); ok {
		_spec.SetField(collection.FieldIndexParams, field.TypeJSON, value)
	}
	if value, ok := cuo.mutation.EmbedderConfig(); ok {
		_spec.SetField(collection.FieldEmbedderConfig, field.TypeJSON, value)
	}
	if value, ok := cuo.mutation.Mappings(); ok {
		_spec.SetField(collection.FieldMappings, field.TypeJSON, value)
	}
	if value, ok := cuo.mutation.AppendedMappings(); ok {
		_spec.AddModifier(func(u *sql.UpdateBuilder) {
			sqljson.Append(u, collection.FieldMappings, value)
		})
	}
	_node = &Collection{config: cuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, cuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{collection.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	cuo.mutation.done = true
	return _node, nil
}
