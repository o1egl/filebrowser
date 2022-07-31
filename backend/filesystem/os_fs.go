package filesystem

import (
	"io"
	"os"
)

type OSFS struct{}

func (o *OSFS) MkDir(fPath string) error {
	return os.MkdirAll(fPath, os.ModePerm)
}

func (o *OSFS) Remove(fPath string) error {
	return os.RemoveAll(fPath)
}

func (o *OSFS) Stat(fPath string) (File, error) {
	fileInfo, err := os.Stat(fPath)
	if err != nil {
		return File{}, err
	}
	return File{
		Name:      fileInfo.Name(),
		Size:      fileInfo.Size(),
		Type:      detectFileType(fileInfo),
		ModTime:   fileInfo.ModTime(),
		IsSymLink: isSymlink(fileInfo),
	}, nil
}

func (o *OSFS) List(fPath string) ([]File, error) {
	entries, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}
	files := make([]File, len(entries))
	for i, entry := range entries {
		fileInfo, err := entry.Info()
		if err != nil {
			return nil, err
		}
		files[i] = File{
			Name:      fileInfo.Name(),
			Size:      fileInfo.Size(),
			Type:      detectFileType(fileInfo),
			ModTime:   fileInfo.ModTime(),
			IsSymLink: isSymlink(fileInfo),
		}
	}
	return files, nil
}

func (o *OSFS) Read(fPath string) (io.ReadCloser, error) {
	return os.Open(fPath)
}

func (o *OSFS) Write(fPath string, reader io.Reader) error {
	f, err := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, reader)
	return err
}

func (o *OSFS) Move(src, dst string) error {
	return os.Rename(src, dst)
}

func detectFileType(info os.FileInfo) FileType {
	switch {
	case info.IsDir():
		return FileTypeDir
	case isSpecialFile(info):
		return FileTypeSpecial
	default:
		return FileTypeRegular
	}
}

// isSpecialFile reports if this file is a special file such as a named pipe,
// device file, or socket. If so it will return a ErrSpecialFile.
func isSpecialFile(fi os.FileInfo) bool {
	if (fi.Mode()&os.ModeDevice) == os.ModeDevice ||
		(fi.Mode()&os.ModeNamedPipe) == os.ModeNamedPipe ||
		(fi.Mode()&os.ModeSocket) == os.ModeSocket ||
		(fi.Mode()&os.ModeCharDevice) == os.ModeCharDevice {

		return true
	}

	return false
}

// isSymlink reports if this file is a symbolic link.
func isSymlink(fi os.FileInfo) bool {
	return (fi.Mode() & os.ModeSymlink) == os.ModeSymlink
}
