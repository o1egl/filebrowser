package auth

import (
	"context"
	"strings"

	"github.com/filebrowser/filebrowser/v3/config"
	"github.com/filebrowser/filebrowser/v3/log"
	pkgAuth "github.com/go-pkgz/auth"
	"github.com/go-pkgz/auth/avatar"
	"github.com/go-pkgz/auth/middleware"
	"github.com/go-pkgz/auth/provider"
	authToken "github.com/go-pkgz/auth/token"
)

func NewAuthenticator(
	cfg *config.Config,
	authService *Service,
	authRefreshCache middleware.RefreshCache,
) (*pkgAuth.Service, error) {
	authenticator := pkgAuth.NewService(pkgAuth.Opts{
		DisableXSRF:    true, // TODO remove it
		URL:            strings.TrimSuffix(cfg.Server.URL, "/"),
		Issuer:         "File Browser",
		TokenDuration:  cfg.Auth.TTL.JWT,
		CookieDuration: cfg.Auth.TTL.Cookie,
		SecureCookies:  strings.HasPrefix(cfg.Server.URL, "https://"),
		SecretReader: authToken.SecretFunc(func(aud string) (string, error) {
			return cfg.Secret, nil
		}),
		ClaimsUpd:        authService,
		BasicAuthChecker: authService.BasicAuthChecker,
		Validator:        authService,
		AvatarStore:      avatar.NewNoOp(),
		Logger:           log.NewLogrAdapter(log.DefaultLogger),
		RefreshCache:     authRefreshCache,
	})

	addAuthProviders(cfg.Auth, authenticator, authService)

	return authenticator, nil
}

func addAuthProviders(cfg config.Auth, authenticator *pkgAuth.Service, localProvider provider.CredChecker) {
	providers := 0

	providers++
	authenticator.AddDirectProvider("local", localProvider)

	if cfg.Google.CID != "" && cfg.Google.CSEC != "" {
		authenticator.AddProvider("google", cfg.Google.CID, cfg.Google.CSEC)
		providers++
	}
	if cfg.Github.CID != "" && cfg.Github.CSEC != "" {
		authenticator.AddProvider("github", cfg.Github.CID, cfg.Github.CSEC)
		providers++
	}
	if cfg.Facebook.CID != "" && cfg.Facebook.CSEC != "" {
		authenticator.AddProvider("facebook", cfg.Facebook.CID, cfg.Facebook.CSEC)
		providers++
	}
	if cfg.Twitter.CID != "" && cfg.Twitter.CSEC != "" {
		authenticator.AddProvider("twitter", cfg.Twitter.CID, cfg.Twitter.CSEC)
		providers++
	}

	if cfg.Dev {
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
