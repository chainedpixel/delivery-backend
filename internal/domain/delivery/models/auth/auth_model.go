package auth

import (
	"domain/delivery/value_objects"
	errPackage "domain/error"
	"time"
)

// Credentials representa los datos necesarios para la autenticación
type Credentials struct {
	Email      string
	Password   string
	DeviceInfo map[string]interface{}
	IPAddress  string
	CreatedAt  time.Time
}

func NewCredentials(email, password string, deviceInfo map[string]interface{}, ipAddress string) (*Credentials, error) {
	domainError := errPackage.NewDomainError("Auth", "NewCredentials", "Failed to create new credentials")

	if !value_objects.NewEmail(email).IsValid() {
		domainError.AddValidationError(errPackage.ErrInvalidEmail)
	}

	if !value_objects.NewPassword(password).IsValid() {
		domainError.AddValidationError(errPackage.ErrInvalidPassword)
	}

	if ipAddress == "" {
		domainError.AddValidationError(errPackage.ErrIPAddressNotFound)
	}

	if domainError.HasValidationErrors() {
		domainError.AsideError(errPackage.ErrValidationErrorsFound)
		return nil, domainError
	}

	return &Credentials{
		Email:      email,
		Password:   password,
		DeviceInfo: deviceInfo,
		IPAddress:  ipAddress,
	}, nil

}
