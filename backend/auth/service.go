package auth

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"strings"

	"github.com/filebrowser/filebrowser/v3/config"
	authToken "github.com/go-pkgz/auth/token"
	"github.com/google/uuid"

	"github.com/filebrowser/filebrowser/v3/hash"
	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/filebrowser/filebrowser/v3/store"
)

type Service struct {
	cfg       *config.Config
	hasher    hash.Hasher
	userStore store.UserStore
}

func NewService(cfg *config.Config, userStore store.UserStore, hasher hash.Hasher) *Service {
	return &Service{cfg: cfg, userStore: userStore, hasher: hasher}
}

// Check implements provider.CredChecker interface
func (s *Service) Check(username, password string) (ok bool, err error) {
	ok, _, err = s.validateUser(context.Background(), username, password)
	fmt.Println(ok, err)
	return ok, err
}

// BasicAuthChecker implements middleware.BasicAuthFunc
func (s *Service) BasicAuthChecker(username, password string) (ok bool, userInfo authToken.User, err error) {
	ok, user, err := s.validateUser(context.Background(), username, password)
	if err != nil {
		return false, authToken.User{}, err
	}
	if !ok {
		return false, authToken.User{}, nil
	}
	return true, authToken.User{
		Name: user.Name,
		ID:   user.ID,
	}, nil
}

func (s *Service) validateUser(ctx context.Context, username, password string) (bool, *store.User, error) {
	user, err := s.userStore.GetByUsername(ctx, username)
	switch {
	case errors.Is(err, store.ErrNotFound):
		return false, nil, nil
	case err != nil:
		return false, nil, err
	}

	if !s.hasher.CheckPassword(password, user.Password) {
		return false, nil, nil
	}

	return true, user, nil
}

// Update implements token.ClaimsUpdater interface
func (s *Service) Update(claims authToken.Claims) authToken.Claims {
	if claims.User == nil {
		return claims
	}

	ctx := context.Background()

	user, err := s.userStore.Get(ctx, claims.User.ID)
	switch {
	// new login with external provider
	case errors.Is(err, store.ErrNotFound):
		username := string(strings.ReplaceAll(uuid.NewString(), "-", "")[10])
		user = &store.User{
			ID:           claims.User.ID,
			Provider:     strings.Split(claims.User.ID, "_")[0],
			Username:     username,
			Password:     "",
			Home:         s.cfg.Auth.User.Home,
			Name:         claims.User.Name,
			Locale:       s.cfg.Locale,
			LockPassword: false,
			Blocked:      false,
		}
		if err := s.userStore.Create(ctx, user); err != nil {
			log.WithContext(ctx).Errorf("failed to create user: %+v", err)
			return claims
		}
	}
	claims.User.SetBoolAttr("blocked", user.Blocked)

	return claims
}

// Validate implements token.Validator
func (s *Service) Validate(token string, claims authToken.Claims) bool {
	if claims.User == nil {
		return false
	}
	return !claims.User.BoolAttr("blocked")
}

// InitAdminUser initializes the admin user with root privileges. Provided password must be encrypted with hash.Password()
func (s *Service) InitAdminUser(ctx context.Context, pwd string) error {
	user := &store.User{
		ID:           "local_" + authToken.HashID(sha1.New(), "admin"), //nolint:gosec,
		Name:         "Admin",
		Provider:     "local",
		Username:     "admin",
		Password:     pwd,
		Home:         "/",
		Locale:       s.cfg.Locale,
		LockPassword: false,
		Blocked:      false,
	}

	_, err := s.userStore.Get(ctx, user.ID)
	switch {
	case errors.Is(err, store.ErrNotFound):
		return s.userStore.Create(ctx, user)
	case err != nil:
		return err
	default:
		return s.userStore.Update(ctx, user)
	}
}

// InitGuestUser initializes the guest user
func (s *Service) InitGuestUser(ctx context.Context) error {
	user := &store.User{
		ID:           "local_" + authToken.HashID(sha1.New(), "guest"), //nolint:gosec,
		Name:         "Guest",
		Provider:     "local",
		Username:     "guest",
		Password:     "",
		Home:         "/",
		Locale:       s.cfg.Locale,
		LockPassword: true,
		Blocked:      false,
	}

	_, err := s.userStore.Get(ctx, user.ID)
	switch {
	case errors.Is(err, store.ErrNotFound):
		return s.userStore.Create(ctx, user)
	case err != nil:
		return err
	default:
		return s.userStore.Update(ctx, user)
	}
}
