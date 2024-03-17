package api

import "net/http"

type OperationsController struct{}

// Get returns operation status
//
//	@Summary		returns status of an operation
//	@Description	returns status of an operation
//	@Tags			operations
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"operation id"
//	@Success		200	{object}	model.OperationStatus
//	@Failure		400	{object}	model.HTTPError
//	@Failure		404	{object}	model.HTTPError
//	@Failure		500	{object}	model.HTTPError
//	@Router			/v1/operations/{id} [get]
func (f *OperationsController) Get(w http.ResponseWriter, r *http.Request) {}
