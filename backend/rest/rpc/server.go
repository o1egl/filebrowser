package rpc

import (
	"github.com/spf13/afero"
	"github.com/twitchtv/twirp"

	filePb "github.com/filebrowser/filebrowser/v3/gen/proto/file/v1"
	userPb "github.com/filebrowser/filebrowser/v3/gen/proto/user/v1"
)

func NewFileServiceServer(pathPrefix string, root afero.Fs) filePb.TwirpServer {
	return filePb.NewFileServiceServer(
		filePb.NewFileServiceWithTwirpValidation(NewFileService(root)),
		twirp.WithServerPathPrefix(pathPrefix),
	)
}

func NewUserServiceServer(pathPrefix string) userPb.TwirpServer {
	return userPb.NewUserServiceServer(
		userPb.NewUserServiceWithTwirpValidation(NewUserService()),
		twirp.WithServerPathPrefix(pathPrefix),
	)
}
