package api

import "net/http"

type FileController struct{}

func NewFileService() *FileController {
	return &FileController{}
}

// List returns list of files
//
//	@Summary		Returns list of files
//	@Description	get files list
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			volume		query		integer				false	"volume id"
//	@Param			path		query		string				true	"path to file"
//	@Param			group_by	query		model.FileGroupBy	false	"group by"	default(none)
//	@Param			sort_by		query		model.SortBy		false	"sort by"
//	@Param			sort_order	query		model.SortOrder		false	"sort order"
//	@Success		200			{array}		model.Group
//	@Failure		default		{object}	model.HTTPError
//	@Router			/v1/files/list [get]
func (f *FileController) List(w http.ResponseWriter, r *http.Request) {
}

// Rename renames a file
//
//	@Summary		renames a file
//	@Description	renames a file
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			volume	query		integer	false	"volume id"
//	@Param			path	query		string	true	"path to file"
//	@Param			name	query		string	true	"new name"
//	@Success		200		{string}	string	"ok"
//	@Failure		default	{object}	model.HTTPError
//	@Router			/v1/files/rename [put]
func (f *FileController) Rename(w http.ResponseWriter, r *http.Request) {}

// Delete renames a file
//
//	@Summary		delete a file
//	@Description	delete a file
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			volume	query		integer	false	"volume id"
//	@Param			path	query		string	true	"path to file"
//	@Success		200		{string}	string	"ok"
//	@Failure		default	{object}	model.HTTPError
//	@Router			/v1/files/delete [delete]
func (f *FileController) Delete(w http.ResponseWriter, r *http.Request) {}

// Copy copies a file
//
//	@Summary		copy a file
//	@Description	copy a file
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.RenameRequest	true	"request"
//	@Success		200		{string}	string				"operation id"
//	@Failure		default	{object}	model.HTTPError
//	@Router			/v1/files/copy [post]
func (f *FileController) Copy(w http.ResponseWriter, r *http.Request) {}

// Move copies a file
//
//	@Summary		copy a file
//	@Description	copy a file
//	@Tags			files
//	@Accept			json
//	@Produce		json
//	@Param			request	body		model.RenameRequest	true	"request"
//	@Success		200		{string}	string				"operation id"
//	@Failure		default	{object}	model.HTTPError
//	@Router			/v1/files/move [post]
func (f *FileController) Move(w http.ResponseWriter, r *http.Request) {}
