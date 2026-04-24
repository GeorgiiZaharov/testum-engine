package deletelecturer

import "errors"

var (
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("user not found")
	ErrNotLecturer  = errors.New("user is not lecturer")
)
