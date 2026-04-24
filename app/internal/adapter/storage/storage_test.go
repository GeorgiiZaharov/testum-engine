package storage

import (
	"bytes"
	"errors"
	"os"
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
func (f *fakeFileInfo) Sys() interface{}   { return nil }

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

		fs.On("Stat", mock.Anything).Return(nil, os.ErrNotExist)

		tmpFile, _ := os.CreateTemp("", "test")

		defer func() { _ = os.Remove(tmpFile.Name()) }()

		fs.On("Create", mock.Anything).Return(tmpFile, nil)

		adapter := NewStorageAdapter(fs)

		err := adapter.UploadFile(bytes.NewBufferString("data"), "file.txt")
		assert.NoError(t, err)
	})

	t.Run("file exists", func(t *testing.T) {
		fs := setupFS()
		fs.On("Stat", mock.Anything).Return(&fakeFileInfo{}, nil)

		adapter := NewStorageAdapter(fs)

		err := adapter.UploadFile(bytes.NewBufferString("data"), "file.txt")
		assert.ErrorIs(t, err, ErrFileExists)
	})

	t.Run("invalid name", func(t *testing.T) {
		fs := setupFS()
		adapter := NewStorageAdapter(fs)

		err := adapter.UploadFile(bytes.NewBufferString("data"), "../bad")
		assert.ErrorIs(t, err, ErrInvalidName)
	})

	t.Run("create error", func(t *testing.T) {
		fs := setupFS()

		fs.On("Stat", mock.Anything).Return(nil, os.ErrNotExist)
		fs.On("Create", mock.Anything).Return(nil, errors.New("err"))

		adapter := NewStorageAdapter(fs)

		err := adapter.UploadFile(bytes.NewBufferString("data"), "file.txt")
		assert.Error(t, err)
	})

	t.Run("copy error", func(t *testing.T) {
		fs := setupFS()

		fs.On("Stat", mock.Anything).Return(nil, os.ErrNotExist)

		tmpFile, _ := os.CreateTemp("", "test")
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		fs.On("Create", mock.Anything).Return(tmpFile, nil)

		adapter := NewStorageAdapter(fs)

		err := adapter.UploadFile(&errReader{}, "file.txt")
		assert.Error(t, err)
	})
}

// ================= GetFile =================
func TestGetFile(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fs := setupFS()

		fs.On("Stat", mock.Anything).Return(&fakeFileInfo{}, nil)

		tmpFile, _ := os.CreateTemp("", "test")
		defer func() { _ = os.Remove(tmpFile.Name()) }()

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
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		fs.On("Create", mock.Anything).Return(tmpFile, nil)

		adapter := NewStorageAdapterWithDeps(fs, salt)

		name, err := adapter.UploadPicture(
			bytes.NewBufferString("img"),
			"pic.png",
			"user1",
		)

		assert.NoError(t, err)
		assert.Equal(t, "pic_abc123.png", name)
	})

	t.Run("invalid filename", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)

		adapter := NewStorageAdapterWithDeps(fs, salt)

		name, err := adapter.UploadPicture(
			bytes.NewBufferString("img"),
			"../bad",
			"user1",
		)

		assert.ErrorIs(t, err, ErrInvalidName)
		assert.Empty(t, name)
	})

	t.Run("invalid login", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)

		adapter := NewStorageAdapterWithDeps(fs, salt)

		name, err := adapter.UploadPicture(
			bytes.NewBufferString("img"),
			"pic.png",
			"../bad",
		)

		assert.ErrorIs(t, err, ErrInvalidName)
		assert.Empty(t, name)
	})

	t.Run("mkdir error", func(t *testing.T) {
		fs := new(MockFS)
		salt := new(mockSalt)

		fs.On("MkdirAll", mock.Anything, mock.Anything).
			Return(errors.New("mkdir err"))

		adapter := NewStorageAdapterWithDeps(fs, salt)

		name, err := adapter.UploadPicture(
			bytes.NewBufferString("img"),
			"pic.png",
			"user1",
		)

		assert.Error(t, err)
		assert.Empty(t, name)
	})

	t.Run("salt error", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)

		salt.On("Generate").Return("", errors.New("salt error"))

		adapter := NewStorageAdapterWithDeps(fs, salt)

		name, err := adapter.UploadPicture(
			bytes.NewBufferString("img"),
			"pic.png",
			"user1",
		)

		assert.Error(t, err)
		assert.Empty(t, name)
	})

	t.Run("create error", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)

		salt.On("Generate").Return("abc", nil)
		fs.On("Create", mock.Anything).Return(nil, errors.New("err"))

		adapter := NewStorageAdapterWithDeps(fs, salt)

		name, err := adapter.UploadPicture(
			bytes.NewBufferString("img"),
			"pic.png",
			"user1",
		)

		assert.Error(t, err)
		assert.Empty(t, name)
	})

	t.Run("copy error", func(t *testing.T) {
		fs := setupFS()
		salt := new(mockSalt)

		salt.On("Generate").Return("abc", nil)

		tmpFile, _ := os.CreateTemp("", "img")
		defer func() { _ = os.Remove(tmpFile.Name()) }()

		fs.On("Create", mock.Anything).Return(tmpFile, nil)

		adapter := NewStorageAdapterWithDeps(fs, salt)

		name, err := adapter.UploadPicture(
			&errReader{},
			"pic.png",
			"user1",
		)

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

	t.Run("stat error (not IsNotExist)", func(t *testing.T) {
		fs := setupFS()

		fs.On("Stat", mock.Anything).Return(nil, errors.New("stat error"))

		adapter := NewStorageAdapter(fs)

		err := adapter.DeletePictures("user1")

		// пойдет дальше и вызовет RemoveAll с тем же путем
		assert.Error(t, err)
	})

	t.Run("removeAll error", func(t *testing.T) {
		fs := setupFS()

		fs.On("Stat", mock.Anything).Return(&fakeFileInfo{}, nil)
		fs.On("RemoveAll", mock.Anything).Return(errors.New("remove error"))

		adapter := NewStorageAdapter(fs)

		err := adapter.DeletePictures("user1")
		assert.Error(t, err)
	})
}
