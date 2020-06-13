//go:generate go-enum --sql --marshal --lower --names --file $GOFILE
package filesystem

import (
	"io"
	"mime"
	"net/http"
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
)
*/
type Type int

func Stat(fs afero.Fs, fPath string) (Info, error) {
	fileInfo, err := fs.Stat(fPath)
	if err != nil {
		return Info{}, err
	}
	fileType := TypeDir
	if !fileInfo.IsDir() {
		fileType = detectFileType(fs, fPath)
	}
	info := Info{
		Path:      fPath,
		Name:      filepath.Base(fPath),
		Size:      fileInfo.Size(),
		Extension: filepath.Ext(fileInfo.Name()),
		ModTime:   fileInfo.ModTime(),
		Mode:      fileInfo.Mode(),
		Type:      fileType,
		IsDir:     fileType == TypeDir,
	}

	return info, nil
}

func ReadDir(fs afero.Fs, fPath string) ([]Info, error) {
	fileInfos, err := afero.ReadDir(fs, fPath)
	if err != nil {
		return nil, err
	}

	infos := make([]Info, 0, len(fileInfos))
	for _, fileInfo := range fileInfos {
		filePath := filepath.Join(fPath, fileInfo.Name())
		fileType := TypeDir
		if !fileInfo.IsDir() {
			fileType = detectFileType(fs, filePath)
		}
		infos = append(infos, Info{
			Path:      filePath,
			Name:      fileInfo.Name(),
			Size:      fileInfo.Size(),
			Extension: filepath.Ext(fileInfo.Name()),
			ModTime:   fileInfo.ModTime(),
			Mode:      fileInfo.Mode(),
			Type:      fileType,
			IsDir:     fileType == TypeDir,
		})
	}

	return infos, nil
}

// failing to detect the type should not return error.
// imagine the situation where a file in a dir with thousands
// of files couldn't be opened: we'd have immediately
// a 500 even though it doesn't matter. So we just log it.
func detectFileType(fs afero.Fs, fileName string) Type {
	header, n, err := readHeader(fs, fileName)
	if err != nil {
		return TypeBlob
	}

	mimetype := mime.TypeByExtension(filepath.Ext(fileName))
	if mimetype == "" {
		mimetype = http.DetectContentType(header[:n])
	}

	switch {
	case strings.HasPrefix(mimetype, "video"):
		return TypeVideo
	case strings.HasPrefix(mimetype, "audio"):
		return TypeAudio
	case strings.HasPrefix(mimetype, "image"):
		return TypeImage
	case strings.HasPrefix(mimetype, "text") || isText(header[:n]):
		return TypeText
	default:
		return TypeBlob
	}
}

func readHeader(fs afero.Fs, fPath string) (b []byte, n int, err error) {
	reader, err := fs.Open(fPath)
	if err != nil {
		return nil, 0, err
	}
	defer reader.Close()

	b = make([]byte, 512)
	n, err = reader.Read(b)
	if err != nil && err != io.EOF {
		return nil, 0, nil
	}
	return b, n, nil
}
