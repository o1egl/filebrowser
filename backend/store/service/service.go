package service

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/engine"
)

// DataStore wraps store.Interface with additional methods
type DataStore struct {
	Engine engine.Interface

	//nolint:unused
	// granular locks
	scopedLocks struct {
		sync.Mutex
		sync.Once
		locks map[string]sync.Locker
	}
}

func (d *DataStore) SaveUser(ctx context.Context, user *store.User) error {
	return d.Engine.SaveUser(ctx, user)
}

func (d *DataStore) FindUserByID(ctx context.Context, userID string) (*store.User, error) {
	users, err := d.Engine.FindUsers(ctx, engine.FindUserRequest{
		ID: userID,
	})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.Wrapf(store.ErrNotFound, "user id: %s", userID)
	}
	return users[0], nil
}

func (d *DataStore) FindUserByUsername(ctx context.Context, username string) (*store.User, error) {
	users, err := d.Engine.FindUsers(ctx, engine.FindUserRequest{
		Username: username,
	})
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, errors.Wrapf(store.ErrNotFound, "username: %s", username)
	}
	return users[0], nil
}

func (d *DataStore) DeleteUser(ctx context.Context, userID string) error {
	return d.Engine.DeleteUser(ctx, userID)
}

func (d *DataStore) IsAdmin(ctx context.Context, userID string) bool {
	user, err := d.FindUserByID(ctx, userID)
	if err != nil {
		return false
	}
	return user.Permissions.Admin
}

// getScopedLocks pull lock from the map if found or create a new one
//nolint:unused
func (d *DataStore) getScopedLocks(id string) (lock sync.Locker) {
	d.scopedLocks.Do(func() { d.scopedLocks.locks = map[string]sync.Locker{} })

	d.scopedLocks.Lock()
	lock, ok := d.scopedLocks.locks[id]
	if !ok {
		lock = &sync.Mutex{}
		d.scopedLocks.locks[id] = lock
	}
	d.scopedLocks.Unlock()

	return lock
}
