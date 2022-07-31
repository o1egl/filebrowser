package conv

import (
	"github.com/AlekSi/pointer"

	"github.com/filebrowser/filebrowser/domain"
	"github.com/filebrowser/filebrowser/storage/sql/model"
)

func UserToDomain(source *model.User) *domain.User {
	return &domain.User{
		ID:       source.ID,
		Username: source.Username,
		Password: pointer.GetString(source.Password),
		Home: domain.HomeVolume{
			Path: source.Home,
			Permissions: domain.VolumePermissions{
				Read:   source.Read,
				Create: source.Create,
				Modify: source.Modify,
				Delete: source.Delete,
				Share:  source.Share,
			},
		},
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
		Home:         source.Home.Path,
		Read:         source.Home.Permissions.Read,
		Create:       source.Home.Permissions.Create,
		Modify:       source.Home.Permissions.Modify,
		Delete:       source.Home.Permissions.Delete,
		Share:        source.Home.Permissions.Share,
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
