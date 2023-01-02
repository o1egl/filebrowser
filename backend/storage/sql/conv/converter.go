package conv

import (
	"github.com/AlekSi/pointer"

	"github.com/filebrowser/filebrowser/domain"
	"github.com/filebrowser/filebrowser/storage/sql/model"
)

func UserToDomain(source *model.User) *domain.User {
	return &domain.User{
		ID:           source.ID,
		Username:     source.Username,
		Password:     pointer.GetString(source.Password),
		Name:         source.Name,
		Locale:       source.Locale,
		LockPassword: source.LockPassword,
		Blocked:      source.Blocked,
	}
}

func UsersToDomain(source []*model.User) []*domain.User {
	res := make([]*domain.User, len(source))
	for i, obj := range source {
		res[i] = UserToDomain(obj)
	}
	return res
}

func UserFromDomain(source *domain.User) *model.User {
	return &model.User{
		ID:           source.ID,
		Username:     source.Username,
		Password:     pointer.ToStringOrNil(source.Password),
		Name:         source.Name,
		Locale:       source.Locale,
		LockPassword: source.LockPassword,
		Blocked:      source.Blocked,
		IsAdmin:      source.IsAdmin,
	}
}

func VolumeToDomain(source *model.Volume) *domain.Volume {
	return &domain.Volume{
		ID:          source.ID,
		Label:       source.Label,
		Path:        source.Path,
		Description: source.Description,
	}
}

func VolumesToDomain(source []*model.Volume) []*domain.Volume {
	res := make([]*domain.Volume, len(source))
	for i, obj := range source {
		res[i] = VolumeToDomain(obj)
	}
	return res
}

func UploadToDomain(source *model.Upload) *domain.Upload {
	return &domain.Upload{
		ID:        source.ID,
		UserID:    source.UserID,
		VolumeID:  source.VolumeID,
		Path:      source.Path,
		TmpPath:   source.TmpPath,
		Size:      source.Size,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
}

func UploadFromDomain(source *domain.Upload) *model.Upload {
	return &model.Upload{
		ID:        source.ID,
		UserID:    source.UserID,
		VolumeID:  source.VolumeID,
		Path:      source.Path,
		TmpPath:   source.TmpPath,
		Size:      source.Size,
		CreatedAt: source.CreatedAt,
		UpdatedAt: source.UpdatedAt,
	}
}
