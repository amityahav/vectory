package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type Collection struct {
	ent.Schema
}

func (Collection) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("index_type"),
		field.String("data_type"),
		field.String("embedder"),
		field.JSON("index_params", map[string]interface{}{}),
	}
}

func (Collection) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("files", File.Type),
	}
}
