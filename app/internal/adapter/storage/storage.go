package storage

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ========================
// INTERFACES
// ========================

type SaltGenerator interface {
	Generate() (string, error)
}

// ========================
// DEFAULT IMPLEMENTATION
// ========================

type cryptoSaltGenerator struct{}

func (c *cryptoSaltGenerator) Generate() (string, error) {
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// ========================
// ADAPTER
// ========================

type StorageAdapter struct {
	fs        FileSystem
	saltGen   SaltGenerator
	basePath  string
	imagePath string
}

func NewStorageAdapter(fs FileSystem) *StorageAdapter {
	return NewStorageAdapterWithDeps(fs, &cryptoSaltGenerator{})
}

func NewStorageAdapterWithDeps(fs FileSystem, saltGen SaltGenerator) *StorageAdapter {
	base := "data"
	images := "../public_html/testum-images"

	_ = fs.MkdirAll(base, os.ModePerm)
	_ = fs.MkdirAll(images, os.ModePerm)

	return &StorageAdapter{
		fs:        fs,
		saltGen:   saltGen,
		basePath:  base,
		imagePath: images,
	}
}

// ========================
// HELPERS
// ========================

func (s *StorageAdapter) validateName(name string) error {
	if name == "" || strings.Contains(name, "/") || strings.Contains(name, "\\") {
		return ErrInvalidName
	}
	return nil
}

func (s *StorageAdapter) fullPath(name string) string {
	return filepath.Join(s.basePath, name)
}

func (s *StorageAdapter) fullImagePath(parts ...string) string {
	all := append([]string{s.imagePath}, parts...)
	return filepath.Join(all...)
}

// ========================
// FILES
// ========================

func (s *StorageAdapter) UploadFile(file io.Reader, fileName string) (string, error) {
	if err := s.validateName(fileName); err != nil {
		return "", err
	}

	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)

	salt, err := s.saltGen.Generate()
	if err != nil {
		return "", err
	}

	newName := name + "_" + salt + ext
	fullPath := s.fullPath(newName)

	// Проверяем существует ли файл с таким именем (редко, но может быть)
	if _, err := s.fs.Stat(fullPath); err == nil {
		return "", ErrFileExists
	}

	f, err := s.fs.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	if _, err := io.Copy(f, file); err != nil {
		return "", err
	}

	return newName, nil
}

func (s *StorageAdapter) GetFile(fileName string) (*os.File, error) {
	if err := s.validateName(fileName); err != nil {
		return nil, err
	}

	path := s.fullPath(fileName)

	if _, err := s.fs.Stat(path); os.IsNotExist(err) {
		return nil, nil
	}

	return s.fs.Open(path)
}

func (s *StorageAdapter) DeleteFile(fileName string) error {
	if err := s.validateName(fileName); err != nil {
		return err
	}

	path := s.fullPath(fileName)

	if _, err := s.fs.Stat(path); os.IsNotExist(err) {
		return nil
	}

	return s.fs.Remove(path)
}

// ========================
// PICTURES (NEW LOGIC)
// ========================

// UploadPicture сохраняет файл в data/images/{login}/
// имя: original + "_" + salt + ext
func (s *StorageAdapter) UploadPicture(file io.Reader, fileName string, login string) (string, error) {
	if err := s.validateName(fileName); err != nil {
		return "", err
	}
	if err := s.validateName(login); err != nil {
		return "", err
	}

	// директория пользователя
	userDir := s.fullImagePath(login)

	// создаём папку если нет
	if err := s.fs.MkdirAll(userDir, os.ModePerm); err != nil {
		return "", err
	}

	ext := filepath.Ext(fileName)
	name := strings.TrimSuffix(fileName, ext)

	salt, err := s.saltGen.Generate()
	if err != nil {
		return "", err
	}

	newName := name + "_" + salt + ext
	fullPath := filepath.Join(userDir, newName)

	f, err := s.fs.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	if _, err := io.Copy(f, file); err != nil {
		return "", err
	}

	return login + "/" + newName, nil
}

// DeletePictures удаляет ВСЮ папку пользователя
func (s *StorageAdapter) DeletePictures(login string) error {
	if err := s.validateName(login); err != nil {
		return err
	}

	userDir := s.fullImagePath(login)

	if _, err := s.fs.Stat(userDir); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return s.fs.RemoveAll(userDir)
}
