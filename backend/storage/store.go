package storage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/filebrowser/filebrowser/config"
	"github.com/filebrowser/filebrowser/domain"
	"os"
	"path/filepath"
)

type Store interface {
	UserStore
	VolumeStore
}

type UserStore interface {
	GetUser(ctx context.Context, id int64) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
}

type VolumeStore interface {
	GetUserVolumes(ctx context.Context, userID int64) ([]*domain.VolumeWithPermissions, error)
	CreateVolume(ctx context.Context, volume *domain.Volume) error
}

func New(cfg config.Store) (Store, error) {
	switch cfg.Type {
	case config.StoreTypeSqlite:
		return newSqliteStore(cfg)
	case config.StoreTypePostgres:
		return newPostgresStore(cfg)
	case config.StoreTypeMysql:
		return newMysqlStore(cfg)
	default:
		return nil, fmt.Errorf("unsupported store type %s", cfg.Type)
	}
}

func newSqliteStore(cfg config.Store) (Store, error) {
	if err := makeDirs(filepath.Dir(cfg.SQLite.File)); err != nil {
		return nil, fmt.Errorf("failed to create sqlite store: %w", err)
	}
	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc&_fk=1", cfg.SQLite.File))
	if err != nil {
		return nil, err
	}
	if err := migrations.Up(db, config.StoreTypeSqlite); err != nil {
		return nil, err
	}
	return sqlite3.NewStore(db), nil
}

func newPostgresStore(cfg config.Store) (Store, error) {
	db, err := sql.Open("postgres", cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}
	if err := migrations.Up(db, config.StoreTypePostgres); err != nil {
		return nil, err
	}
	return postgres.NewStore(db), nil
}

func newMysqlStore(cfg config.Store) (Store, error) {
	db, err := sql.Open("mysql", cfg.Postgres.DSN)
	if err != nil {
		return nil, err
	}
	if err := migrations.Up(db, config.StoreTypeMysql); err != nil {
		return nil, err
	}
	return mysql.NewStore(db), nil
}

// mkdir -p for all dirs
func makeDirs(dirs ...string) error {
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0700); err != nil { // If path is already a directory, MkdirAll does nothing
			return fmt.Errorf("can't make directory %s: %w", dir, err)
		}
	}
	return nil
}
