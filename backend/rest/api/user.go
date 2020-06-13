package api

import (
	"context"
	"net/http"

	"github.com/go-pkgz/auth/token"

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
func MustGetUser(r *http.Request) *store.User {
	user := UserFromContext(r.Context())
	if user == nil {
		panic("user not found in request context")
	}
	return user
}

type userGetterStore interface {
	FindUserByID(ctx context.Context, id string) (*store.User, error)
}

func UserMiddleware(userStore userGetterStore) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			userInfo, err := token.GetUserInfo(r)
			if err != nil {
				//rest.SendErrorJSONChi(w, r, http.StatusUnauthorized, err, "unauthorized user", rest.ErrCodeUnauthorized)
				return
			}

			user, err := userStore.FindUserByID(r.Context(), userInfo.ID)
			if err != nil {
				//rest.SendErrorJSONChi(w, r, http.StatusUnauthorized, err, "unauthorized user", rest.ErrCodeUnauthorized)
				return
			}
			r = r.WithContext(UserToContext(r.Context(), user))
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
