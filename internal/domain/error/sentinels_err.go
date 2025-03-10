package error

import "errors"

var (
	ErrFailedToParseJSON   = errors.New("failed to marshal content")
	ErrFailedToUnparseJSON = errors.New("failed to unmarshal content")

	ErrFailedToSignToken       = errors.New("failed to sign token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenExpired            = errors.New("token has expired")
	ErrInvalidToken            = errors.New("token is invalid")

	ErrInvalidEmail               = errors.New("invalid email format")
	ErrInvalidPassword            = errors.New("invalid password format, minimum 8 characters, at least one uppercase letter, one lowercase letter, one number and one special character")
	ErrValidationErrorsFound      = errors.New("validation errors has been found")
	ErrIPAddressNotFound          = errors.New("ip address not found and is required")
	ErrIDsNotFound                = errors.New("all IDs must be provided, company, branch and client")
	ErrTrackingNumber             = errors.New("tracking number is required")
	ErrStatusRequired             = errors.New("status is required")
	ErrOrderDetailsRequired       = errors.New("order details are required")
	ErrDeliveryAddress            = errors.New("delivery address is required")
	ErrPickupAddress              = errors.New("pickup address is required")
	ErrPackageDetails             = errors.New("package details are required")
	ErrCompanyIDRequired          = errors.New("company ID is required")
	ErrCompanyPickUpIDRequired    = errors.New("company pick up ID is required")
	ErrBranchIDRequired           = errors.New("branch ID is required")
	ErrClientIDRequired           = errors.New("client ID is required")
	ErrUserNotFoundOrUnauthorized = errors.New("the user is not found or is unauthorized to perform this action")
)
