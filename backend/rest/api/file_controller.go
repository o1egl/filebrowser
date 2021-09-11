package api

import (
	"net/http"
	"strconv"

	"github.com/filebrowser/filebrowser/v3/hash"
	"github.com/filebrowser/filebrowser/v3/rest"
	"github.com/filebrowser/filebrowser/v3/service/filebrowser"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	_ "github.com/speps/go-hashids/v2"
)

const homeVolumeID = "home"

type fileController struct {
	fileBrowserSvc filebrowser.Service
	hasher         hash.Hasher
}

func newFileController(fileBrowserSvc filebrowser.Service, hasher hash.Hasher) *fileController {
	return &fileController{fileBrowserSvc: fileBrowserSvc, hasher: hasher}
}

func (fc *fileController) ListHandler(c *gin.Context) {
	ctx := c.Request.Context()
	params, err := fc.parseListHandlerParams(c)
	if err != nil {
		rest.SendBadRequestError(c, err, "failed to parse input params")
		return
	}

	user := rest.MustGetUser(c)

	list, err := fc.fileBrowserSvc.List(ctx, user, *params)
	if err != nil {
		rest.SendServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, list)
}

func (fc *fileController) parseListHandlerParams(c *gin.Context) (*filebrowser.ListParams, error) {
	filename := c.Param("path")
	groupByInput := c.Query("group_by")
	sortByInput := c.Query("sort_by")
	orderByInput := c.Query("order_by")
	offsetInput := c.Query("offset")
	limitInput := c.Query("limit")

	params := &filebrowser.ListParams{
		Volume:   filebrowser.HomeVolumeID,
		Filename: filename,
		GroupBy:  filebrowser.DefaultGroupBy,
		SortBy:   filebrowser.DefaultSortBy,
		OrderBy:  filebrowser.DefaultOrderBy,
		Offset:   0,
		Limit:    filebrowser.NoLimit,
	}

	var err error
	params.Volume, err = fc.parseVolumeFromPath(c)
	if err != nil {
		return nil, err
	}
	if groupByInput != "" {
		if params.GroupBy, err = filebrowser.ParseGroupBy(groupByInput); err != nil {
			return nil, errors.Wrap(err, "incorrect group_by param")
		}
	}
	if sortByInput != "" {
		if params.SortBy, err = filebrowser.ParseSortBy(sortByInput); err != nil {
			return nil, errors.Wrap(err, "incorrect sort_by param")
		}
	}
	if orderByInput != "" {
		if params.OrderBy, err = filebrowser.ParseOrderBy(orderByInput); err != nil {
			return nil, errors.Wrap(err, "incorrect order_by param")
		}
	}
	if offsetInput != "" {
		if params.Offset, err = strconv.Atoi(offsetInput); err != nil {
			return nil, errors.Wrap(err, "incorrect offset param")
		}
		if params.Offset < 0 {
			return nil, errors.New("offset must be positive")
		}
	}
	if limitInput != "" {
		if params.Limit, err = strconv.Atoi(limitInput); err != nil {
			return nil, errors.Wrap(err, "incorrect limit param")
		}
		if params.Limit == 0 {
			return nil, errors.New("limit must be greater than 0")
		}
	}

	return params, nil
}

func (fc *fileController) DeleteHandler(c *gin.Context) {
	ctx := c.Request.Context()
	user := rest.MustGetUser(c)

	filename := c.Param("path")
	volumeID, err := fc.parseVolumeFromPath(c)
	if err != nil {
		rest.SendBadRequestError(c, err, "failed to parse input params")
		return
	}

	err = fc.fileBrowserSvc.Delete(ctx, user, volumeID, filename)
	if err != nil {
		rest.SendServiceError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

func (fc *fileController) parseVolumeFromPath(c *gin.Context) (int64, error) {
	volumeInput := c.Param("volume")
	switch volumeInput {
	case homeVolumeID:
		return filebrowser.HomeVolumeID, nil
	default:
		volumeID, err := fc.hasher.DecodeInt64(volumeInput)
		return volumeID, errors.Wrap(err, "incorrect volume id")
	}
}

/*func (fc *fileController) listHandler(c *gin.Context, fSys afero.Fs, params *ListHandlerParams) {
	info, err := filesystem.Stat(fSys, params.Filename)
	if err != nil {
		switch {
		case errors.Is(err, afero.ErrFileNotFound):
			rest.SendNotFoundError(c, err, "resource not found")
		default:
			rest.SendInternalError(c, err, "can't read requested resource")
		}
		return
	}

	response := FileListResponse{
		Info: info,
	}
	if info.Mode == filesystem.TypeDir {
		children, err := filesystem.ReadDir(fSys, params.Filename)
		if err != nil {
			rest.SendErrorJSON(c, http.StatusInternalServerError, err, "can't open requested resource", rest.ErrCodeInternal)
			return
		}
		response.Children = sortResources(children, params.GroupBy, params.SortBy, params.OrderBy)
	}

	// Add metadata
	for _, res := range response.Children {
		if res.IsDir {
			response.Metadata.DirsCount++
		} else {
			response.Metadata.FilesCount++
		}
		response.Metadata.TotalCount++
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



func (fc *fileController) DeleteHandler(c *gin.Context) {
	filename := c.Param("path")

	user := rest.MustGetUser(c)
	userFs := afero.NewBasePathFs(fc.rootFS, user.Home)
	if err := userFs.RemoveAll(filename); err != nil {
		rest.SendInternalError(c, err, "failed to delete file")
		return
	}
	c.Status(http.StatusOK)
}*/

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

	if fileInfo.Mode == filesystem.TypeDir {
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
