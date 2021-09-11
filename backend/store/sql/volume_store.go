package sql

import (
	"context"

	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/sql/conv"
	"github.com/filebrowser/filebrowser/v3/store/sql/conv/generated"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent/group"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent/user"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent/volume"
)

type VolumeStore struct {
	client    *ent.Client
	converter conv.VolumeConverter
}

func NewVolumeStore(client *ent.Client) *VolumeStore {
	return &VolumeStore{client: client, converter: &generated.VolumeConverterImpl{}}
}

func (v *VolumeStore) GetUserVolumes(ctx context.Context, userID string) ([]*store.Volume, error) {
	volumes, err := v.client.Volume.Query().Where(
		volume.Or(
			volume.HasUsersWith(user.ID(userID)),
			volume.HasGroupsWith(group.HasUsersWith(user.ID(userID))),
		),
	).All(ctx)
	switch {
	case ent.IsNotFound(err):
		return nil, store.ErrNotFound
	case err != nil:
		return nil, err
	}
	return v.converter.ConvertSlice(volumes), nil
}
