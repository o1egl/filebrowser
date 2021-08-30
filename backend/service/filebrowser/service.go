package filebrowser

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/filebrowser/filebrowser/v3/auth"
	"github.com/filebrowser/filebrowser/v3/filesystem"
	"github.com/filebrowser/filebrowser/v3/mathx"
	"github.com/filebrowser/filebrowser/v3/service"
	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/spf13/afero"
)

type Service struct {
	rootFs      afero.Fs
	userStore   store.UserStore
	volumeStore store.VolumeStore
}

func New(rootFs afero.Fs, userStore store.UserStore, volumeStore store.VolumeStore) *Service {
	return &Service{rootFs: rootFs, userStore: userStore, volumeStore: volumeStore}
}

func (s *Service) List(ctx context.Context, user auth.User, params service.ListParams) (*service.FileWithChildren, error) {
	volume, err := s.loadUserVolume(ctx, user.ID, params.Volume)
	if err != nil {
		return nil, err
	}

	info, err := filesystem.Stat(volume.Fs, params.Filename)
	if err != nil {
		return nil, fsError(err, volume.Path, params.Filename)
	}

	var children []service.File
	if info.IsDir {
		dir, err := filesystem.ReadDir(volume.Fs, params.Filename)
		if err != nil {
			return nil, fsError(err, volume.Path, params.Filename)
		}
		children = make([]service.File, len(dir))
		for i, file := range dir {
			children[i] = service.File(file)
		}
		children = sortFiles(children, params.GroupBy, params.SortBy, params.OrderBy)
	}

	// fill metadata
	fileMetaData := service.FileMetaData{
		FilesCount: 0,
		DirsCount:  0,
		TotalCount: len(children),
	}
	for _, child := range children {
		if child.IsDir {
			fileMetaData.DirsCount++
			continue
		}
		fileMetaData.FilesCount++
	}

	// apply offset/limit
	offset := mathx.MinInt(params.Offset, len(children))
	limit := len(children)
	if params.Limit > 0 {
		limit = mathx.MinInt(offset+params.Limit, len(children))
	}
	children = children[offset:limit]

	return &service.FileWithChildren{
		File:     service.File(info),
		Children: children,
		Meta:     fileMetaData,
	}, nil
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

func (s *Service) loadUser(ctx context.Context, userID string) (*store.User, error) {
	user, err := s.userStore.Get(ctx, userID)
	switch {
	case errors.Is(err, store.ErrNotFound):
		return nil, service.NewNotFoundError(err, service.ResourceUser, "id", userID)
	case err != nil:
		return nil, err
	}
	return user, nil
}

type VolumeInfo struct {
	ID    int64
	Label string
	Path  string
	Fs    afero.Fs
}

func (s *Service) loadUserVolume(ctx context.Context, userID string, volumeID int64) (*VolumeInfo, error) {
	if volumeID == service.HomeVolumeID {
		user, err := s.loadUser(ctx, userID)
		if err != nil {
			return nil, err
		}
		return &VolumeInfo{
			ID:    0,
			Label: service.HomeVolumeLabel,
			Path:  user.Home,
			Fs:    afero.NewBasePathFs(s.rootFs, user.Home),
		}, nil
	}

	volumes, err := s.volumeStore.GetUserVolumes(ctx, userID)
	if err != nil {
		return nil, err
	}
	for _, volume := range volumes {
		if volume.ID == volumeID {
			return &VolumeInfo{
				ID:    volume.ID,
				Label: volume.Label,
				Path:  volume.Path,
				Fs:    afero.NewBasePathFs(s.rootFs, volume.Path),
			}, nil
		}
	}
	return nil, service.NewAccessDeniedError(fmt.Sprintf("user %s", userID), fmt.Sprintf("volume %d", volumeID))
}

func fsError(err error, volumePath, filePath string) error {
	if errors.Is(err, store.ErrNotFound) {
		return service.NewNotFoundError(err, service.ResourceFile, "volume", volumePath, "path", filePath)
	}
	return err
}
