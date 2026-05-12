package storage

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ================= MOCKS =================
type mockSalt struct {
	mock.Mock
}

func (m *mockSalt) Generate() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

type errReader struct{}

func (e *errReader) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}

// ================= COMMON =================
func setupFS() *MockFS {
	fs := new(MockFS)
	fs.On("MkdirAll", mock.Anything, mock.Anything).Return(nil)
	return fs
}

type fakeFileInfo struct{}

func (f *fakeFileInfo) Name() string       { return "file" }
func (f *fakeFileInfo) Size() int64        { return 0 }
func (f *fakeFileInfo) Mode() os.FileMode  { return 0 }
func (f *fakeFileInfo) ModTime() time.Time { return time.Now() }
func (f *fakeFileInfo) IsDir() bool        { return false }
func (f *fakeFileInfo) Sys() any           { return nil }

// ================= CONSTRUCTOR =================
func TestNewStorageAdapterWithDeps(t *testing.T) {
	fs := setupFS()
	salt := new(mockSalt)

	adapter := NewStorageAdapterWithDeps(fs, salt)

	assert.NotNil(t, adapter)
	assert.Equal(t, "data", adapter.basePath)
	assert.Equal(t, "data/images", adapter.imagePath)
}

// ================= UploadFile =================
func TestUploadFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)
		salt.On("Generate").Return("abc123", nil)

		tmpFile, _ := os.CreateTemp("", "test")
		defer os.Remove(tmpFile.Name())

		fileName := "file.txt"
		expectedName := "file_abc123.txt"
		fullPath := filepath.Join("data", expectedName)

		fs.On("Stat", fullPath).Return(nil, os.ErrNotExist)
		fs.On("Create", fullPath).Return(tmpFile, nil)

		adapter := NewStorageAdapterWithDeps(fs, salt)

		name, err := adapter.UploadFile(bytes.NewBufferString("data"), fileName)
		assert.NoError(t, err)
		assert.Equal(t, expectedName, name)
	})

	t.Run("invalid name", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)
		adapter := NewStorageAdapterWithDeps(fs, salt)

		name, err := adapter.UploadFile(bytes.NewBufferString("data"), "../bad")
		assert.ErrorIs(t, err, ErrInvalidName)
		assert.Empty(t, name)
	})

	t.Run("salt error", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)
		salt.On("Generate").Return("", errors.New("salt error"))

		adapter := NewStorageAdapterWithDeps(fs, salt)
		name, err := adapter.UploadFile(bytes.NewBufferString("data"), "file.txt")
		assert.Error(t, err)
		assert.Empty(t, name)
	})

	t.Run("create error", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)
		salt.On("Generate").Return("abc123", nil)

		fullPath := filepath.Join("data", "file_abc123.txt")
		fs.On("Stat", fullPath).Return(nil, os.ErrNotExist)
		fs.On("Create", fullPath).Return(nil, errors.New("create error"))

		adapter := NewStorageAdapterWithDeps(fs, salt)
		name, err := adapter.UploadFile(bytes.NewBufferString("data"), "file.txt")
		assert.Error(t, err)
		assert.Empty(t, name)
	})

	t.Run("copy error", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)
		salt.On("Generate").Return("abc123", nil)

		tmpFile, _ := os.CreateTemp("", "test")
		defer os.Remove(tmpFile.Name())

		fullPath := filepath.Join("data", "file_abc123.txt")
		fs.On("Stat", fullPath).Return(nil, os.ErrNotExist)
		fs.On("Create", fullPath).Return(tmpFile, nil)

		adapter := NewStorageAdapterWithDeps(fs, salt)
		name, err := adapter.UploadFile(&errReader{}, "file.txt")
		assert.Error(t, err)
		assert.Empty(t, name)
	})
}

// ================= GetFile =================
func TestGetFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fs := setupFS()
		fs.On("Stat", mock.Anything).Return(&fakeFileInfo{}, nil)

		tmpFile, _ := os.CreateTemp("", "test")
		defer os.Remove(tmpFile.Name())

		fs.On("Open", mock.Anything).Return(tmpFile, nil)

		adapter := NewStorageAdapter(fs)
		file, err := adapter.GetFile("file.txt")
		assert.NoError(t, err)
		assert.NotNil(t, file)
	})

	t.Run("not found", func(t *testing.T) {
		fs := setupFS()
		fs.On("Stat", mock.Anything).Return(nil, os.ErrNotExist)

		adapter := NewStorageAdapter(fs)
		file, err := adapter.GetFile("file.txt")
		assert.NoError(t, err)
		assert.Nil(t, file)
	})

	t.Run("invalid name", func(t *testing.T) {
		fs := setupFS()
		adapter := NewStorageAdapter(fs)
		file, err := adapter.GetFile("bad/name")
		assert.ErrorIs(t, err, ErrInvalidName)
		assert.Nil(t, file)
	})

	t.Run("open error", func(t *testing.T) {
		fs := setupFS()
		fs.On("Stat", mock.Anything).Return(&fakeFileInfo{}, nil)
		fs.On("Open", mock.Anything).Return(nil, errors.New("err"))

		adapter := NewStorageAdapter(fs)
		file, err := adapter.GetFile("file.txt")
		assert.Error(t, err)
		assert.Nil(t, file)
	})
}

// ================= DeleteFile =================
func TestDeleteFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fs := setupFS()
		fs.On("Stat", mock.Anything).Return(&fakeFileInfo{}, nil)
		fs.On("Remove", mock.Anything).Return(nil)

		adapter := NewStorageAdapter(fs)
		err := adapter.DeleteFile("file.txt")
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		fs := setupFS()
		fs.On("Stat", mock.Anything).Return(nil, os.ErrNotExist)

		adapter := NewStorageAdapter(fs)
		err := adapter.DeleteFile("file.txt")
		assert.NoError(t, err)
	})

	t.Run("invalid name", func(t *testing.T) {
		fs := setupFS()
		adapter := NewStorageAdapter(fs)
		err := adapter.DeleteFile("../bad")
		assert.ErrorIs(t, err, ErrInvalidName)
	})

	t.Run("remove error", func(t *testing.T) {
		fs := setupFS()
		fs.On("Stat", mock.Anything).Return(&fakeFileInfo{}, nil)
		fs.On("Remove", mock.Anything).Return(errors.New("err"))

		adapter := NewStorageAdapter(fs)
		err := adapter.DeleteFile("file.txt")
		assert.Error(t, err)
	})
}

// ================= UploadPicture =================
func TestUploadPicture(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)
		salt.On("Generate").Return("abc123", nil)

		tmpFile, _ := os.CreateTemp("", "img")
		defer os.Remove(tmpFile.Name())

		expectedName := "pic_abc123.png"
		fullPath := filepath.Join("data/images", "user1", expectedName)
		fs.On("MkdirAll", filepath.Join("data/images", "user1"), os.ModePerm).Return(nil)
		fs.On("Stat", fullPath).Return(nil, os.ErrNotExist)
		fs.On("Create", fullPath).Return(tmpFile, nil)

		adapter := NewStorageAdapterWithDeps(fs, salt)
		name, err := adapter.UploadPicture(bytes.NewBufferString("img"), "pic.png", "user1")

		assert.NoError(t, err)
		assert.Equal(t, expectedName, name)
	})

	t.Run("invalid filename", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)
		adapter := NewStorageAdapterWithDeps(fs, salt)
		name, err := adapter.UploadPicture(bytes.NewBufferString("img"), "../bad", "user1")
		assert.ErrorIs(t, err, ErrInvalidName)
		assert.Empty(t, name)
	})

	t.Run("invalid login", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)
		adapter := NewStorageAdapterWithDeps(fs, salt)
		name, err := adapter.UploadPicture(bytes.NewBufferString("img"), "pic.png", "../bad")
		assert.ErrorIs(t, err, ErrInvalidName)
		assert.Empty(t, name)
	})

	t.Run("mkdir error", func(t *testing.T) {
		fs := new(MockFS)
		salt := new(mockSalt)
		fs.On("MkdirAll", mock.Anything, mock.Anything).Return(errors.New("mkdir err"))

		adapter := NewStorageAdapterWithDeps(fs, salt)
		name, err := adapter.UploadPicture(bytes.NewBufferString("img"), "pic.png", "user1")
		assert.Error(t, err)
		assert.Empty(t, name)
	})
}

// ================= DeletePictures =================
func TestDeletePictures(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fs := setupFS()
		fs.On("Stat", mock.Anything).Return(&fakeFileInfo{}, nil)
		fs.On("RemoveAll", mock.Anything).Return(nil)

		adapter := NewStorageAdapter(fs)
		err := adapter.DeletePictures("user1")
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		fs := setupFS()
		fs.On("Stat", mock.Anything).Return(nil, os.ErrNotExist)

		adapter := NewStorageAdapter(fs)
		err := adapter.DeletePictures("user1")
		assert.NoError(t, err)
	})

	t.Run("invalid login", func(t *testing.T) {
		fs := setupFS()
		adapter := NewStorageAdapter(fs)
		err := adapter.DeletePictures("../bad")
		assert.ErrorIs(t, err, ErrInvalidName)
	})
}
