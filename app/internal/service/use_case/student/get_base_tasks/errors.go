package getbasetasks

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrAccessDenied  = errors.New("access denied")
	ErrTestCompleted = errors.New("test already completed")
)
