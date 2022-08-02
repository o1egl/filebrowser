package osfs

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/filebrowser/filebrowser/fsys"
)

type FS struct{}

func New() *FS {
	return &FS{}
}

func (o *FS) MkDir(fPath string) error {
	return os.MkdirAll(fPath, os.ModePerm)
}

func (o *FS) Remove(fPath string) error {
	return os.RemoveAll(fPath)
}

func (o *FS) Stat(fPath string) (fsys.File, error) {
	fileInfo, err := os.Lstat(fPath)
	if err != nil {
		return fsys.File{}, err
	}

	symlink := isSymlink(fileInfo)
	// follow symlink
	if symlink && !fileInfo.IsDir() {
		fileInfo, err = os.Stat(fPath)
		if err != nil {
			return fsys.File{}, err
		}
	}

	return fsys.File{
		Name:      fileInfo.Name(),
		Size:      fileInfo.Size(),
		Type:      detectFileType(fileInfo),
		ModTime:   fileInfo.ModTime(),
		IsSymLink: symlink,
	}, nil
}

func (o *FS) List(fPath string) ([]fsys.File, error) {
	entries, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}
	files := make([]fsys.File, len(entries))
	for i, entry := range entries {
		stat, err := o.Stat(filepath.Join(fPath, entry.Name()))
		if err != nil {
			return nil, err
		}
		files[i] = stat
	}
	return files, nil
}

func (o *FS) Read(fPath string) (io.ReadCloser, error) {
	fileInfo, err := o.Stat(fPath)
	if err != nil {
		return nil, err
	}
	if fileInfo.Type != fsys.FileTypeRegular {
		return nil, fsys.ErrNotARegularFile
	}
	return os.Open(fPath)
}

func (o *FS) Write(fPath string, reader io.Reader) (written int64, err error) {
	fileInfo, err := o.Stat(fPath)
	if err != nil && !os.IsNotExist(err) {
		return 0, err
	}
	if fileInfo.Type != fsys.FileTypeRegular {
		return 0, fsys.ErrNotARegularFile
	}

	f, err := os.OpenFile(fPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.ModePerm)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return io.Copy(f, reader)
}

func (o *FS) Move(src, dst string) error {
	_, err := o.Stat(src)
	if err != nil {
		return err
	}
	_, err = o.Stat(dst)
	switch {
	case errors.Is(err, os.ErrNotExist): // do nothing
	case err != nil:
		return err
	}

	err = os.Rename(src, dst)
	if err == nil {
		return nil
	}
	// if err is not nil, try to copy and delete
	if err := o.Copy(src, dst); err != nil {
		return err
	}
	return o.Remove(src)
}

func (o *FS) Copy(src, dst string) error {
	srcInfo, err := o.Stat(src)
	if err != nil {
		return err
	}

	_, err = o.Stat(dst)
	switch {
	case errors.Is(err, os.ErrNotExist): // do nothing
		break
	case err != nil:
		return err
	default:
		return fsys.ErrExist
	}

	switch srcInfo.Type {
	case fsys.FileTypeDir:
		if err := o.MkDir(dst); err != nil {
			return err
		}
		entries, err := o.List(src)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if err := o.Copy(filepath.Join(src, entry.Name), filepath.Join(dst, entry.Name)); err != nil {
				return err
			}
		}
	case fsys.FileTypeRegular:
		reader, err := o.Read(src)
		if err != nil {
			return err
		}
		_, err = o.Write(dst, reader)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported file type: %s", srcInfo.Type)
	}

	return nil
}

func detectFileType(info os.FileInfo) fsys.FileType {
	switch {
	case info.IsDir():
		return fsys.FileTypeDir
	case isSpecialFile(info):
		return fsys.FileTypeSpecial
	default:
		return fsys.FileTypeRegular
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
