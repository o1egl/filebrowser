package model

type Volume struct {
	ID          int64  `db:"id" goqu:"skipinsert,skipupdate"`
	Label       string `db:"label"`
	Path        string `db:"path"`
	Description string `db:"description"`
}
