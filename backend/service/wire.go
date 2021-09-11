package service

import (
	"github.com/filebrowser/filebrowser/v3/service/filebrowser"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	filebrowser.New,
	wire.Bind(new(filebrowser.Service), new(*filebrowser.ServiceImpl)),
)
