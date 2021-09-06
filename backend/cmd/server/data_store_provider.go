package server

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/sql"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent"
	"github.com/google/wire"
	"github.com/pkg/errors"
)

var DataStoreSet = wire.NewSet(
	EntClientProvider,
	sql.NewUserStore,
	wire.Bind(new(store.UserStore), new(*sql.UserStore)),
	sql.NewVolumeStore,
	wire.Bind(new(store.VolumeStore), new(*sql.VolumeStore)),
)

func EntClientProvider(ctx context.Context, srvCmd *ServerCommand) (client *ent.Client, err error) {
	switch srvCmd.Store.Type {
	case "sqlite":
		if err = makeDirs(filepath.Dir(srvCmd.Store.SQLite.File)); err != nil {
			return nil, errors.Wrap(err, "failed to create sqlite store")
		}
		client, err = ent.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc&_fk=1", srvCmd.Store.SQLite.File))
	case "postgres":
		client, err = ent.Open("postgres", srvCmd.Store.Postgres.DSN)
	case "mysql":
		client, err = ent.Open("mysql", srvCmd.Store.Postgres.DSN)
	default:
		return nil, errors.Errorf("unsupported store type %s", srvCmd.Store.Type)
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
