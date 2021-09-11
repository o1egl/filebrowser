package sql

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/filebrowser/filebrowser/v3/config"
	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

var Set = wire.NewSet(
	EntClientProvider,
	NewUserStore,
	wire.Bind(new(store.UserStore), new(*UserStore)),
	NewVolumeStore,
	wire.Bind(new(store.VolumeStore), new(*VolumeStore)),
)

func EntClientProvider(ctx context.Context, cfg *config.Config) (client *ent.Client, err error) {
	switch cfg.Store.Type {
	case config.StoreTypeSqlite:
		if err = makeDirs(filepath.Dir(cfg.Store.SQLite.File)); err != nil {
			return nil, errors.Wrap(err, "failed to create sqlite store")
		}
		client, err = ent.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc&_fk=1", cfg.Store.SQLite.File))
	case config.StoreTypePostgres:
		client, err = ent.Open("postgres", cfg.Store.Postgres.DSN)
	case config.StoreTypeMysql:
		client, err = ent.Open("mysql", cfg.Store.Postgres.DSN)
	default:
		return nil, errors.Errorf("unsupported store type %s", cfg.Store.Type)
	}
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize data store")
	}

	log.WithContext(ctx).Debugf("Apply schema migrations")
	if err := client.Schema.Create(ctx); err != nil {
		return nil, err
	}

	return client, nil
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
