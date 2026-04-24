package gethardtasks

import "errors"

var (
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidInput = errors.New("invalid input")
	ErrAccessDenied = errors.New("access denied")
	ErrGetHardTasks = errors.New("failed to get hard tasks")
)

