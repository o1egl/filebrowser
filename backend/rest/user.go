package rest

import (
	"github.com/filebrowser/filebrowser/v3/auth"
	"github.com/gin-gonic/gin"
)

// MustGetUser fails if can't extract user data from the request.
// should be called from authed controllers only
func MustGetUser(c *gin.Context) auth.User {
	user, ok := auth.UserFromContext(c.Request.Context())
	if !ok {
		panic("user not found in request context")
	}
	return user
}
