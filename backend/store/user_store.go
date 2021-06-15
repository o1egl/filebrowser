//go:generate mockgen -source=$GOFILE -destination mock/user_store.go
package store

import "context"

type UserStore interface {
	Get(ctx context.Context, id string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
}

// User holds user-related info
type User struct {
	ID           string `json:"id"`
	Provider     string `json:"provider"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Home         string `json:"home"`
	Name         string `json:"name"`
	Locale       string `json:"locale"`
	LockPassword bool   `json:"lock_password"`
	Blocked      bool   `json:"blocked"`
}
