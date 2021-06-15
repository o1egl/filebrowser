package rest

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/filebrowser/filebrowser/v3/log"
)

// All error codes for UI mapping and translation
const (
	ErrCodeInternal      = 0 // any internal error
	ErrCodeDecode        = 1 // failed to unmarshal incoming request
	ErrCodeUnauthorized  = 2 // rejected by auth
	ErrCodeNoPermissions = 3 // no permissions
	ErrCodeUserBlocked   = 4 // user blocked
	ErrCodeReadOnly      = 5 // write failed on read only
	ErrBadRequest        = 6 // bad request
	ErrNotFound          = 7 // not found
	ErrFileExist         = 8 // file with provided name already exist exist
	ErrFolderExist       = 9 // folder with provided name already exist exist
)

type HttpError struct {
	cause    error
	details  string
	errCode  int
	httpCode int
}

func NewHttpError(cause error, details string, errCode int, httpCode int) *HttpError {
	return &HttpError{cause: cause, details: details, errCode: errCode, httpCode: httpCode}
}

func (e *HttpError) Error() string {
	return e.cause.Error()
}

func (e *HttpError) Details() string {
	return e.details
}

func (e *HttpError) ErrCode() int {
	return e.errCode
}

func (e *HttpError) HttpCode() int {
	if e.httpCode == 0 {
		return http.StatusInternalServerError
	}
	return e.httpCode
}

// it's also a wrapper.
func (e *HttpError) Cause() error  { return e.cause }
func (e *HttpError) Unwrap() error { return e.cause }

type stackTracer interface {
	StackTrace() errors.StackTrace
}

// SendErrorJSON makes {"error": "blah", "details": "blah", "code": 0} json body and responds with error code
func SendErrorJSON(c *gin.Context, httpStatusCode int, err error, details string, errCode int) {
	fields := log.Fields{
		"http_status": httpStatusCode,
		"details":     details,
		"error_code":  errCode,
	}
	if err, ok := errors.Cause(err).(stackTracer); ok {
		fields["trace"] = fmt.Sprintf("%+v", err.StackTrace()[:5])
	}
	log.WithContext(c.Request.Context()).WithFields(fields).Warnf("%v", err)

	c.AbortWithStatusJSON(httpStatusCode, gin.H{"error": err.Error(), "details": details, "code": errCode})
}

func SendNotFoundError(c *gin.Context, err error, details string) {
	SendErrorJSON(c, http.StatusNotFound, err, details, ErrNotFound)
}

func SendInternalError(c *gin.Context, err error, details string) {
	SendErrorJSON(c, http.StatusInternalServerError, err, details, ErrCodeInternal)
}

func SendBadRequestError(c *gin.Context, err error, details string) {
	SendErrorJSON(c, http.StatusBadRequest, err, details, ErrBadRequest)
}
