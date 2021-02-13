package conv

import (
	"github.com/golang/protobuf/ptypes"

	"github.com/filebrowser/filebrowser/v3/filesystem"
	filePb "github.com/filebrowser/filebrowser/v3/gen/proto/file/v1"
)

func FileTypeToPb(fileType filesystem.Type) filePb.FileType {
	switch fileType {
	case filesystem.TypeBlob:
		return filePb.FileType_FILE_TYPE_BLOB
	case filesystem.TypeVideo:
		return filePb.FileType_FILE_TYPE_VIDEO
	case filesystem.TypeAudio:
		return filePb.FileType_FILE_TYPE_AUDIO
	case filesystem.TypeImage:
		return filePb.FileType_FILE_TYPE_IMAGE
	case filesystem.TypeText:
		return filePb.FileType_FILE_TYPE_TEXT
	case filesystem.TypeDir:
		return filePb.FileType_FILE_TYPE_DIR
	case filesystem.TypeSpecial:
		return filePb.FileType_FILE_TYPE_SPECIAL
	default:
		return filePb.FileType_FILE_TYPE_UNSPECIFIED
	}
}

func FileInfoToPb(info filesystem.Info) *filePb.FileInfo {
	ts, _ := ptypes.TimestampProto(info.ModTime)
	return &filePb.FileInfo{
		Path:    info.Path,
		Name:    info.Name,
		Size:    info.Size,
		Type:    FileTypeToPb(info.Type),
		ModTime: ts,
		Mode:    uint32(info.Mode),
	}
}

func FileInfosToPb(infos []filesystem.Info) []*filePb.FileInfo {
	resp := make([]*filePb.FileInfo, len(infos))
	for i, info := range infos {
		resp[i] = FileInfoToPb(info)
	}
	return resp
}
