package rest

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/filebrowser/filebrowser/v3/backend/store"
)

type userKeyType int

const userKey userKeyType = iota

func UserToContext(ctx context.Context, user *store.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) *store.User {
	user, ok := ctx.Value(userKey).(*store.User)
	if ok {
		return user
	}
	return nil
}

// MustGetUserInfo fails if can't extract user data from the request.
// should be called from authed controllers only
func MustGetUser(c *gin.Context) *store.User {
	user := UserFromContext(c.Request.Context())
	if user == nil {
		panic("user not found in request context")
	}
	return user
}
