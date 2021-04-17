package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	authToken "github.com/go-pkgz/auth/token"

	"github.com/filebrowser/filebrowser/v3/rest"
	"github.com/filebrowser/filebrowser/v3/store"
)

func User(userStore store.UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := authToken.GetUserInfo(c.Request)
		if err != nil {
			rest.SendErrorJSON(c, http.StatusUnauthorized, err, "unauthorized user", rest.ErrCodeUnauthorized)
			return
		}
		user, err := userStore.Get(c.Request.Context(), userInfo.ID)
		if err != nil {
			rest.SendErrorJSON(c, http.StatusUnauthorized, err, "unauthorized user", rest.ErrCodeUnauthorized)
			return
		}
		c.Request = c.Request.WithContext(rest.UserToContext(c.Request.Context(), user))

		c.Next()
	}
}
