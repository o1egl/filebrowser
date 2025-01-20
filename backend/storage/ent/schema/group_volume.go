package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

type GroupVolume struct {
	ent.Schema
}

func (GroupVolume) Annotations() []schema.Annotation {
	return []schema.Annotation{
		field.ID("group_id", "volume_id"),
	}
}

func (GroupVolume) Fields() []ent.Field {
	return []ent.Field{
		field.Int("group_id"),
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

func (GroupVolume) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("group", Group.Type).
			Required().
			Unique().
			Field("group_id"),
		edge.To("volume", Volume.Type).
			Required().
			Unique().
			Field("volume_id"),
	}
}
