package osfs

import (
	"bytes"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/filebrowser/filebrowser/fsys"
)

func TestOSFS(t *testing.T) {
	osfs := New()
	tmp := createTmp(t)
	folder := filepath.Join(tmp, "testdir", "subfolder")
	err := osfs.MkDir(folder)
	require.NoError(t, err)

	fileStat, err := os.Stat(folder)
	require.NoError(t, err)
	assert.True(t, fileStat.IsDir())
}

func TestOSFS_Remove(t *testing.T) {
	osfs := New()

	t.Run("should remove a file", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "testfile")
		f, err := os.Create(filename)
		require.NoError(t, err)
		err = f.Close()
		require.NoError(t, err)

		err = osfs.Remove(filename)
		require.NoError(t, err)

		_, err = os.Stat(filename)
		assert.True(t, os.IsNotExist(err))
	})
	t.Run("should remove a folder", func(t *testing.T) {
		tmp := createTmp(t)
		folder := filepath.Join(tmp, "testdir")
		err := osfs.MkDir(folder)
		require.NoError(t, err)

		err = osfs.Remove(folder)
		require.NoError(t, err)

		_, err = os.Stat(folder)
		assert.True(t, os.IsNotExist(err))
	})
	t.Run("should remove a folder with its content", func(t *testing.T) {
		tmp := createTmp(t)
		folder := filepath.Join(tmp, "testdir")
		err := osfs.MkDir(folder)
		require.NoError(t, err)

		filename := filepath.Join(folder, "testfile")
		f, err := os.Create(filename)
		require.NoError(t, err)
		err = f.Close()
		require.NoError(t, err)

		err = osfs.Remove(folder)
		require.NoError(t, err)

		_, err = os.Stat(filename)
		assert.True(t, os.IsNotExist(err))
		_, err = os.Stat(folder)
		assert.True(t, os.IsNotExist(err))
	})
}

func createTmp(t *testing.T) string {
	temp, err := os.MkdirTemp("", "filebrowser-*")
	require.NoError(t, err)

	t.Cleanup(func() {
		err = os.RemoveAll(temp)
		require.NoError(t, err)
	})
	return temp
}

func TestOSFS_Stat(t *testing.T) {
	osfs := New()
	t.Run("should return a regular file stat", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "testfile")
		err := os.WriteFile(filename, []byte("test"), 0644)
		require.NoError(t, err)

		fileStat, err := osfs.Stat(filename)
		require.NoError(t, err)
		assert.Equal(t, filepath.Base(filename), fileStat.Name)
		assert.Equal(t, int64(4), fileStat.Size)
		assert.Equal(t, fsys.FileTypeRegular, fileStat.Type)
		assert.Equal(t, false, fileStat.IsSymLink)
	})
	t.Run("should return a symlink stat", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "testfile")
		err := os.WriteFile(filename, []byte("test"), 0644)
		require.NoError(t, err)

		link := filepath.Join(tmp, "testlink")
		err = os.Symlink(filename, link)
		require.NoError(t, err)

		fileStat, err := osfs.Stat(link)
		require.NoError(t, err)
		assert.Equal(t, filepath.Base(link), fileStat.Name)
		assert.Equal(t, int64(4), fileStat.Size)
		assert.Equal(t, fsys.FileTypeRegular, fileStat.Type)
		assert.Equal(t, true, fileStat.IsSymLink)
	})
	t.Run("should return a folder stat", func(t *testing.T) {
		tmp := createTmp(t)
		folder := filepath.Join(tmp, "testdir")
		err := osfs.MkDir(folder)
		require.NoError(t, err)

		fileStat, err := osfs.Stat(folder)
		require.NoError(t, err)
		assert.Equal(t, filepath.Base(folder), fileStat.Name)
		assert.Equal(t, fsys.FileTypeDir, fileStat.Type)
		assert.Equal(t, false, fileStat.IsSymLink)
	})
	t.Run("should return a special file stat", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "test.sock")
		l, err := net.Listen("unix", filename)
		require.NoError(t, err)
		t.Cleanup(func() {
			err := l.Close()
			require.NoError(t, err)
		})

		fileStat, err := osfs.Stat(filename)
		require.NoError(t, err)
		assert.Equal(t, filepath.Base(filename), fileStat.Name)
		assert.Equal(t, fsys.FileTypeSpecial, fileStat.Type)
		assert.Equal(t, false, fileStat.IsSymLink)
	})
	t.Run("should return an error if the file does not exist", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "testfile")
		_, err := osfs.Stat(filename)
		assert.ErrorIs(t, err, fsys.ErrNotExist)
	})
}

func TestOSFS_List(t *testing.T) {
	osfs := New()
	t.Run("should return a list of files", func(t *testing.T) {
		tmp := createTmp(t)
		folder := filepath.Join(tmp, "testdir")
		err := osfs.MkDir(folder)
		require.NoError(t, err)

		filename := filepath.Join(tmp, "testfile")
		err = os.WriteFile(filename, []byte("test"), 0644)
		require.NoError(t, err)

		link := filepath.Join(tmp, "testfilelink")
		err = os.Symlink(filename, link)
		require.NoError(t, err)

		files, err := osfs.List(tmp)
		require.NoError(t, err)
		require.Equal(t, 3, len(files))

		assert.Equal(t, filepath.Base(folder), files[0].Name)
		assert.Equal(t, fsys.FileTypeDir, files[0].Type)
		assert.Equal(t, false, files[0].IsSymLink)

		assert.Equal(t, filepath.Base(filename), files[1].Name)
		assert.Equal(t, int64(4), files[1].Size)
		assert.Equal(t, fsys.FileTypeRegular, files[1].Type)
		assert.Equal(t, false, files[1].IsSymLink)

		assert.Equal(t, filepath.Base(link), files[2].Name)
		assert.Equal(t, int64(4), files[2].Size)
		assert.Equal(t, fsys.FileTypeRegular, files[2].Type)
		assert.Equal(t, true, files[2].IsSymLink)
	})
	t.Run("should return an error if the folder does not exist", func(t *testing.T) {
		tmp := createTmp(t)
		_, err := osfs.List(filepath.Join(tmp, "testdir"))
		assert.ErrorIs(t, err, fsys.ErrNotExist)
	})
}

func TestOSFS_Read(t *testing.T) {
	osfs := New()
	t.Run("should read a file", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "testfile")
		err := os.WriteFile(filename, []byte("test"), 0644)
		require.NoError(t, err)

		file, err := osfs.Read(filename)
		require.NoError(t, err)

		data, err := io.ReadAll(file)
		require.NoError(t, err)

		assert.Equal(t, []byte("test"), data)

		err = file.Close()
		require.NoError(t, err)
	})
	t.Run("should return an error if the file does not exist", func(t *testing.T) {
		tmp := createTmp(t)
		_, err := osfs.Read(filepath.Join(tmp, "testfile"))
		assert.ErrorIs(t, err, fsys.ErrNotExist)
	})
	t.Run("should return an error if the file is a directory", func(t *testing.T) {
		tmp := createTmp(t)
		err := osfs.MkDir(filepath.Join(tmp, "testdir"))
		require.NoError(t, err)

		_, err = osfs.Read(filepath.Join(tmp, "testdir"))
		assert.ErrorIs(t, err, fsys.ErrNotARegularFile)
	})
}

func TestOSFS_Write(t *testing.T) {
	osfs := New()
	t.Run("should write a file", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "testfile")
		fileContent := []byte("test")
		written, err := osfs.Write(filename, bytes.NewReader(fileContent))
		require.NoError(t, err)
		assert.Equal(t, int64(len(fileContent)), written)

		data, err := ioutil.ReadFile(filename)
		require.NoError(t, err)

		assert.Equal(t, []byte("test"), data)
	})
	t.Run("should append to existing file", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "testfile")

		hello := []byte("hello ")
		written, err := osfs.Write(filename, bytes.NewReader(hello))
		require.NoError(t, err)
		assert.Equal(t, int64(len(hello)), written)

		world := []byte("world")
		written, err = osfs.Write(filename, bytes.NewReader(world))
		require.NoError(t, err)
		assert.Equal(t, int64(len(world)), written)

		data, err := ioutil.ReadFile(filename)
		require.NoError(t, err)

		assert.Equal(t, append(hello, world...), data)
	})
	t.Run("should return an error if the file is not a regular file", func(t *testing.T) {
		tmp := createTmp(t)
		err := osfs.MkDir(filepath.Join(tmp, "testdir"))
		require.NoError(t, err)

		_, err = osfs.Write(filepath.Join(tmp, "testdir"), bytes.NewReader([]byte("test")))
		assert.ErrorIs(t, err, fsys.ErrNotARegularFile)
	})
}

func TestOSFS_Move(t *testing.T) {
	osfs := New()
	t.Run("should move a file", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "testfile")
		err := os.WriteFile(filename, []byte("test"), 0644)
		require.NoError(t, err)

		newFilename := filepath.Join(tmp, "newfile")
		err = osfs.Move(filename, newFilename)
		require.NoError(t, err)

		data, err := ioutil.ReadFile(newFilename)
		require.NoError(t, err)
		assert.Equal(t, []byte("test"), data)

		_, err = os.Stat(filename)
		assert.ErrorIs(t, err, fsys.ErrNotExist)
	})
	t.Run("should return an error if the file does not exist", func(t *testing.T) {
		tmp := createTmp(t)
		err := osfs.Move(filepath.Join(tmp, "testfile"), filepath.Join(tmp, "newfile"))
		assert.ErrorIs(t, err, fsys.ErrNotExist)
	})
}

func TestOSFS_Copy(t *testing.T) {
	osfs := New()
	t.Run("should copy a file", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "testfile")
		err := os.WriteFile(filename, []byte("test"), 0644)
		require.NoError(t, err)

		newFilename := filepath.Join(tmp, "newfile")
		err = osfs.Copy(filename, newFilename)
		require.NoError(t, err)

		data, err := ioutil.ReadFile(newFilename)
		require.NoError(t, err)
		assert.Equal(t, []byte("test"), data)
	})
	t.Run("should copy a directory", func(t *testing.T) {
		tmp := createTmp(t)
		dirname := filepath.Join(tmp, "testdir")
		err := osfs.MkDir(dirname)
		require.NoError(t, err)

		file := filepath.Join(dirname, "testfile")
		err = os.WriteFile(file, []byte("test"), 0644)

		copyDirname := filepath.Join(tmp, "testdir_copy")
		err = osfs.Copy(dirname, copyDirname)
		require.NoError(t, err)

		newDirStat, err := os.Stat(copyDirname)
		require.NoError(t, err)
		assert.True(t, newDirStat.IsDir())

		newFileContent, err := os.ReadFile(filepath.Join(copyDirname, "testfile"))
		require.NoError(t, err)
		assert.Equal(t, []byte("test"), newFileContent)
	})
	t.Run("should return an error if the dst file exists", func(t *testing.T) {
		tmp := createTmp(t)
		filename := filepath.Join(tmp, "testfile")
		err := os.WriteFile(filename, []byte("test"), 0644)
		require.NoError(t, err)

		copyFilename := filepath.Join(tmp, "testfile_copy")
		err = os.WriteFile(copyFilename, []byte("test copy"), 0644)
		require.NoError(t, err)

		err = osfs.Copy(filename, copyFilename)
		require.ErrorIs(t, err, fsys.ErrExist)
	})
}
