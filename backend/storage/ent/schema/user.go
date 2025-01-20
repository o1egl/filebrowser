package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.Time{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("username").
			Unique().
			NotEmpty(),
		field.String("password").
			Sensitive().
			NotEmpty(),
		field.String("name").
			NotEmpty(),
		field.String("locale").
			Default("en"),
		field.Int("home_volume_id").
			Optional(),
		field.Bool("blocked").
			Default(false),
		field.Bool("admin").
			Default(false),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("volumes", Volume.Type).
			Through("user_volume", UserVolume.Type),
		edge.To("groups", Group.Type),
	}
}
