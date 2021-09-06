package server

import (
	"context"
	"strings"

	"github.com/filebrowser/filebrowser/v3/auth"
	"github.com/filebrowser/filebrowser/v3/log"
	pkgAuth "github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/middleware"
	"github.com/go-pkgz/auth/provider"
	authToken "github.com/go-pkgz/auth/token"
	"github.com/google/wire"
)

var AuthenticatorSet = wire.NewSet(
	AuthenticatorProvider,
	auth.NewInMemoryAuthRefreshCache,
	wire.Bind(new(middleware.RefreshCache), new(*auth.InMemoryAuthRefreshCache)),
)

func AuthenticatorProvider(
	srvCmd *ServerCommand,
	authService *auth.Service,
	authRefreshCache middleware.RefreshCache,
) (*pkgAuth.Service, error) {
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

func addAuthProviders(srvCmd *ServerCommand, authenticator *pkgAuth.Service, localProvider provider.CredChecker) {
	providers := 0

	providers++
	authenticator.AddDirectProvider("local", localProvider)

	if srvCmd.Auth.Google.CID != "" && srvCmd.Auth.Google.CSEC != "" {
		authenticator.AddProvider("google", srvCmd.Auth.Google.CID, srvCmd.Auth.Google.CSEC)
		providers++
	}
	if srvCmd.Auth.Github.CID != "" && srvCmd.Auth.Github.CSEC != "" {
		authenticator.AddProvider("github", srvCmd.Auth.Github.CID, srvCmd.Auth.Github.CSEC)
		providers++
	}
	if srvCmd.Auth.Facebook.CID != "" && srvCmd.Auth.Facebook.CSEC != "" {
		authenticator.AddProvider("facebook", srvCmd.Auth.Facebook.CID, srvCmd.Auth.Facebook.CSEC)
		providers++
	}
	if srvCmd.Auth.Twitter.CID != "" && srvCmd.Auth.Twitter.CSEC != "" {
		authenticator.AddProvider("twitter", srvCmd.Auth.Twitter.CID, srvCmd.Auth.Twitter.CSEC)
		providers++
	}

	if srvCmd.Auth.Dev {
		log.Warnf("dev oauth provider is enabled")
		authenticator.AddProvider("dev", "", "")
		providers++

		// run dev/test oauth2 server on :8084
		go func() {
			devAuthServer, err := authenticator.DevAuth() // peak dev oauth2 server
			if err != nil {
				log.Fatalf("failed to start dev oauth2 server, %v", err)
			}
			devAuthServer.Run(context.Background())
		}()
	}

	if providers == 0 {
		log.Warnf("no auth providers defined")
	}
}
