package rpc

import (
	"context"

	pb "github.com/filebrowser/filebrowser/v3/gen/proto/filebrowser/v1"
	"github.com/filebrowser/filebrowser/v3/service/filebrowser"
)

type FileService struct {
	fbSvc filebrowser.Service
}

func NewFileService(fbSvc filebrowser.Service) *FileService {
	return &FileService{fbSvc: fbSvc}
}

func (f FileService) List(ctx context.Context, request *pb.FileServiceListRequest) (*pb.FileServiceListResponse, error) {
	panic("implement me")
}
