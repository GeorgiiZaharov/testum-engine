package gettestresult

import "errors"

var (
	ErrAccessDenied   = errors.New("access denied")
	ErrResultNotFound = errors.New("result not found")
	ErrInvalidInput   = errors.New("invalid input")
)

