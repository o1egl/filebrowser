package mysql

import (
	"database/sql"

	"github.com/filebrowser/filebrowser/storage/sql/common"
)

type Store struct {
	*common.Store
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		Store: common.NewStore("mysql", db),
	}
}
