package filebrowser

import (
	"fmt"
	"testing"
	"time"

	"github.com/filebrowser/filebrowser/v3/service"
	"github.com/stretchr/testify/assert"

	"github.com/filebrowser/filebrowser/v3/filesystem"
)

func Test_sortResources(t *testing.T) {
	modTime := time.Date(2021, 1, 2, 15, 32, 14, 0, time.UTC)
	files := []service.File{
		{
			Path:      "/file1",
			Name:      "file1",
			Size:      10,
			Extension: "",
			ModTime:   modTime,
			Mode:      22,
			Type:      filesystem.TypeFile,
			IsSymlink: false,
			IsDir:     false,
		},
		{
			Path:      "/file2",
			Name:      "file2",
			Size:      20,
			Extension: "",
			ModTime:   modTime,
			Mode:      22,
			Type:      filesystem.TypeFile,
			IsSymlink: false,
			IsDir:     false,
		},
		{
			Path:      "/dir1",
			Name:      "dir1",
			Size:      11,
			Extension: "",
			ModTime:   modTime,
			Mode:      22,
			Type:      filesystem.TypeDir,
			IsSymlink: false,
			IsDir:     true,
		},
		{
			Path:      "/dir2",
			Name:      "dir2",
			Size:      22,
			Extension: "",
			ModTime:   modTime,
			Mode:      22,
			Type:      filesystem.TypeDir,
			IsSymlink: false,
			IsDir:     true,
		},
	}
	testCases := map[string]struct {
		resources []service.File
		groupBy   service.GroupBy
		sortBy    service.SortBy
		orderBy   service.OrderBy
		want      []service.File
	}{
		"group by: type, sort by: name, order: asc": {
			resources: files,
			groupBy:   service.GroupByType,
			sortBy:    service.SortByName,
			orderBy:   service.OrderByAsc,
			want: []service.File{
				{
					Path:      "/dir1",
					Name:      "dir1",
					Size:      11,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeDir,
					IsSymlink: false,
					IsDir:     true,
				},
				{
					Path:      "/dir2",
					Name:      "dir2",
					Size:      22,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeDir,
					IsSymlink: false,
					IsDir:     true,
				},
				{
					Path:      "/file1",
					Name:      "file1",
					Size:      10,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeFile,
					IsSymlink: false,
					IsDir:     false,
				},
				{
					Path:      "/file2",
					Name:      "file2",
					Size:      20,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeFile,
					IsSymlink: false,
					IsDir:     false,
				},
			},
		},
		"group by: type, sort by: name, order: desc": {
			resources: files,
			groupBy:   service.GroupByType,
			sortBy:    service.SortByName,
			orderBy:   service.OrderByDesc,
			want: []service.File{
				{
					Path:      "/dir2",
					Name:      "dir2",
					Size:      22,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeDir,
					IsSymlink: false,
					IsDir:     true,
				},
				{
					Path:      "/dir1",
					Name:      "dir1",
					Size:      11,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeDir,
					IsSymlink: false,
					IsDir:     true,
				},
				{
					Path:      "/file2",
					Name:      "file2",
					Size:      20,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeFile,
					IsSymlink: false,
					IsDir:     false,
				},
				{
					Path:      "/file1",
					Name:      "file1",
					Size:      10,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeFile,
					IsSymlink: false,
					IsDir:     false,
				},
			},
		},
		"group by: none, sort by: name, order: asc": {
			resources: files,
			groupBy:   service.GroupByType,
			sortBy:    service.SortByName,
			orderBy:   service.OrderByAsc,
			want: []service.File{
				{
					Path:      "/dir1",
					Name:      "dir1",
					Size:      11,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeDir,
					IsSymlink: false,
					IsDir:     true,
				},
				{
					Path:      "/dir2",
					Name:      "dir2",
					Size:      22,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeDir,
					IsSymlink: false,
					IsDir:     true,
				},
				{
					Path:      "/file1",
					Name:      "file1",
					Size:      10,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeFile,
					IsSymlink: false,
					IsDir:     false,
				},
				{
					Path:      "/file2",
					Name:      "file2",
					Size:      20,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeFile,
					IsSymlink: false,
					IsDir:     false,
				},
			},
		},
		"group by: none, sort by: size, order: asc": {
			resources: files,
			groupBy:   service.GroupByNone,
			sortBy:    service.SortBySize,
			orderBy:   service.OrderByDesc,
			want: []service.File{
				{
					Path:      "/dir2",
					Name:      "dir2",
					Size:      22,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeDir,
					IsSymlink: false,
					IsDir:     true,
				},
				{
					Path:      "/file2",
					Name:      "file2",
					Size:      20,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeFile,
					IsSymlink: false,
					IsDir:     false,
				},
				{
					Path:      "/dir1",
					Name:      "dir1",
					Size:      11,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeDir,
					IsSymlink: false,
					IsDir:     true,
				},
				{
					Path:      "/file1",
					Name:      "file1",
					Size:      10,
					Extension: "",
					ModTime:   modTime,
					Mode:      22,
					Type:      filesystem.TypeFile,
					IsSymlink: false,
					IsDir:     false,
				},
			},
		},
	}
	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			got := sortFiles(tt.resources, tt.groupBy, tt.sortBy, tt.orderBy)
			fmt.Println(got)
			assert.Equal(t, tt.want, got)
		})
	}
}
