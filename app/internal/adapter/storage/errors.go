package storage

import "errors"

var (
	ErrFileExists  = errors.New("file already exists")
	ErrInvalidName = errors.New("invalid file name")
)
