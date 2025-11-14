package errorlib

import "errors"

var (
	ErrWrongPassword   = errors.New("validate: wrong password")
	ErrUserNotFound    = errors.New("record: user not found")
	ErrEmailRegistered = errors.New("record: user with this email already registered")
)
