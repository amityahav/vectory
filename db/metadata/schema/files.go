package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

type File struct {
	ent.Schema
}

func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("file_name").Unique(),
	}
}
