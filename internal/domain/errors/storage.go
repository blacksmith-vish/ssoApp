package errs

import "errors"

var (
	ErrUserExists   = errors.New("user exists already")
	ErrUserNotFound = errors.New("user not found")
	ErrAppNotFound  = errors.New("app not found")
)
