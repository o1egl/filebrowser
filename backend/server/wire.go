package server

import (
	"github.com/filebrowser/filebrowser/v3/server/rpc"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	rpc.Set,
	NewServer,
)
