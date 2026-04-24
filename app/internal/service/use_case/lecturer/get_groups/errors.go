package getgroups

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrAccessDenied = errors.New("access denied")
)
