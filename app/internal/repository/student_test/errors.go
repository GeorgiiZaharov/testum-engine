package studenttest

import "errors"

var (
	ErrTestNotFound        = errors.New("test not found")
	ErrStudentTestNotFound = errors.New("student test not found")
	ErrAlreadyStarted      = errors.New("test already started for this student")
	ErrFinishFailed        = errors.New("failed to finish test")
)
