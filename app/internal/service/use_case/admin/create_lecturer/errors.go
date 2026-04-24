package createlecturer

import "errors"

var (
	ErrForbidden    = errors.New("forbidden")
	ErrInvalidInput = errors.New("invalid input")
	ErrLDAPFailed   = errors.New("failed to fetch user from ldap")
)
