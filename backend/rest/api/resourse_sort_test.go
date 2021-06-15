package api

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/filebrowser/filebrowser/v3/filesystem"
)

func Test_sortResources(t *testing.T) {
	modTime := time.Date(2021, 1, 2, 15, 32, 14, 0, time.UTC)
	resources := []filesystem.Info{
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
		resources []filesystem.Info
		groupBy   GroupBy
		sortBy    SortBy
		orderBy   OrderBy
		want      []filesystem.Info
	}{
		"group by: type, sort by: name, order: asc": {
			resources: resources,
			groupBy:   GroupByType,
			sortBy:    SortByName,
			orderBy:   OrderByAsc,
			want: []filesystem.Info{
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
			resources: resources,
			groupBy:   GroupByType,
			sortBy:    SortByName,
			orderBy:   OrderByDesc,
			want: []filesystem.Info{
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
			resources: resources,
			groupBy:   GroupByType,
			sortBy:    SortByName,
			orderBy:   OrderByAsc,
			want: []filesystem.Info{
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
			resources: resources,
			groupBy:   GroupByNone,
			sortBy:    SortBySize,
			orderBy:   OrderByDesc,
			want: []filesystem.Info{
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
			got := sortResources(tt.resources, tt.groupBy, tt.sortBy, tt.orderBy)
			fmt.Println(got)
			assert.Equal(t, tt.want, got)
		})
	}
}
