package uploadtest

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrAccessDenied  = errors.New("access denied")
	ErrStorageFailed = errors.New("failed to upload file")
)

