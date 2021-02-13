package rpc

import (
	"io/fs"
	"os"
	"path/filepath"
)

type NewUserFS func(root, userScope string) fs.FS

func DirFS(root, userScope string) fs.FS {
	return os.DirFS(filepath.Join(root, userScope))
}
