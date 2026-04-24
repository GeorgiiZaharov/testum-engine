package file

import "errors"

var (
	ErrFilesNotFound  = errors.New("files not found")
	ErrImagesNotFound = errors.New("images not found")
)

