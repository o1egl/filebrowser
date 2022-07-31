package model

import "time"

type Upload struct {
	ID        int64     `db:"id" goqu:"skipinsert,skipupdate"`
	UserID    int64     `db:"user_id"`
	VolumeID  int64     `db:"volume_id"`
	Path      string    `db:"path"`
	TmpPath   string    `db:"tmp_path"`
	Size      int64     `db:"size"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
