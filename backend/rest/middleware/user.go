package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	authToken "github.com/go-pkgz/auth/token"

	"github.com/filebrowser/filebrowser/v3/backend/rest"
	"github.com/filebrowser/filebrowser/v3/backend/store"
)

type UserStore interface {
	FindUserByID(ctx context.Context, id string) (*store.User, error)
}

func User(userStore UserStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInfo, err := authToken.GetUserInfo(c.Request)
		if err != nil {
			rest.SendErrorJSON(c, http.StatusUnauthorized, err, "unauthorized user", rest.ErrCodeUnauthorized)
			return
		}

		user, err := userStore.FindUserByID(c.Request.Context(), userInfo.ID)
		if err != nil {
			rest.SendErrorJSON(c, http.StatusUnauthorized, err, "unauthorized user", rest.ErrCodeUnauthorized)
			return
		}
		c.Request = c.Request.WithContext(rest.UserToContext(c.Request.Context(), user))

		c.Next()
	}
}
