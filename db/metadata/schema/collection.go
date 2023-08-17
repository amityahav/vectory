package schema

import (
	"entgo.io/ent"
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
		field.String("embedder_type"),
		field.JSON("index_params", map[string]interface{}{}),
		field.JSON("embedder_config", map[string]interface{}{}),
		field.Strings("mappings"),
	}
}
