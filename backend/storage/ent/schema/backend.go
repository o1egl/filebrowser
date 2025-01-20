package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type Backend struct {
	ent.Schema
}

func (Backend) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Backend) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			NotEmpty(),
		field.String("type").
			NotEmpty(),
		field.String("config").
			Optional(),
	}
}

func (Backend) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("volumes", Volume.Type).
			Ref("backend"),
	}
}
