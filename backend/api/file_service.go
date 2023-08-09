package api

import (
	connect "connectrpc.com/connect"
	"context"
	v1 "github.com/filebrowser/filebrowser/gen/proto/filebrowser/v1"
)

type FileService struct {
}

func NewFileService() *FileService {
	return &FileService{}
}

func (f *FileService) List(ctx context.Context, c *connect.Request[v1.FileServiceListRequest]) (*connect.Response[v1.FileServiceListResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (f *FileService) Rename(ctx context.Context, c *connect.Request[v1.FileServiceRenameRequest]) (*connect.Response[v1.FileServiceRenameResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (f *FileService) Move(ctx context.Context, c *connect.Request[v1.FileServiceMoveRequest]) (*connect.Response[v1.FileServiceMoveResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (f *FileService) Copy(ctx context.Context, c *connect.Request[v1.FileServiceCopyRequest]) (*connect.Response[v1.FileServiceCopyResponse], error) {
	//TODO implement me
	panic("implement me")
}

func (f *FileService) Delete(ctx context.Context, c *connect.Request[v1.FileServiceDeleteRequest]) (*connect.Response[v1.FileServiceDeleteResponse], error) {
	//TODO implement me
	panic("implement me")
}
