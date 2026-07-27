package auth

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already in use")

	ErrInvalidCredentials = errors.New("invalid email or password")

	ErrSessionExpired = errors.New("session expired")
)
