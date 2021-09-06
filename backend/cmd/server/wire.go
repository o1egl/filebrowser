//+build wireinject

package server

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/filebrowser/filebrowser/v3/auth"
	"github.com/filebrowser/filebrowser/v3/log"
	"github.com/filebrowser/filebrowser/v3/rest/api"
	"github.com/filebrowser/filebrowser/v3/service"
	"github.com/filebrowser/filebrowser/v3/service/filebrowser"
	"github.com/filebrowser/filebrowser/v3/store"
	"github.com/filebrowser/filebrowser/v3/store/sql"
	"github.com/filebrowser/filebrowser/v3/store/sql/ent"
	"github.com/filebrowser/filebrowser/v3/token"
	pkgAuth "github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	authToken "github.com/go-pkgz/auth/token"
	"github.com/google/wire"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func EntClientProvider(ctx context.Context, srvCmd *ServerCommand) (client *ent.Client, err error) {
	switch srvCmd.Store.Type {
	case "sqlite":
		if err = makeDirs(filepath.Dir(srvCmd.Store.SQLite.File)); err != nil {
			return nil, errors.Wrap(err, "failed to create sqlite store")
		}
		client, err = ent.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&mode=rwc&_fk=1", srvCmd.Store.SQLite.File))
	case "postgres":
		client, err = ent.Open("postgres", srvCmd.Store.Postgres.DSN)
	case "mysql":
		client, err = ent.Open("mysql", srvCmd.Store.Postgres.DSN)
	default:
		return nil, errors.Errorf("unsupported store type %s", srvCmd.Store.Type)
	}
	if err != nil {
		return nil, errors.Wrap(err, "can't initialize data store")
	}

	log.WithContext(ctx).Debugf("Apply schema migrations")
	if err := client.Schema.Create(ctx); err != nil {
		return nil, err
	}

	return client, nil
}

var DataStoreSet = wire.NewSet(
	EntClientProvider,
	sql.NewUserStore,
	wire.Bind(new(store.UserStore), new(*sql.UserStore)),
	sql.NewVolumeStore,
	wire.Bind(new(store.VolumeStore), new(*sql.VolumeStore)),
)

var ServiceSet = wire.NewSet(
	filebrowser.New,
	wire.Bind(new(service.FileBrowser), new(*filebrowser.Service)),
)

func AuthServiceProvider(srvCmd *ServerCommand, userStore store.UserStore) *auth.Service {
	return auth.NewService(userStore, srvCmd.Auth.User.Home, srvCmd.Locale)
}

func TokenServiceProvider(srvCmd *ServerCommand) *token.Service {
	return token.New(srvCmd.Secret)
}

func ApiServerOptionsProvider(srvCmd *ServerCommand, sslConfig api.SSLConfig) api.Options {
	return api.Options{
		Host:         srvCmd.Host,
		Port:         srvCmd.Port,
		ServerURL:    srvCmd.ServerURL,
		SharedSecret: srvCmd.Secret,
		Revision:     srvCmd.Revision,
		AccessLog:    srvCmd.AccessLog.Enable,
		Anonymous:    srvCmd.Auth.Anonymous.Enable,
		SSLConfig:    sslConfig,
	}
}

func RootFSProvider(srvCmd *ServerCommand) afero.Fs {
	return afero.NewBasePathFs(afero.NewOsFs(), srvCmd.RootPath)
}

func AuthenticatorProvider(
	srvCmd *ServerCommand,
	authService *auth.Service,
	authRefreshCache *authRefreshCache, //nolint:interfacer
) (*pkgAuth.Service, error) { //nolint:unparam
	authenticator := pkgAuth.NewService(pkgAuth.Opts{
		DisableXSRF:    true, // TODO remove it
		URL:            strings.TrimSuffix(srvCmd.ServerURL, "/"),
		Issuer:         "File Browser",
		TokenDuration:  srvCmd.Auth.TTL.JWT,
		CookieDuration: srvCmd.Auth.TTL.Cookie,
		SecureCookies:  strings.HasPrefix(srvCmd.ServerURL, "https://"),
		SecretReader: authToken.SecretFunc(func(aud string) (string, error) {
			return srvCmd.Secret, nil
		}),
		ClaimsUpd:        authService,
		BasicAuthChecker: authService.BasicAuthChecker,
		Validator:        authService,
		AvatarStore:      avatar.NewNoOp(),
		Logger:           log.NewLogrAdapter(log.DefaultLogger),
		RefreshCache:     authRefreshCache,
	})

	addAuthProviders(srvCmd, authenticator, authService)

	return authenticator, nil
}

func InitializeServer(ctx context.Context, srvCmd *ServerCommand) (*serverApp, error) {
	wire.Build(
		NewServerApp,
		NewSSLConfig,
		AuthServiceProvider,
		newAuthRefreshCache,
		AuthenticatorProvider,
		TokenServiceProvider,
		ApiServerOptionsProvider,
		RootFSProvider,
		api.NewServer,
		DataStoreSet,
		ServiceSet,
	)
	return &serverApp{}, nil
}
