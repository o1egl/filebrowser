package engine

import (
	"context"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"

	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/filebrowser/filebrowser/v3/store"
)

// BoltDB implements store.Interface
type BoltDB struct {
	db *storm.DB
}

// ensure that BoltDB implements storage interface
var _ Interface = &BoltDB{}

// NewBoltDB makes persistent boltdb-based store. For each site new boltdb file created
func NewBoltDB(ctx context.Context, fileName string, options *bolt.Options) (*BoltDB, error) {
	log.WithContext(ctx).Infof("bolt store: options %+v", *options)
	db, err := storm.Open(fileName, storm.BoltOptions(0600, options)) //nolint:gomnd
	if err != nil {
		return nil, errors.Wrap(err, "failed to open bolt db")
	}

	res := &BoltDB{db: db}

	err = res.withTx(true, func(tx storm.Node) error {
		if err := tx.Init(&store.User{}); err != nil {
			return errors.Wrap(err, "failed to create user bucket")
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize buckets")
	}

	return res, nil
}

func (b *BoltDB) SaveUser(ctx context.Context, user *store.User) error {
	err := b.withTx(true, func(tx storm.Node) error {
		return tx.Save(user)
	})
	if err != nil {
		return err
	}
	return nil
}

func (b *BoltDB) FindUsers(ctx context.Context, request FindUserRequest) ([]*store.User, error) {
	var matchers []q.Matcher
	if request.ID != "" {
		matchers = append(matchers, q.Eq("ID", request.ID))
	}
	if request.Username != "" {
		matchers = append(matchers, q.Eq("Username", request.Username))
	}
	var users []*store.User
	if err := b.db.Select(matchers...).Find(&users); err != nil && err != storm.ErrNotFound {
		return nil, errors.Wrap(err, "failed execute user find query")
	}
	return users, nil
}

func (b *BoltDB) DeleteUser(ctx context.Context, userID string) error {
	if err := b.db.Delete("User", userID); err != nil {
		return errors.Wrapf(err, "can't delete user with id %s", userID)
	}
	return nil
}

func (b *BoltDB) withTx(writable bool, fn func(tx storm.Node) error) error {
	tx, err := b.db.Begin(writable)
	if err != nil {
		return err
	}
	if fnErr := fn(tx); fnErr != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return errors.Wrapf(err, "failed (rollback error %v)", errRollback)
		}
		return err
	}
	return tx.Commit()
}
