package middleware

import (
	"net/http"

	"github.com/filebrowser/filebrowser/v3/auth"
	"github.com/gin-gonic/gin"
	authToken "github.com/go-pkgz/auth/token"

	"github.com/filebrowser/filebrowser/v3/rest"
)

func User() gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := authToken.GetUserInfo(c.Request)
		if err != nil {
			rest.SendErrorJSON(c, http.StatusUnauthorized, err, "unauthorized user", rest.ErrCodeUnauthorized)
			return
		}
		user := auth.FromToken(userInfo)
		c.Request = c.Request.WithContext(auth.UserToContext(c.Request.Context(), user))

		c.Next()
	}
}
