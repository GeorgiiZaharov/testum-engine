package login

import "errors"

var (
	ErrInvalidInput  = errors.New("invalid input")
	ErrUnauthorized  = errors.New("invalid credentials")
	ErrLDAPFailed    = errors.New("ldap failed")
	ErrTokenGenerate = errors.New("failed to generate token")
)

