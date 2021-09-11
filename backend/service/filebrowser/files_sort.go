package filebrowser

import (
	"sort"

	"github.com/filebrowser/filebrowser/v3/filesystem"
)

func sortFiles(files []File, groupBy GroupBy, sortBy SortBy, orderBy OrderBy) []File {
	var (
		grouperFn    func(file File) string
		groupWeights = make(map[string]int)
	)

	// group files
	switch groupBy {
	case GroupByType:
		grouperFn = func(info File) string {
			return info.Type.String()
		}
		groupWeights[filesystem.TypeDir.String()] = 10
	case GroupByNone:
		fallthrough
	default:
		grouperFn = func(_ File) string { return "" }
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
	result := make([]File, 0, len(files))
	for _, name := range orderedGroupNames {
		result = append(result, fileGroups[name]...)
	}
	return result
}

func groupFiles(files []File, grouper func(File) string) map[string][]File {
	groupedFiles := make(map[string][]File)
	for _, file := range files {
		groupedFiles[grouper(file)] = append(groupedFiles[grouper(file)], file)
	}
	return groupedFiles
}

func sortFileGroup(files []File, sortBy SortBy, orderBy OrderBy) {
	sort.Slice(files, func(i, j int) bool {
		var result bool

		fileA := files[i]
		fileB := files[j]

		switch sortBy {
		case SortBySize:
			result = fileA.Size < fileB.Size
		case SortByModified:
			result = fileA.ModTime.Unix() < fileB.ModTime.Unix()
		case SortByName:
			fallthrough
		default:
			result = fileA.Name < fileB.Name
		}

		if orderBy == OrderByDesc {
			return !result
		}
		return result
	})
}
