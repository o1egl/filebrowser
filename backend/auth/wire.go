package auth

import (
	"github.com/go-pkgz/auth/middleware"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewService,
	NewAuthenticator,
	NewInMemoryAuthRefreshCache,
	wire.Bind(new(middleware.RefreshCache), new(*InMemoryAuthRefreshCache)),
)
