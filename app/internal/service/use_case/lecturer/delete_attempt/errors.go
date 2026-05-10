package deleteattempt

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrAccessDenied = errors.New("access denied")
	ErrDeleteFailed = errors.New("failed to delete attempt")
)
