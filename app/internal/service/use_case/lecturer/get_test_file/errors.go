package gettestfile

import "errors"

var (
	ErrForbidden    = errors.New("no access to test")
	ErrFileNotFound = errors.New("test file not found")
)
