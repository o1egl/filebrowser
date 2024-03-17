//go:generate go-enum --marshal
package model

type File struct {
	Name         string       `json:"name"`
	Size         int64        `json:"size"`
	Capabilities Capabilities `json:"capabilities"`
}

type Capabilities struct {
	Read   bool `json:"read"`
	Write  bool `json:"write"`
	Delete bool `json:"delete"`
	Rename bool `json:"rename"`
	Share  bool `json:"share"`
}

type RenameRequest struct {
	Src        FileLocation `json:"src"`
	Dst        FileLocation `json:"dst"`
	OnConflict OnConflict   `json:"on_conflict"`
}

type FileLocation struct {
	Volume int64  `json:"volume"`
	Path   string `json:"path"`
}

type Group struct {
	Name  string `json:"name"`
	Files []File `json:"files"`
}

// ENUM(none, kind, modified, size)
type FileGroupBy string

// ENUM(name, size, modified)
type SortBy string

// ENUM(asc, desc)
type SortOrder string

// ENUM(skip, override, rename)
type OnConflict string
