package rpc

import (
	pb "github.com/filebrowser/filebrowser/v3/gen/proto/filebrowser/v1"
	"github.com/google/wire"
)

var Set = wire.NewSet(
	NewFileService,
	wire.Bind(new(pb.FileService), new(*FileService)),
)
