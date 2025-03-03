package error

import "errors"

var (
	ErrFailedToParseJSON   = errors.New("failed to marshal content")
	ErrFailedToUnparseJSON = errors.New("failed to unmarshal content")

	ErrFailedToSignToken       = errors.New("failed to sign token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenExpired            = errors.New("token has expired")
	ErrInvalidToken            = errors.New("token is invalid")

	ErrInvalidEmail          = errors.New("invalid email format")
	ErrInvalidPassword       = errors.New("invalid password format, minimum 8 characters, at least one uppercase letter, one lowercase letter, one number and one special character")
	ErrValidationErrorsFound = errors.New("validation errors has been found")
	ErrIPAddressNotFound     = errors.New("ip address not found and is required")
)
