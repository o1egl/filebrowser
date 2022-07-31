package common

import (
	"context"
	"strconv"

	"github.com/doug-martin/goqu/v9"

	"github.com/filebrowser/filebrowser/domain"
	"github.com/filebrowser/filebrowser/storage/sql/conv"
	"github.com/filebrowser/filebrowser/storage/sql/model"
)

const (
	usersTable            = "users"
	userAuthProviderTable = "user_auth_providers"
)

func (s *Store) GetUser(ctx context.Context, id int64) (*domain.User, error) {
	var userModel model.User
	found, err := s.db.Select().
		From(usersTable).
		Where(goqu.C("id").Eq(id)).
		Executor().
		ScanStructContext(ctx, &userModel)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, domain.NewNotFoundError(err, domain.ResourceUser, "id", strconv.FormatInt(id, 10))
	}

	return conv.UserToDomain(&userModel), nil
}

func (s *Store) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var userModel model.User
	found, err := s.db.Select().
		From(usersTable).
		Where(goqu.C("username").Eq(username)).
		Executor().
		ScanStructContext(ctx, &userModel)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, domain.NewNotFoundError(err, domain.ResourceUser, "username", username)
	}

	return conv.UserToDomain(&userModel), nil
}

func (s *Store) GetUserByAuthProvider(ctx context.Context, provider string, token string) (*domain.User, error) {
	var userModel model.User
	found, err := s.db.Select().
		From(usersTable).As("s").
		Join(goqu.T(userAuthProviderTable).As("p"), goqu.On(goqu.Ex{"s.id": goqu.I("p.user_id")})).
		Where(goqu.I("p.id").Eq(token)).
		Executor().
		ScanStructContext(ctx, &userModel)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, domain.NewNotFoundError(err, domain.ResourceUser)
	}

	return conv.UserToDomain(&userModel), nil
}

func (s *Store) CreateUser(ctx context.Context, user *domain.User) error {
	res, err := s.db.Insert(usersTable).
		Rows(conv.UserFromDomain(user)).
		Executor().
		ExecContext(ctx)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (s *Store) UpdateUser(ctx context.Context, user *domain.User) error {
	userModel := conv.UserFromDomain(user)
	_, err := s.db.Update(usersTable).
		Set(userModel).
		Where(goqu.C("id").Eq(user.ID)).Executor().ExecContext(ctx)
	return err
}
