//+build wireinject

package server

import (
	"context"

	"github.com/google/wire"
)

func InitializeServer(ctx context.Context, srvCmd *ServerCommand) (*serverApp, error) {
	wire.Build(
		NewServerApp,
		RootFSProvider,
		AuthenticatorSet,
		ApiServerSet,
		DataStoreSet,
		ServiceSet,
	)
	return &serverApp{}, nil
}
