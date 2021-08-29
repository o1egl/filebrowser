package filebrowser

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/filebrowser/filebrowser/v3/auth"
	"github.com/filebrowser/filebrowser/v3/service"
	"github.com/filebrowser/filebrowser/v3/store"
)

type Service struct {
	userStore   store.UserStore
	volumeStore store.VolumeStore
}

func (s *Service) List(ctx context.Context, user auth.User, params service.ListParams) (*service.FileWithChildren, error) {
	dbUser, err := s.userStore.Get(ctx, user.ID)
	switch {
	case errors.Is(err, store.ErrNotFound):
		return nil, service.NewNotFoundError(err, service.ResourceUser, "id", user.ID)
	case err != nil:
		return nil, err
	}

	volume := store.Volume{
		ID:    0,
		Label: "Home",
		Path:  dbUser.Home,
	}
	if params.Volume != service.HomeVolumeID {
		volume, err = s.loadUserVolume(ctx, dbUser, params.Volume)
	}

	return nil, nil
}

func (s *Service) Create(ctx context.Context, user auth.User, volume int64, filename string, content io.Reader) error {
	panic("implement me")
}

func (s *Service) Update(ctx context.Context, user auth.User, volume int64, filename string, content io.Reader) error {
	panic("implement me")
}

func (s *Service) Delete(ctx context.Context, user auth.User, volume int64, filename string) error {
	panic("implement me")
}

func (s *Service) loadUserVolume(ctx context.Context, userID string, volumeID int64) (*store.Volume, error) {
	volumes, err := s.volumeStore.GetUserVolumes(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, volume := range volumes {
		if volume.ID == volumeID {
			return volume, nil
		}
	}
	return nil, service.NewAccessDeniedError(fmt.Sprintf("user %s", userID), fmt.Sprintf("volume %d", volumeID))
}
