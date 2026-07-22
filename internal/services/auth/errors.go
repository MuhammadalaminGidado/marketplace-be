package auth

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already in use")

	ErrPasswordMismatch = errors.New("password confirmation doesn't match")

	ErrInvalidCredentials = errors.New("invalid email or password")

	ErrSessionExpired = errors.New("session expired")
)
