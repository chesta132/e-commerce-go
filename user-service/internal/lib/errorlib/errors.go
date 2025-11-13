package errorlib

import "errors"

var (
	ErrWrongPassword = errors.New("validate: wrong password")
	ErrUserNotFound  = errors.New("record: user not found")
)
