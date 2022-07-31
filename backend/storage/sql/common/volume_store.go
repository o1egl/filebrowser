package common

import (
	"context"

	"github.com/filebrowser/filebrowser/domain"
)

func (s *Store) GetUserVolumes(ctx context.Context, userID int64) ([]*domain.VolumeWithPermissions, error) {
	return nil, nil
}

func (s *Store) CreateVolume(ctx context.Context, volume *domain.Volume) error {
	return nil
}

func (s *Store) UpdateVolume(ctx context.Context, volume *domain.Volume) error {
	return nil
}
