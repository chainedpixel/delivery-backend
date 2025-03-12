package error

import "errors"

var (
	ErrFailedToParseJSON            = errors.New("failed to marshal content")
	ErrFailedToUnparseJSON          = errors.New("failed to unmarshal content")
	ErrCannotDeleteOrder            = errors.New("the order cannot be deleted, only orders with status 'pending', 'cancelled' or 'restored' can be deleted")
	ErrCannotUpdateOrder            = errors.New("the order cannot be updated, only orders with status 'pending', 'completed', 'picked up', 'in transit or 'in warehouse' can be updated")
	ErrOrderAlreadyDeleted          = errors.New("the order has already been deleted")
	ErrOrderNotDeleted              = errors.New("the order has not been deleted")
	ErrOrderDeleted                 = errors.New("the order has been deleted")
	ErrDeliveryDeadlineBeforePickup = errors.New("delivery deadline must be after pickup deadline")

	ErrFailedToSignToken       = errors.New("failed to sign token")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrTokenExpired            = errors.New("token has expired")
	ErrInvalidToken            = errors.New("token is invalid")

	ErrUserDeactivated                      = errors.New("user is deactivated")
	ErrUserCannotActivateOrDeactivateItself = errors.New("user cannot activate or deactivate itself")
	ErrUserDeleted                          = errors.New("user is deleted")
	ErrUserAlreadyActiveOrInactive          = errors.New("user is already active or inactive")
	ErrUserNotDeleted                       = errors.New("user is not deleted")

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
