package getme

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrNotFound     = errors.New("user not found")
	ErrLDAPFailed   = errors.New("failed to fetch user from ldap")
)
