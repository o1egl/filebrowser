package engine

import (
	"context"

	"github.com/filebrowser/filebrowser/v3/store"
)

// Interface defines methods provided by low-level storage engine
type Interface interface {
	SaveUser(ctx context.Context, user *store.User) error
	FindUsers(ctx context.Context, request FindUserRequest) ([]*store.User, error)
	DeleteUser(ctx context.Context, userID string) error
}

type FindUserRequest struct {
	ID       string
	Username string
}
