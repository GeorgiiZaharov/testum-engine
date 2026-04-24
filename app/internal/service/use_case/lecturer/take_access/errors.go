package takeaccess

import "errors"

var (
	ErrAccessDenied = errors.New("access denied")
	ErrInvalidInput = errors.New("invalid input")
)

