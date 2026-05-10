package result

import "errors"

var (
	ErrResultNotFound = errors.New("result not found")
	ErrDeleteFailed   = errors.New("failed to delete attempt")
)
