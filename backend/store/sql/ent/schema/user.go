package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique(),
		field.String("provider"),
		field.String("username").Unique(),
		field.String("password").Optional(),
		field.String("home"),
		field.String("name"),
		field.String("locale"),
		field.Bool("lock_password"),
		field.Bool("blocked"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("provider"),
	}
}

// Edges of the user.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("volumes", Volume.Type),
		edge.To("groups", Group.Type),
	}
}
