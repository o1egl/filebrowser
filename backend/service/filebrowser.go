//go:generate go-enum --sql --marshal --nocase --names --file $GOFILE
package service

import (
	"context"
	"io"
	"os"
	"time"

	"github.com/filebrowser/filebrowser/v3/auth"
	"github.com/filebrowser/filebrowser/v3/filesystem"
)

const (
	HomeVolumeLabel       = "Home"
	HomeVolumeID    int64 = 0

	DefaultGroupBy = GroupByType
	DefaultSortBy  = SortByName
	DefaultOrderBy = OrderByAsc
	NoLimit        = -1
)

type FileBrowser interface {
	List(ctx context.Context, user auth.User, params ListParams) (*FileWithChildren, error)
	Create(ctx context.Context, user auth.User, volume int64, filename string, content io.Reader) error
	Update(ctx context.Context, user auth.User, volume int64, filename string, content io.Reader) error
	Delete(ctx context.Context, user auth.User, volume int64, filename string) error
}

type ListParams struct {
	Volume   int64
	Filename string
	GroupBy  GroupBy
	SortBy   SortBy
	OrderBy  OrderBy
	Offset   int
	Limit    int
}

type FileWithChildren struct {
	File
	Children []File       `json:"children,omitempty"`
	Meta     FileMetaData `json:"meta"`
}

type File struct {
	Path      string          `json:"path"`
	Name      string          `json:"name"`
	Size      int64           `json:"size"`
	Extension string          `json:"extension,omitempty"`
	ModTime   time.Time       `json:"modified"`
	Mode      os.FileMode     `json:"mode"`
	Type      filesystem.Type `json:"type"`
	IsSymlink bool            `json:"is_symlink"`
	IsDir     bool            `json:"is_dir"`
}

type FileMetaData struct {
	FilesCount int `json:"files_count"`
	DirsCount  int `json:"dirs_count"`
	TotalCount int `json:"total_count"`
}

/*
ENUM(
none
type
)
*/
type GroupBy int

/*
ENUM(
name
size
modified
)
*/
type SortBy int

/*
ENUM(
asc
desc
)
*/
type OrderBy int
