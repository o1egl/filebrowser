package filebrowser

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/filebrowser/filebrowser/v3/auth"
	"github.com/filebrowser/filebrowser/v3/domain"
	"github.com/filebrowser/filebrowser/v3/filesystem"
	"github.com/filebrowser/filebrowser/v3/mathx"
	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/spf13/afero"
)

type ServiceImpl struct {
	rootFs      afero.Fs
	userStore   store.UserStore
	volumeStore store.VolumeStore
}

func New(rootFs afero.Fs, userStore store.UserStore, volumeStore store.VolumeStore) *ServiceImpl {
	return &ServiceImpl{rootFs: rootFs, userStore: userStore, volumeStore: volumeStore}
}

func (s *ServiceImpl) List(ctx context.Context, user auth.User, params ListParams) (*FileWithChildren, error) {
	volume, err := s.loadUserVolume(ctx, user.ID, params.Volume)
	if err != nil {
		return nil, err
	}

	info, err := filesystem.Stat(volume.Fs, params.Filename)
	if err != nil {
		return nil, fsError(err, volume.Path, params.Filename)
	}

	var children []File
	if info.IsDir {
		dir, err := filesystem.ReadDir(volume.Fs, params.Filename)
		if err != nil {
			return nil, fsError(err, volume.Path, params.Filename)
		}
		children = make([]File, len(dir))
		for i, file := range dir {
			children[i] = File(file)
		}
		children = sortFiles(children, params.GroupBy, params.SortBy, params.OrderBy)
	}

	// fill metadata
	fileMetaData := FileMetaData{
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

	return &FileWithChildren{
		File:     File(info),
		Children: children,
		Meta:     fileMetaData,
	}, nil
}

func (s *ServiceImpl) Create(ctx context.Context, user auth.User, volume int64, filename string, content io.Reader) error {
	panic("implement me")
}

func (s *ServiceImpl) Update(ctx context.Context, user auth.User, volume int64, filename string, content io.Reader) error {
	panic("implement me")
}

func (s *ServiceImpl) Delete(ctx context.Context, user auth.User, volume int64, filename string) error {
	panic("implement me")
}

func (s *ServiceImpl) loadUser(ctx context.Context, userID string) (*store.User, error) {
	user, err := s.userStore.Get(ctx, userID)
	switch {
	case errors.Is(err, store.ErrNotFound):
		return nil, domain.NewNotFoundError(err, domain.ResourceUser, "id", userID)
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

func (s *ServiceImpl) loadUserVolume(ctx context.Context, userID string, volumeID int64) (*VolumeInfo, error) {
	if volumeID == HomeVolumeID {
		user, err := s.loadUser(ctx, userID)
		if err != nil {
			return nil, err
		}
		return &VolumeInfo{
			ID:    0,
			Label: HomeVolumeLabel,
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
	return nil, domain.NewAccessDeniedError(fmt.Sprintf("user %s", userID), fmt.Sprintf("volume %d", volumeID))
}

func fsError(err error, volumePath, filePath string) error {
	if errors.Is(err, store.ErrNotFound) {
		return domain.NewNotFoundError(err, domain.ResourceFile, "volume", volumePath, "path", filePath)
	}
	return err
}
