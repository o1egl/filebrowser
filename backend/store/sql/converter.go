//go:generate goverter -packageName conv -output ./conv/converter_gen.go github.com/filebrowser/filebrowser/v3/store/sql
package sql

import (
	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent"
)

// goverter:converter
type UserConverter interface {
	Convert(source *ent.User) *store.User
	ConvertSlice(source []*ent.User) []*store.User
}
