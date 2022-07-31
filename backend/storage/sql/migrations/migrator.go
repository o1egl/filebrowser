package migrations

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migrateDatabase "github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source"
	"github.com/golang-migrate/migrate/v4/source/iofs"

	"github.com/filebrowser/filebrowser/config"
)

//go:embed */*.sql
var fs embed.FS

func Up(db *sql.DB, storeType config.StoreType) error {
	src, err := makeSource(storeType)
	if err != nil {
		return err
	}

	driver, err := makeDriver(db, storeType)
	if err != nil {
		return fmt.Errorf("failed to initialize migration driver: %w", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", src, "database", driver)
	if err != nil {
		return err
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func makeSource(storeType config.StoreType) (source.Driver, error) {
	var path string
	switch storeType {
	case config.StoreTypeSqlite:
		path = "sqlite"
	case config.StoreTypePostgres:
		path = "postgres"
	case config.StoreTypeMysql:
		path = "mysql"
	default:
		return nil, fmt.Errorf("unsupported store type %s", storeType)
	}
	return iofs.New(fs, path)
}

func makeDriver(db *sql.DB, storeType config.StoreType) (migrateDatabase.Driver, error) {
	switch storeType {
	case config.StoreTypeSqlite:
		return sqlite3.WithInstance(db, &sqlite3.Config{})
	case config.StoreTypePostgres:
		return postgres.WithInstance(db, &postgres.Config{})
	case config.StoreTypeMysql:
		return mysql.WithInstance(db, &mysql.Config{})
	default:
		return nil, fmt.Errorf("unsupported store type %s", storeType)
	}
}
