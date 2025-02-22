package error

import "errors"

var (
	ErrFailedToMarshalClaims   = errors.New("failed to marshal claims")
	ErrFailedToUnmarshalClaims = errors.New("failed to unmarshal claims")

	ErrFailedToSignToken       = errors.New("failed to sign token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenExpired            = errors.New("token has expired")
	ErrInvalidToken            = errors.New("token is invalid")
)
