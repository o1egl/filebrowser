package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type UserVolume struct {
	ent.Schema
}

func (UserVolume) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("user_id", "volume_id"),
	}
}

func (UserVolume) Fields() []ent.Field {
	return []ent.Field{
		field.Int("user_id"),
		field.Int("volume_id"),
		field.Bool("read").
			Default(true),
		field.Bool("write").
			Default(true),
		field.Bool("delete").
			Default(true),
		field.Bool("share").
			Default(true),
	}
}

func (UserVolume) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Unique().
			Required().
			Field("user_id"),
		edge.To("volume", Volume.Type).
			Unique().
			Required().
			Field("volume_id"),
	}
}
