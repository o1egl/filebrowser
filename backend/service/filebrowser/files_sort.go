package filebrowser

import (
	"sort"

	"github.com/filebrowser/filebrowser/v3/filesystem"
	"github.com/filebrowser/filebrowser/v3/service"
)

func sortFiles(files []service.File, groupBy service.GroupBy, sortBy service.SortBy, orderBy service.OrderBy) []service.File {
	var (
		grouperFn    func(file service.File) string
		groupWeights = make(map[string]int)
	)

	// group files
	switch groupBy {
	case service.GroupByType:
		grouperFn = func(info service.File) string {
			return info.Type.String()
		}
		groupWeights[filesystem.TypeDir.String()] = 10
	case service.GroupByNone:
		fallthrough
	default:
		grouperFn = func(_ service.File) string { return "" }
	}

	fileGroups := groupFiles(files, grouperFn)

	// sort files in a group
	for groupName, fileGroup := range fileGroups {
		sortFileGroup(fileGroup, sortBy, orderBy)
		// add missed groups to the groupWeights
		if _, ok := groupWeights[groupName]; !ok {
			groupWeights[groupName] = -1
		}
	}

	// prepare group names list
	orderedGroupNames := make([]string, 0, len(fileGroups))
	for name := range groupWeights {
		orderedGroupNames = append(orderedGroupNames, name)
	}
	sort.Slice(orderedGroupNames, func(i, j int) bool {
		nameI := orderedGroupNames[i]
		nameJ := orderedGroupNames[j]
		// order by name is added for keeping consistent order of the groups with the same weights
		if groupWeights[nameI] == groupWeights[nameJ] {
			return nameI < nameJ
		}
		// group with the higher weight should be placed first
		return groupWeights[nameI] > groupWeights[nameJ]
	})

	// merge groups
	result := make([]service.File, 0, len(files))
	for _, name := range orderedGroupNames {
		result = append(result, fileGroups[name]...)
	}
	return result
}

func groupFiles(files []service.File, grouper func(service.File) string) map[string][]service.File {
	groupedFiles := make(map[string][]service.File)
	for _, file := range files {
		groupedFiles[grouper(file)] = append(groupedFiles[grouper(file)], file)
	}
	return groupedFiles
}

func sortFileGroup(files []service.File, sortBy service.SortBy, orderBy service.OrderBy) {
	sort.Slice(files, func(i, j int) bool {
		var result bool

		fileA := files[i]
		fileB := files[j]

		switch sortBy {
		case service.SortBySize:
			result = fileA.Size < fileB.Size
		case service.SortByModified:
			result = fileA.ModTime.Unix() < fileB.ModTime.Unix()
		case service.SortByName:
			fallthrough
		default:
			result = fileA.Name < fileB.Name
		}

		if orderBy == service.OrderByDesc {
			return !result
		}
		return result
	})
}
