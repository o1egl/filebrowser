package sql

import (
	"context"

	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/sql/conv"
	"github.com/filebrowser/filebrowser/v3/store/sql/conv/generated"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent/user"
)

type UserStore struct {
	client    *ent.Client
	converter conv.UserConverter
}

func NewUserStore(client *ent.Client) *UserStore {
	return &UserStore{
		client:    client,
		converter: &generated.UserConverterImpl{},
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

func (u *UserStore) GetByUsername(ctx context.Context, username string) (*store.User, error) {
	usr, err := u.client.User.Query().Where(user.UsernameEqualFold(username)).Only(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, store.ErrNotFound
	case err != nil:
		return nil, err
	}
	return u.converter.Convert(usr), nil
}

func (u *UserStore) Create(ctx context.Context, user *store.User) error {
	userBuilder := u.client.User.Create().
		SetID(user.ID).
		SetProvider(user.Provider).
		SetUsername(user.Username).
		SetName(user.Name).
		SetHome(user.Home).
		SetLocale(user.Locale).
		SetLockPassword(user.LockPassword).
		SetBlocked(user.Blocked)

	if user.Password != "" {
		userBuilder = userBuilder.SetPassword(user.Password)
	}

	_, err := userBuilder.Save(ctx)
	return err
}

func (u *UserStore) Update(ctx context.Context, user *store.User) error {
	userBuilder := u.client.User.UpdateOneID(user.ID).
		SetUsername(user.Username).
		SetName(user.Name).
		SetHome(user.Home).
		SetLocale(user.Locale).
		SetLockPassword(user.LockPassword).
		SetBlocked(user.Blocked)

	if user.Password != "" {
		userBuilder = userBuilder.SetPassword(user.Password)
	}

	_, err := userBuilder.Save(ctx)
	return err
}
