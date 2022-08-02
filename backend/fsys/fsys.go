//go:generate ${TOOLS_BIN}/go-enum --marshal --nocase --names --file $GOFILE
package fsys

import (
	"errors"
	"io"
	"os"
	"time"
)

var (
	ErrNotExist        = os.ErrNotExist
	ErrExist           = os.ErrExist
	ErrNotARegularFile = errors.New("not a regular file")
)

/*
ENUM(
regular
dir
special
)
*/
type FileType int

type File struct {
	Name      string
	Size      int64
	Type      FileType
	ModTime   time.Time
	IsSymLink bool
}

type FS interface {
	MkDir(fPath string) error
	Remove(fPath string) error
	Stat(fPath string) (File, error)
	List(fPath string) ([]File, error)
	Read(fPath string) (io.ReadCloser, error)
	Write(fPath string, reader io.Reader) (written int64, err error)
	Move(src, dst string) error
}
