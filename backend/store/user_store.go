package store

import "context"

type UserStore interface {
	Get(ctx context.Context, id string) (*User, error)
	GetByUsernameAndProvider(ctx context.Context, username, provider string) (*User, error)
	Save(ctx context.Context, user *User) error
}

// User holds user-related info
type User struct {
	ID           string `json:"id"`
	Provider     string `json:"provider"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Name         string `json:"name"`
	Scope        string `json:"scope"`
	Locale       string `json:"locale"`
	LockPassword bool   `json:"lockPassword"`
	Blocked      bool   `json:"blocked"`
}
