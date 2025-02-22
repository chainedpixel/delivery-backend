package error

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrSessionExpired     = errors.New("session has expired")
	ErrUserInactive       = errors.New("user is inactive")
	ErrUnauthorized       = errors.New("unauthorized access")
)
