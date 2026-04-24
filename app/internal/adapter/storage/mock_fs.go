package storage

import (
	"os"

	"github.com/stretchr/testify/mock"
)

// MockFS реализует FileSystem интерфейс для тестов
type MockFS struct {
	mock.Mock
}

func (m *MockFS) Stat(name string) (os.FileInfo, error) {
	args := m.Called(name)

	if fi := args.Get(0); fi != nil {
		return fi.(os.FileInfo), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFS) Create(name string) (*os.File, error) {
	args := m.Called(name)

	if f := args.Get(0); f != nil {
		return f.(*os.File), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFS) Open(name string) (*os.File, error) {
	args := m.Called(name)

	if f := args.Get(0); f != nil {
		return f.(*os.File), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockFS) Remove(name string) error {
	args := m.Called(name)
	return args.Error(0)
}

func (m *MockFS) RemoveAll(path string) error {
	args := m.Called(path)
	return args.Error(0)
}

func (m *MockFS) MkdirAll(path string, perm os.FileMode) error {
	args := m.Called(path, perm)
	return args.Error(0)
}
