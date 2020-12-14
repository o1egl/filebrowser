//go:generate go-enum --sql --marshal --lower --names --file $GOFILE
package api

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v3/backend/filesystem"
	"github.com/filebrowser/filebrowser/v3/backend/log"
	"github.com/filebrowser/filebrowser/v3/backend/rest"
)

type fileController struct {
	root afero.Fs
}

type Resource struct {
	filesystem.Info
	Content string            `json:"content"`
	Items   []filesystem.Info `json:"items"`
}

func (fc *fileController) ListHandler(c *gin.Context) {
	filename := c.Param("path")
	sortBy := c.Param("sort_by")
	sortOrder := c.Param("order")

	user := rest.MustGetUser(c)
	userFs := afero.NewBasePathFs(fc.root, user.Scope)

	info, err := filesystem.Stat(userFs, filename)
	if err != nil {
		switch {
		case errors.Is(err, afero.ErrFileNotFound):
			rest.SendErrorJSON(c, http.StatusBadRequest, err, "resource not found", rest.ErrNotFound)
		default:
			rest.SendErrorJSON(c, http.StatusInternalServerError, err, "can't open requested resource", rest.ErrCodeInternal)
		}
		return
	}

	resource := Resource{Info: info}
	switch info.Type {
	case filesystem.TypeDir:
		infos, err := filesystem.ReadDir(userFs, filename)
		if err != nil {
			rest.SendErrorJSON(c, http.StatusInternalServerError, err, "can't open requested resource", rest.ErrCodeInternal)
			return
		}
		resource.Items = infos
		sortResources(resource.Items, sortBy, sortOrder)
	case filesystem.TypeText:
		b, err := afero.ReadFile(userFs, filename)
		if err != nil {
			rest.SendErrorJSON(c, http.StatusInternalServerError, err, "can't read file", rest.ErrCodeInternal)
			return
		}
		resource.Content = string(b)
	}
	c.JSON(http.StatusOK, resource)
}

func sortResources(resources []filesystem.Info, sortBy, order string) {
	sort.Slice(resources, func(i, j int) bool {
		var result bool

		switch sortBy {
		case "size":
			result = resources[i].Size < resources[j].Size
		case "modified":
			result = resources[i].ModTime.Unix() < resources[j].ModTime.Unix()
		case "name":
			fallthrough
		default:
			result = resources[i].Name < resources[j].Name
		}

		if order == "asc" {
			return !result
		}
		return result
	})
}

func (fc *fileController) ModifyHandler(c *gin.Context) {
	filename := c.Param("path")
	isDir := strings.HasSuffix(filename, "/")
	override := c.Query("override") == "true" || c.Request.Method == http.MethodPut

	user := rest.MustGetUser(c)
	userFs := afero.NewBasePathFs(fc.root, user.Scope)

	switch {
	case c.Request.Method == http.MethodPost && !user.Permissions.CanCreate():
		msg := "no permissions to create resource"
		rest.SendErrorJSON(c, http.StatusForbidden, errors.New(msg), msg, rest.ErrCodeNoPermissions)
		return
	case c.Request.Method == http.MethodPut && !user.Permissions.CanModify():
		msg := "no permissions to update resource"
		rest.SendErrorJSON(c, http.StatusForbidden, errors.New(msg), msg, rest.ErrCodeNoPermissions)
		return
	}

	if err := checkFileExistence(userFs, filename, isDir, override); err != nil {
		rest.SendErrorJSON(c, err.HttpCode(), err, err.Details(), err.ErrCode())
		return
	}

	// handle directory operations
	if isDir {
		if err := userFs.MkdirAll(filename, 0775); err != nil {
			rest.SendErrorJSON(c, http.StatusInternalServerError, err, "failed to create directory", rest.ErrCodeInternal)
			return
		}
		c.Status(http.StatusCreated)
		return
	}

	defer func() {
		_, _ = io.Copy(ioutil.Discard, c.Request.Body)
		_ = c.Request.Body.Close()
	}()

	// handle file operations
	err := func() error {
		dir, _ := filepath.Split(filename)
		if err := userFs.MkdirAll(dir, 0775); err != nil {
			return errors.Wrap(err, "failed to create directory")
		}
		file, err := userFs.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
		if err != nil {
			return errors.Wrap(err, "failed to open file")
		}
		defer file.Close()

		_, err = io.Copy(file, c.Request.Body)
		if err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		if err := userFs.RemoveAll(filename); err != nil {
			log.WithContext(c.Request.Context()).Errorf("failed to clean resource: %v", err)
			return
		}
		rest.SendErrorJSON(c, http.StatusInternalServerError, err, "failed to modify resource", rest.ErrCodeInternal)
		return
	}
}

func checkFileExistence(fs afero.Fs, filename string, isDir, override bool) *rest.HttpError {
	fileInfo, err := filesystem.Stat(fs, filename)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return rest.NewHttpError(err, "failed to get file info", rest.ErrCodeInternal, http.StatusInternalServerError)
	}

	if fileInfo.Type == filesystem.TypeDir {
		err := errors.Errorf("folder %s already exist", filename)
		return rest.NewHttpError(err, err.Error(), rest.ErrFolderExist, http.StatusConflict)
	}

	if isDir || !override {
		err := errors.Errorf("file %s already exist", filename)
		return rest.NewHttpError(err, err.Error(), rest.ErrFileExist, http.StatusConflict)
	}

	return nil
}

/*
ENUM(
copy
move
)
*/
type FileAction int

/*
ENUM(
error
override
rename
)
*/
type OnConflictAction int

func (fc *fileController) MoveHandler(c *gin.Context) {
	filename := c.Param("path")
	action, err := ParseFileAction(c.Query("action"))
	if err != nil {
		rest.NewHttpError(err, err.Error(), rest.ErrBadRequest, http.StatusBadRequest)
		return
	}
	_, err = ParseOnConflictAction(c.Query("on-conflict"))
	if err != nil {
		rest.NewHttpError(err, err.Error(), rest.ErrBadRequest, http.StatusBadRequest)
		return
	}

	user := rest.MustGetUser(c)
	userFs := afero.NewBasePathFs(fc.root, user.Scope)

	switch action {
	case FileActionCopy:
		if !user.Permissions.CanCreate() {
			msg := "no permissions to create resource"
			rest.SendErrorJSON(c, http.StatusForbidden, errors.New(msg), msg, rest.ErrCodeNoPermissions)
		}
	case FileActionMove:
		if !user.Permissions.CanRename() {
			msg := "no permissions to rename resource"
			rest.SendErrorJSON(c, http.StatusForbidden, errors.New(msg), msg, rest.ErrCodeNoPermissions)
		}
	}

	if err := userFs.RemoveAll(filename); err != nil {
		rest.SendErrorJSON(c, http.StatusInternalServerError, err, "failed to delete file", rest.ErrCodeInternal)
		return
	}
}

func (fc *fileController) DeleteHandler(c *gin.Context) {
	filename := c.Param("path")

	user := rest.MustGetUser(c)
	userFs := afero.NewBasePathFs(fc.root, user.Scope)

	if !user.Permissions.CanDelete() {
		msg := "no permissions to delete resource"
		rest.SendErrorJSON(c, http.StatusForbidden, errors.New(msg), msg, rest.ErrCodeNoPermissions)
	}

	if err := userFs.RemoveAll(filename); err != nil {
		rest.SendErrorJSON(c, http.StatusInternalServerError, err, "failed to delete file", rest.ErrCodeInternal)
		return
	}
}
