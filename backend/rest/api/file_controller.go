//go:generate go-enum --sql --marshal --nocase --names --file $GOFILE
package api

import (
	"net/http"
	"strconv"

	"github.com/filebrowser/filebrowser/v3/mathx"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/filebrowser/filebrowser/v3/filesystem"
	"github.com/filebrowser/filebrowser/v3/rest"
)

const (
	defaultGroupBy = GroupByType
	defaultSortBy  = SortByName
	defaultOrderBy = OrderByAsc
)

type fileController struct {
	rootFS afero.Fs
}

type Resource struct {
	filesystem.Info
	Children []filesystem.Info `json:"children"`
	Meta     FileMeta          `json:"meta"`
}

type FileMeta struct {
	FilesCount int `json:"files_count"`
	DirsCount  int `json:"dirs_count"`
	TotalCount int `json:"total_count"`
}

/*
ENUM(
none
type
)
*/
type GroupBy int

/*
ENUM(
name
size
modified
)
*/
type SortBy int

/*
ENUM(
asc
desc
)
*/
type OrderBy int

type ListHandlerParams struct {
	Filename string
	GroupBy  GroupBy
	SortBy   SortBy
	OrderBy  OrderBy
	Offset   int
	Limit    int
}

func (fc *fileController) ListHandler(c *gin.Context) {
	params, err := parseListHandlerParams(c)
	if err != nil {
		rest.SendBadRequestError(c, err, "failed to parse input params")
		return
	}

	user := rest.MustGetUser(c)
	userFs := afero.NewBasePathFs(fc.rootFS, user.Home)

	info, err := filesystem.Stat(userFs, params.Filename)
	if err != nil {
		switch {
		case errors.Is(err, afero.ErrFileNotFound):
			rest.SendNotFoundError(c, err, "resource not found")
		default:
			rest.SendInternalError(c, err, "can't read requested resource")
		}
		return
	}

	response := Resource{
		Info: info,
	}
	if info.Type == filesystem.TypeDir {
		children, err := filesystem.ReadDir(userFs, params.Filename)
		if err != nil {
			rest.SendErrorJSON(c, http.StatusInternalServerError, err, "can't open requested resource", rest.ErrCodeInternal)
			return
		}
		response.Children = sortResources(children, params.GroupBy, params.SortBy, params.OrderBy)
	}

	// Add metadata
	for _, res := range response.Children {
		if res.IsDir {
			response.Meta.DirsCount++
		} else {
			response.Meta.FilesCount++
		}
		response.Meta.TotalCount++
	}

	// apply offset/limit
	offset := mathx.MinInt(params.Offset, len(response.Children))
	limit := len(response.Children)
	if params.Limit > 0 {
		limit = mathx.MinInt(offset+params.Limit, len(response.Children))
	}
	response.Children = response.Children[offset:limit]

	c.JSON(http.StatusOK, response)
}

func parseListHandlerParams(c *gin.Context) (*ListHandlerParams, error) {
	filename := c.Param("path")
	groupByInput := c.Query("group_by")
	sortByInput := c.Query("sort_by")
	orderByInput := c.Query("order_by")
	offsetInput := c.Query("offset")
	limitInput := c.Query("limit")

	params := &ListHandlerParams{
		Filename: filename,
		GroupBy:  defaultGroupBy,
		SortBy:   defaultSortBy,
		OrderBy:  defaultOrderBy,
		Offset:   0,
		Limit:    -1,
	}

	var err error
	if groupByInput != "" {
		if params.GroupBy, err = ParseGroupBy(groupByInput); err != nil {
			return nil, errors.Wrap(err, "incorrect group_by param")
		}
	}
	if sortByInput != "" {
		if params.SortBy, err = ParseSortBy(sortByInput); err != nil {
			return nil, errors.Wrap(err, "incorrect sort_by param")
		}
	}
	if orderByInput != "" {
		if params.OrderBy, err = ParseOrderBy(orderByInput); err != nil {
			return nil, errors.Wrap(err, "incorrect order_by param")
		}
	}
	if offsetInput != "" {
		if params.Offset, err = strconv.Atoi(offsetInput); err != nil {
			return nil, errors.Wrap(err, "incorrect offset param")
		}
		if params.Offset < 0 {
			return nil, errors.New("offset must be negative")
		}
	}
	if limitInput != "" {
		if params.Limit, err = strconv.Atoi(limitInput); err != nil {
			return nil, errors.Wrap(err, "incorrect limit param")
		}
		if params.Limit == 0 {
			return nil, errors.New("offset must be greater than 0")
		}
	}

	return params, nil
}

/*func (fc *fileController) ModifyHandler(c *gin.Context) {
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
}*/

/*func checkFileExistence(fs afero.Fs, filename string, isDir, override bool) *rest.HttpError {
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
}*/

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

/*func (fc *fileController) MoveHandler(c *gin.Context) {
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
}*/
