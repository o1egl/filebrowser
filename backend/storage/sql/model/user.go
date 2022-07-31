package model

type User struct {
	ID           int64   `db:"id" goqu:"skipinsert,skipupdate"`
	Username     string  `db:"username"`
	Password     *string `db:"password"`
	Name         string  `db:"name"`
	Locale       string  `db:"locale"`
	LockPassword bool    `db:"lock_password"`
	Blocked      bool    `db:"blocked"`
	Home         string  `db:"home"`
	Read         bool    `db:"read"`
	Create       bool    `db:"create"`
	Modify       bool    `db:"modify"`
	Delete       bool    `db:"delete"`
	Share        bool    `db:"share"`
	IsAdmin      bool    `db:"is_admin"`
}
