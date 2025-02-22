package error

import "errors"

var (
	ErrFailedToParseJSON   = errors.New("failed to marshal content")
	ErrFailedToUnparseJSON = errors.New("failed to unmarshal content")

	ErrFailedToSignToken       = errors.New("failed to sign token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenExpired            = errors.New("token has expired")
	ErrInvalidToken            = errors.New("token is invalid")
)
