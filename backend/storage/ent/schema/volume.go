package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"regexp"
)

// Volume holds the schema definition for the Volume entity.
type Volume struct {
	ent.Schema
}

func (Volume) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the Volume.
func (Volume) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").
			Unique().
			NotEmpty(),
		field.String("path").
			NotEmpty().
			Match(regexp.MustCompile("^/.*$")),
		field.Int("backend_id"),
	}
}

// Edges of the Volume.
func (Volume) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).
			Ref("volumes").
			Through("user_volume", UserVolume.Type),
		edge.From("groups", Group.Type).
			Ref("volumes").
			Through("group_volume", GroupVolume.Type),
		edge.To("backend", Backend.Type).
			Unique().
			Field("backend_id").
			Required(),
	}
}
