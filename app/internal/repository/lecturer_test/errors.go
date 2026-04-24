package lecturertest

import "errors"

var (
	ErrTestNotFound = errors.New("test not found")
	ErrCreateFailed = errors.New("failed to create test")
	ErrDeleteFailed = errors.New("failed to delete test")
)
