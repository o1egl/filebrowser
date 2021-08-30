package auth

import (
	"context"

	"github.com/go-pkgz/auth/token"
)

type userKeyType int

const userKey userKeyType = iota

type User struct {
	ID       string `json:"id"`
	Provider string `json:"provider"`
	Username string `json:"username"`
	Name     string `json:"name"`
	Locale   string `json:"locale"`
	Blocked  bool   `json:"blocked"`
}

func UserToContext(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(userKey).(User)
	return user, ok
}

const (
	attrProvider = "provider"
	attrUsername = "username"
	attrLocale   = "locale"
	attrBlocked  = "blocked"
)

func FromToken(tokenUser token.User) User {
	user := User{
		ID:       tokenUser.ID,
		Provider: tokenUser.StrAttr(attrProvider),
		Username: tokenUser.StrAttr(attrUsername),
		Name:     tokenUser.Name,
		Locale:   tokenUser.StrAttr(attrLocale),
		Blocked:  tokenUser.BoolAttr(attrBlocked),
	}
	return user
}
