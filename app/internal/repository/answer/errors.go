package answer

import "errors"

var (
	ErrSaveFailed    = errors.New("failed to save answers")
	ErrDeleteAttempt = errors.New("failed to delete attempt")
	ErrGetAnswers    = errors.New("failed to get answers")
	ErrNotFound      = errors.New("answers not found")
)

