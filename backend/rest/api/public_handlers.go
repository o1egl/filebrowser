package api

import (
	"encoding/json"
	"net/http"
	"path"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"

	"github.com/filebrowser/filebrowser/v3/backend/rest"
)

// publicHandlers provides router for all requests with no required auth
type publicHandlers struct {
	BasePath  string
	Revision  string
	Anonymous bool
}

func (h *publicHandlers) indexHandler(c *gin.Context) {
	data := map[string]interface{}{
		"Name":            "File Browser",
		"DisableExternal": true,
		"BaseURL":         h.BasePath,
		"Version":         h.Revision,
		"StaticURL":       path.Join(h.BasePath, "/static"),
		"Signup":          true,
		"NoAuth":          h.Anonymous,
		"AuthMethod":      "json",
		"LoginPage":       true,
		"CSS":             false,
		"ReCaptcha":       false,
		"Theme":           "",
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		rest.SendErrorJSON(c, http.StatusInternalServerError, err, "data encoding error", rest.ErrCodeInternal)
		return
	}
	data["Json"] = string(jsonData)

	c.HTML(http.StatusOK, "/frontend/dist/index.html", data)
}

// staticHandler returns static assets
func (h *publicHandlers) staticHandler(c *gin.Context) {
	fPath := c.Param("path")
	if !strings.HasSuffix(fPath, ".js") {
		c.FileFromFS(fPath, pkger.Dir("/frontend/dist"))
		return
	}

	c.Header("Content-Type", "text/javascript")

	data := map[string]interface{}{
		"StaticURL": path.Join(h.BasePath, "/static"),
	}
	c.HTML(http.StatusOK, path.Join("/frontend/dist", fPath), data)
}
