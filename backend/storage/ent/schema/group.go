package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

type Group struct {
	ent.Schema
}

func (Group) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

func (Group) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			NotEmpty(),
	}
}

func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).
			Ref("groups"),
		edge.To("volumes", Volume.Type).
			Through("group_volume", GroupVolume.Type),
	}
}
