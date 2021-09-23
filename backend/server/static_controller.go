package server

import (
	"encoding/json"
	"net/http"
	"path"
	"path/filepath"
	"strings"

	"github.com/filebrowser/filebrowser/v3/config"
	"github.com/filebrowser/filebrowser/v3/domain"
	"github.com/gin-gonic/gin"

	"github.com/filebrowser/filebrowser/v3/assets"
	"github.com/filebrowser/filebrowser/v3/rest"
)

// staticController provides router for all requests with no required auth
type staticController struct {
	cfg     *config.Config
	version domain.Version
}

func newStaticController(cfg *config.Config, version domain.Version) *staticController {
	return &staticController{cfg: cfg, version: version}
}

func (h *staticController) indexHandler(c *gin.Context) {
	data := map[string]interface{}{
		"Name":       "File Browser",
		"BaseURL":    h.cfg.Server.BasePath(),
		"Version":    h.version,
		"StaticURL":  path.Join(h.cfg.Server.BasePath(), "/static"),
		"Signup":     true,
		"NoAuth":     h.cfg.Auth.Anonymous.Enabled,
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
		"StaticURL": path.Join(h.cfg.Server.BasePath(), "/static"),
	}
	c.HTML(http.StatusOK, fPath, data)
}
