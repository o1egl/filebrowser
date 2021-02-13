package filesystem

import "mime"

func init() {
	fileExtMap := map[string]string{
		".gitignore": "text/plain",
		".md":        "text/plain",
		".yml":       "text/yaml",
		".yaml":      "text/yaml",
		".json":      "text/json",
		".js":        "text/js",
	}
	for ext, fileType := range fileExtMap {
		_ = mime.AddExtensionType(ext, fileType)
	}
}
