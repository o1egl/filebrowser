package server

import (
	"github.com/filebrowser/filebrowser/v3/auth"
	"github.com/filebrowser/filebrowser/v3/service"
	"github.com/filebrowser/filebrowser/v3/service/filebrowser"
	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/token"
	"github.com/google/wire"
	"github.com/spf13/afero"
)

var ServiceSet = wire.NewSet(
	filebrowser.New,
	wire.Bind(new(service.FileBrowser), new(*filebrowser.Service)),
	TokenServiceProvider,
	AuthServiceProvider,
)

func TokenServiceProvider(srvCmd *ServerCommand) *token.Service {
	return token.New(srvCmd.Secret)
}

func AuthServiceProvider(srvCmd *ServerCommand, userStore store.UserStore) *auth.Service {
	return auth.NewService(userStore, srvCmd.Auth.User.Home, srvCmd.Locale)
}

func RootFSProvider(srvCmd *ServerCommand) afero.Fs {
	return afero.NewBasePathFs(afero.NewOsFs(), srvCmd.RootPath)
}
