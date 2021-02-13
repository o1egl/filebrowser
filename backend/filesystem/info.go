//go:generate go-enum --sql --marshal --lower --names --file $GOFILE
package filesystem

import (
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/afero"
)

type Info struct {
	Path      string      `json:"path"`
	Name      string      `json:"name"`
	Size      int64       `json:"size"`
	Extension string      `json:"extension,omitempty"`
	ModTime   time.Time   `json:"modified"`
	Mode      os.FileMode `json:"mode"`
	Type      Type        `json:"type"`
	IsSymlink bool        `json:"isSymlink"`
	IsDir     bool        `json:"isDir"`
}

/*
ENUM(
blob
video
audio
image
text
dir
special
)
*/
type Type int

func toRelative(fPath string) string {
	fPath = strings.TrimPrefix(fPath, "/")
	if fPath == "" {
		fPath = "."
	}
	return fPath
}

func Stat(fSys afero.Fs, fPath string) (Info, error) {
	fileInfo, err := fSys.Stat(fPath)
	if err != nil {
		return Info{}, err
	}
	fileType := detectFileType(fileInfo)
	info := Info{
		Path:      fPath,
		Name:      filepath.Base(fPath),
		Size:      fileInfo.Size(),
		Extension: filepath.Ext(fileInfo.Name()),
		ModTime:   fileInfo.ModTime(),
		Mode:      fileInfo.Mode(),
		Type:      fileType,
		IsSymlink: IsSymlink(fileInfo),
		IsDir:     fileType == TypeDir,
	}

	return info, nil
}

// IsSpecialFile reports if this file is a special file such as a named pipe,
// device file, or socket. If so it will return a ErrSpecialFile.
func IsSpecialFile(fi os.FileInfo) bool {
	if (fi.Mode()&os.ModeDevice) == os.ModeDevice ||
		(fi.Mode()&os.ModeNamedPipe) == os.ModeNamedPipe ||
		(fi.Mode()&os.ModeSocket) == os.ModeSocket ||
		(fi.Mode()&os.ModeCharDevice) == os.ModeCharDevice {

		return true
	}

	return false
}

// IsSymlink reports if this file is a symbolic link.
func IsSymlink(fi os.FileInfo) bool {
	return (fi.Mode() & os.ModeSymlink) == os.ModeSymlink
}

func ReadDir(fSys afero.Fs, dirPath string) ([]Info, error) {
	entries, err := afero.ReadDir(fSys, dirPath)
	if err != nil {
		return nil, err
	}

	infos := make([]Info, 0, len(entries))
	for _, entry := range entries {
		filePath := filepath.Join(dirPath, entry.Name())
		fileType := detectFileType(entry)
		infos = append(infos, Info{
			Path:      filePath,
			Name:      entry.Name(),
			Size:      entry.Size(),
			Extension: filepath.Ext(entry.Name()),
			ModTime:   entry.ModTime(),
			Mode:      entry.Mode(),
			Type:      fileType,
			IsSymlink: IsSymlink(entry),
			IsDir:     fileType == TypeDir,
		})
	}

	return infos, nil
}

func detectFileType(info os.FileInfo) Type {
	mimetype := mime.TypeByExtension(filepath.Ext(info.Name()))
	switch {
	case info.IsDir():
		return TypeDir
	case IsSpecialFile(info):
		return TypeSpecial
	case strings.HasPrefix(mimetype, "video"):
		return TypeVideo
	case strings.HasPrefix(mimetype, "audio"):
		return TypeAudio
	case strings.HasPrefix(mimetype, "image"):
		return TypeImage
	case strings.HasPrefix(mimetype, "text"):
		return TypeText
	default:
		return TypeBlob
	}
}
