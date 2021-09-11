package api

import (
	"encoding/json"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/filebrowser/filebrowser/v3/domain"
	"github.com/gin-gonic/gin"

	"github.com/filebrowser/filebrowser/v3/assets"
	"github.com/filebrowser/filebrowser/v3/rest"
)

// staticController provides router for all requests with no required auth
type staticController struct {
	BasePath  string
	Version   domain.Version
	Anonymous bool
}

func newStaticController(basePath string, version domain.Version, anonymous bool) *staticController {
	return &staticController{BasePath: basePath, Version: version, Anonymous: anonymous}
}

func (h *staticController) indexHandler(c *gin.Context) {
	data := map[string]interface{}{
		"Name":       "File Browser",
		"BaseURL":    h.BasePath,
		"Version":    h.Version,
		"StaticURL":  path.Join(h.BasePath, "/static"),
		"Signup":     true,
		"NoAuth":     h.Anonymous,
		"AuthMethod": "json",
		"LoginPage":  true,
		"CSS":        false,
		"ReCaptcha":  false,
		"Theme":      "",
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		rest.SendErrorJSON(c, http.StatusInternalServerError, err, "data encoding error", rest.ErrCodeInternal)
		return
	}
	data["Json"] = string(jsonData)

	c.HTML(http.StatusOK, "web/dist/index.html", data)
}

// staticHandler returns static assets
func (h *staticController) staticHandler(c *gin.Context) {
	fPath := filepath.Join("web/dist", c.Param("path"))
	if !strings.HasSuffix(fPath, ".js") {
		c.FileFromFS(fPath, http.FS(assets.FS()))
		return
	}

	c.Header("Content-Mode", "text/javascript")

	data := map[string]interface{}{
		"StaticURL": path.Join(h.BasePath, "/static"),
	}
	c.HTML(http.StatusOK, fPath, data)
}
