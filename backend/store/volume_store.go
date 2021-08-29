package store

import "context"

type VolumeStore interface {
	GetUserVolumes(ctx context.Context, userID string) ([]*Volume, error)
}

type Volume struct {
	ID    int64
	Label string
	Path  string
}
