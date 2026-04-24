package refresh

import "errors"

var (
	ErrInvalidUserID = errors.New("invalid user id")
	ErrAuthFailed    = errors.New("failed to generate auth tokens")
)

