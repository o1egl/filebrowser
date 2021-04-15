package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").Unique().DefaultFunc(uuid.NewString),
		field.String("provider"),
		field.String("username"),
		field.String("password").Optional(),
		field.String("name").Optional(),
		field.String("scope").Default("/"),
		field.String("locale"),
		field.Bool("lock_password"),
		field.Bool("blocked"),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		// unique index.
		index.Fields("provider", "username").Unique(),
	}
}
