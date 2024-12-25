package types

import "errors"

var (
	ErrInvalidEmail  = errors.New("invalid email format")
	ErrWeakPassword  = errors.New("password does not meet security requirements")
	ErrUserNotFound  = errors.New("user not found")
	ErrUserExists    = errors.New("user already exists")
	ErrInternalError = errors.New("internal server error")
)
