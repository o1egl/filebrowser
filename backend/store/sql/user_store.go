package sql

import (
	"context"

	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/sql/conv"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent"
)

type UserStore struct {
	client    *ent.Client
	converter UserConverter
}

func NewUserStore(client *ent.Client) *UserStore {
	return &UserStore{
		client:    client,
		converter: &conv.UserConverterImpl{},
	}
}

func (u *UserStore) Get(ctx context.Context, id string) (*store.User, error) {
	user, err := u.client.User.Get(ctx, id)
	switch {
	case ent.IsNotFound(err):
		return nil, store.ErrNotFound
	case err != nil:
		return nil, err
	}
	return u.converter.Convert(user), nil
}

func (u *UserStore) Save(ctx context.Context, user *store.User) (id string, err error) {
	newUser, err := u.client.User.Create().
		SetProvider(user.Provider).
		SetUsername(user.Username).
		SetPassword(user.Password).
		SetName(user.Name).
		SetScope(user.Scope).
		SetLocale(user.Locale).
		SetLockPassword(user.LockPassword).
		SetBlocked(user.Blocked).
		Save(ctx)
	if err != nil {
		return "", err
	}

	return newUser.ID, nil
}
