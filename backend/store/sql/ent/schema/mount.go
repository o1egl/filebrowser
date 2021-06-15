package schema

import (
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/pkg/errors"
)

// Mount holds the schema definition for the Mount entity.
type Mount struct {
	ent.Schema
}

// Fields of the Mount.
func (Mount) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").Unique(),
		field.String("path").Validate(func(s string) error {
			if !strings.HasPrefix(s, "/") {
				return errors.New("path must start from /")
			}
			return nil
		}),
	}
}

// Edges of the mount.
func (Mount) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("mounts"),
		edge.From("groups", Group.Type).Ref("mounts"),
	}
}
