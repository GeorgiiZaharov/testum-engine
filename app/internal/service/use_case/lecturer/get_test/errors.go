package gettest

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrAccessDenied = errors.New("access denied")
	ErrNotFound     = errors.New("test not found")
)

