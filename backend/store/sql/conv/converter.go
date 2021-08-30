//go:generate goverter github.com/filebrowser/filebrowser/v3/store/sql/conv
package conv

import (
	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent"
)

// goverter:converter
type UserConverter interface {
	Convert(source *ent.User) *store.User
	ConvertSlice(source []*ent.User) []*store.User
}

// goverter:converter
// goverter:extend IntToInt64
type VolumeConverter interface {
	Convert(source *ent.Volume) *store.Volume
	ConvertSlice(source []*ent.Volume) []*store.Volume
}

func IntToInt64(i int) int64 {
	return int64(i)
}
