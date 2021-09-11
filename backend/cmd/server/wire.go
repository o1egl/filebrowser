//+build wireinject

package server

import (
	"context"
	"path/filepath"

	"github.com/filebrowser/filebrowser/v3/auth"
	"github.com/filebrowser/filebrowser/v3/config"
	"github.com/filebrowser/filebrowser/v3/domain"
	"github.com/filebrowser/filebrowser/v3/hash"
	"github.com/filebrowser/filebrowser/v3/rest/api"
	"github.com/filebrowser/filebrowser/v3/service"
	"github.com/filebrowser/filebrowser/v3/store/sql"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func InitializeServer(ctx context.Context, cfg *config.Config, version domain.Version) (*app, error) {
	wire.Build(
		rootFSProvider,
		secretProvider,
		auth.Set,
		hash.Set,
		sql.Set,
		service.Set,
		api.NewServer,
		newApp,
	)
	return &app{}, nil
}

func secretProvider(cfg *config.Config) domain.Secret {
	return domain.Secret(cfg.Secret)
}

func rootFSProvider(cfg *config.Config) (afero.Fs, error) {
	absRootPath, err := filepath.Abs(cfg.RootPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get abs path")
	}
	return afero.NewBasePathFs(afero.NewOsFs(), absRootPath), nil
}
