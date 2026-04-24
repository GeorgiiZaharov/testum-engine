package ldap

import "errors"

var (
	// Общие ошибки
	ErrConnectionFailed = errors.New("ldap: connection failed")
	ErrSearchFailed     = errors.New("ldap: search failed")

	// Бизнес-ошибки
	ErrUserNotFound     = errors.New("ldap: user not found")
	ErrInvalidPassword  = errors.New("ldap: invalid credentials")
	ErrEmptyCredentials = errors.New("ldap: empty credentials")
	ErrEmptyLogin       = errors.New("ldap: empty login")
)
