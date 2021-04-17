package sql

import (
	"context"

	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/sql/conv"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent"
	user "github.com/filebrowser/filebrowser/v3/store/sql/ent/user"
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
	usr, err := u.client.User.Get(ctx, id)
	switch {
	case ent.IsNotFound(err):
		return nil, store.ErrNotFound
	case err != nil:
		return nil, err
	}
	return u.converter.Convert(usr), nil
}

func (u *UserStore) GetByUsernameAndProvider(ctx context.Context, username, provider string) (*store.User, error) {
	usr, err := u.client.User.Query().Where(user.And(
		user.UsernameEqualFold(username),
		user.ProviderEQ(provider),
	)).Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, store.ErrNotFound
	case err != nil:
		return nil, err
	}
	return u.converter.Convert(usr), nil
}

func (u *UserStore) Save(ctx context.Context, user *store.User) error {
	_, err := u.client.User.Create().
		SetID(user.ID).
		SetProvider(user.Provider).
		SetUsername(user.Username).
		SetPassword(user.Password).
		SetName(user.Name).
		SetScope(user.Scope).
		SetLocale(user.Locale).
		SetLockPassword(user.LockPassword).
		SetBlocked(user.Blocked).
		Save(ctx)

	return err
}
