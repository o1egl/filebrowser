package schema

import (
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/pkg/errors"
)

// Volume holds the schema definition for the Volume entity.
type Volume struct {
	ent.Schema
}

// Fields of the Volume.
func (Volume) Fields() []ent.Field {
	return []ent.Field{
		field.String("label").Unique(),
		field.String("path").Validate(func(s string) error {
			if !strings.HasPrefix(s, "/") {
				return errors.New("path must start from /")
			}
			return nil
		}),
	}
}

// Edges of the volume.
func (Volume) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("users", User.Type).Ref("volumes"),
		edge.From("groups", Group.Type).Ref("volumes"),
	}
}
