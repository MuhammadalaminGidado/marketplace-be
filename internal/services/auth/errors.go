package auth

import "errors"

var (
	ErrEmailAlreadyExists = errors.New("email already in use")

	ErrInvalidCredentials = errors.New("invalid email or password")

	ErrSessionExpired = errors.New("session expired")

	ErrInvalidOrExpiredCode = errors.New("invalid or expired code")

	ErrOtpThrottled = errors.New("too many requests")
)
