package errorlib

import "errors"

var (
	ErrWrongPassword   = errors.New("validate: wrong password")
	ErrUserNotFound    = errors.New("record: user not found")
	ErrEmailRegistered = errors.New("record: user with this email already registered")
	ErrInvalidToken    = errors.New("token: invalid session, please re-login to refresh your session")
	ErrTokenExpired    = errors.New("token: session is expire, please re-login to refresh your session")
	ErrAlreadyRequest  = errors.New("conflict: already requested, please wait to be accepted")
	ErrAdminOnly       = errors.New("unauthorized: this resource only available for admin")
)
