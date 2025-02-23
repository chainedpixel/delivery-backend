package dto

import (
	"domain/delivery/models/auth"
	errPackage "domain/error"
	"encoding/json"
	"fmt"
	"io"
)

type LoginRequest struct {
	// User email
	// @example user@example.com
	Email string `json:"email" validate:"required,email"`
	// User password
	// @example mySecurePassword123
	Password string `json:"password" validate:"required"`
	// Device information
	DeviceInfo map[string]interface{} `json:"device_info,omitempty"`
}

// LoginResponse representa la estructura de la respuesta de login
type LoginResponse struct {
	// JWT Token
	// @example eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
	Token string `json:"token"`
}

func NewLoginRequest(body io.ReadCloser) (*LoginRequest, error) {
	var request LoginRequest
	err := json.NewDecoder(body).Decode(&request)
	if err != nil {
		return nil, errPackage.ErrFailedToUnparseJSON
	}

	if err := request.Validate(); err != nil {
		return nil, err
	}

	return &request, nil
}

func (r *LoginRequest) Validate() error {
	if r.Email == "" {
		return fmt.Errorf("email is required")
	}
	if r.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}

func (r *LoginRequest) ParseToCredentialsModel(ipAddress string) *auth.Credentials {
	return &auth.Credentials{
		Email:      r.Email,
		Password:   r.Password,
		DeviceInfo: r.DeviceInfo,
		IPAddress:  ipAddress,
	}
}
