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

	ErrInvalidCredentials     = errors.New("invalid email or password")
	ErrInactiveUser           = errors.New("users is inactive")
	ErrInvalidUser            = errors.New("email, firstName, lastName, phone and password are required, please fill them")
	ErrInvalidProfileUser     = errors.New("document type, document number, birth date, emergency contact and phone in profile section are required, please fill them")
	ErrMissingProfileSection  = errors.New("profile section is required, please fill it")
	ErrRoleMissing            = errors.New("role_id is required, provide them")
	ErrInvalidRole            = errors.New("role is invalid, please provide a valid role")
	ErrReasonToDeactivateUser = errors.New("when you want deactivate user reason field must be provide")
	ErrMissingRoles           = errors.New("at least one role is required, please provide them")

	ErrNilOrder = errors.New("order cannot be nil, please provide a valid order")
	ErrNilQR    = errors.New("qr code cannot be nil")

	ErrSessionNotFound   = errors.New("the session assigned to the token was not found, probably was deleted or expired")
	ErrSessionDBNotFound = errors.New("the session assigned to the token was not found")
	ErrGenericDBError    = errors.New("an error occurred while trying to execute the operation in the database")

	ErrAuthorizationHeaderNotFound = errors.New("authorization header not found, please provide a valid token")
	ErrInvalidAuthorizationFormat  = errors.New("invalid authorization format, the format should be 'Bearer <token>'")
	ErrTokenExpiredOrTampered      = errors.New("token is expired or has been tampered with, please provide a valid token")

	ErrInvalidEmailConfig   = errors.New("invalid email configuration")
	ErrSMTPConnectionFailed = errors.New("failed to connect to SMTP server")

	ErrNoRecipients        = errors.New("no recipients specified")
	ErrEmailSendFailed     = errors.New("failed to send email")
	ErrInvalidEmailAddress = errors.New("invalid email address")

	ErrTemplateNotFound         = errors.New("email template not found")
	ErrInvalidTemplate          = errors.New("invalid email template")
	ErrTemplateProcessingFailed = errors.New("failed to process email template")

	ErrEmptyEmailBody    = errors.New("email body cannot be empty")
	ErrEmptyEmailSubject = errors.New("email subject cannot be empty")

	ErrInvalidAttachment  = errors.New("invalid email attachment")
	ErrAttachmentTooLarge = errors.New("attachment exceeds size limit")
)
