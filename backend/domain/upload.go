package domain

import "time"

type Upload struct {
	ID        int64
	UserID    int64
	VolumeID  int64
	Path      string
	TmpPath   string
	Size      int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
