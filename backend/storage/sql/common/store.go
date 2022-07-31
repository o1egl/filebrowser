package common

import (
	"database/sql"

	"github.com/doug-martin/goqu/v9"
)

type Store struct {
	db *goqu.Database
}

func NewStore(dialect string, db *sql.DB) *Store {
	return &Store{
		db: goqu.New(dialect, db),
	}
}
