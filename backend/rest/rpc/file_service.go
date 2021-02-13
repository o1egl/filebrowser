package rpc

import (
	"context"
	"errors"
	"sort"

	"github.com/spf13/afero"
	"github.com/twitchtv/twirp"

	"github.com/filebrowser/filebrowser/v3/filesystem"
	filePb "github.com/filebrowser/filebrowser/v3/gen/proto/file/v1"
	"github.com/filebrowser/filebrowser/v3/rest"
	"github.com/filebrowser/filebrowser/v3/rest/rpc/conv"
)

type FileService struct {
	rootFS afero.Fs
}

func NewFileService(rootFS afero.Fs) *FileService {
	return &FileService{
		rootFS: rootFS,
	}
}

func (f *FileService) FileInfo(ctx context.Context, request *filePb.FileInfoRequest) (*filePb.FileInfoResponse, error) {
	user := rest.UserFromContext(ctx)
	userFs := afero.NewBasePathFs(f.rootFS, user.Scope)

	info, err := filesystem.Stat(userFs, request.Path)
	if err != nil {
		switch {
		case errors.Is(err, afero.ErrFileNotFound):
			return nil, twirp.NotFoundError("resource not found")
		default:
			return nil, twirp.InternalErrorWith(err)
		}
	}

	return &filePb.FileInfoResponse{
		Info: conv.FileInfoToPb(info),
	}, nil
}

func (f *FileService) List(ctx context.Context, request *filePb.ListRequest) (*filePb.ListResponse, error) {
	user := rest.UserFromContext(ctx)
	userFs := afero.NewBasePathFs(f.rootFS, user.Scope)

	info, err := filesystem.Stat(userFs, request.Path)
	if err != nil {
		switch {
		case errors.Is(err, afero.ErrFileNotFound):
			return nil, twirp.NotFoundError("resource not found")
		default:
			return nil, twirp.InternalErrorWith(err)
		}
	}
	var children []filesystem.Info
	if info.IsDir {
		children, err = filesystem.ReadDir(userFs, request.Path)
		if err != nil {
			return nil, twirp.InternalErrorWith(err)
		}
	}

	return &filePb.ListResponse{
		Info:     conv.FileInfoToPb(info),
		Children: sortResources(conv.FileInfosToPb(children), request.SortBy, request.SortOrder),
	}, nil
}

func sortResources(resources []*filePb.FileInfo, sortBy filePb.FileSortBy, order filePb.SortOrder) []*filePb.FileInfo {
	sort.Slice(resources, func(i, j int) bool {
		var result bool

		switch sortBy {
		case filePb.FileSortBy_FILE_SORT_BY_SIZE:
			result = resources[i].Size < resources[j].Size
		case filePb.FileSortBy_FILE_SORT_BY_MOD_TIME:
			result = resources[i].ModTime.GetSeconds() < resources[j].ModTime.GetSeconds()
		case filePb.FileSortBy_FILE_SORT_BY_NAME:
			fallthrough
		default:
			result = resources[i].Name < resources[j].Name
		}

		if order == filePb.SortOrder_SORT_ORDER_ASC {
			return !result
		}
		return result
	})
	return resources
}

func (f *FileService) CreateFile(ctx context.Context, request *filePb.CreateFileRequest) (*filePb.CreateFileResponse, error) {
	panic("implement me")
}

func (f *FileService) CreateDir(ctx context.Context, request *filePb.CreateDirRequest) (*filePb.CreateDirResponse, error) {
	panic("implement me")
}

func (f *FileService) Remove(ctx context.Context, request *filePb.RemoveRequest) (*filePb.RemoveResponse, error) {
	panic("implement me")
}
