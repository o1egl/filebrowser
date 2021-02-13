package assets

import (
	"embed"
	"io/fs"
)

//go:embed web/*
var fsys embed.FS

func FS() fs.FS {
	return fsys
}
