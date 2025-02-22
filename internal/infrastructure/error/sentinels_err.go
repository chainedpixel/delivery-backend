package error

import "errors"

var (
	ErrFailedToConnectRedis = errors.New("failed to connect to redis")
	ErrFailedToPingRedis    = errors.New("failed to ping redis")
	ErrFailedToSetKeyRedis  = errors.New("failed to set key in redis")
	ErrTokenNotFound        = errors.New("token not found")
	ErrFailedToGetKey       = errors.New("failed to get key from redis")
	ErrFailedToDeleteKey    = errors.New("failed to delete key from redis")
	ErrFailedToCloseRedis   = errors.New("failed to close redis connection")
	ErrFailedRPush          = errors.New("failed to execute RPush command in redis")
	ErrFailedLPush          = errors.New("failed to execute LPush command in redis")
	ErrFailedLRange         = errors.New("failed to execute LRange command in redis")
	ErrFailedLLen           = errors.New("failed to execute LLen command in redis")
	ErrFailedLTrim          = errors.New("failed to execute LTrim command in redis")

	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInactiveUser       = errors.New("user is inactive")

	ErrAuthorizationHeaderNotFound = errors.New("authorization header not found, please provide a valid token")
	ErrInvalidAuthorizationFormat  = errors.New("invalid authorization format, the format should be 'Bearer <token>'")
	ErrTokenExpiredOrTampered      = errors.New("token is expired or has been tampered with, please provide a valid token")
)
