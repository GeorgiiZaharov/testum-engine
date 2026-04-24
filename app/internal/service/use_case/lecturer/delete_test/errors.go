package deletetest

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrAccessDenied  = errors.New("access denied")
	ErrTestNotFound  = errors.New("test not found")
	ErrStorageFailed = errors.New("failed to delete file from storage")
)

