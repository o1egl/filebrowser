package api

import (
	"sort"

	"github.com/filebrowser/filebrowser/v3/filesystem"
)

func sortResources(resources []filesystem.Info, groupBy GroupBy, sortBy SortBy, orderBy OrderBy) []filesystem.Info {
	var (
		grouperFn    func(filesystem.Info) string
		groupWeights = make(map[string]int)
	)

	// group resources
	switch groupBy {
	case GroupByType:
		grouperFn = func(info filesystem.Info) string {
			return info.Type.String()
		}
		groupWeights[filesystem.TypeDir.String()] = 10
	case GroupByNone:
		fallthrough
	default:
		grouperFn = func(_ filesystem.Info) string { return "" }
	}

	resourceGroups := groupResources(resources, grouperFn)

	// sort resources in a group
	for groupName, resGroup := range resourceGroups {
		sortResourceGroup(resGroup, sortBy, orderBy)
		// add missed groups to the groupWeights
		if _, ok := groupWeights[groupName]; !ok {
			groupWeights[groupName] = -1
		}
	}

	// prepare group names list
	orderedGroupNames := make([]string, 0, len(resourceGroups))
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
	result := make([]filesystem.Info, 0, len(resources))
	for _, name := range orderedGroupNames {
		result = append(result, resourceGroups[name]...)
	}
	return result
}

func groupResources(resources []filesystem.Info, grouper func(filesystem.Info) string) map[string][]filesystem.Info {
	groupedResources := make(map[string][]filesystem.Info)
	for _, resource := range resources {
		groupedResources[grouper(resource)] = append(groupedResources[grouper(resource)], resource)
	}
	return groupedResources
}

func sortResourceGroup(resources []filesystem.Info, sortBy SortBy, orderBy OrderBy) {
	sort.Slice(resources, func(i, j int) bool {
		var result bool

		resA := resources[i]
		resB := resources[j]

		switch sortBy {
		case SortBySize:
			result = resA.Size < resB.Size
		case SortByModified:
			result = resA.ModTime.Unix() < resB.ModTime.Unix()
		case SortByName:
			fallthrough
		default:
			result = resA.Name < resB.Name
		}

		if orderBy == OrderByDesc {
			return !result
		}
		return result
	})
}
