package dto

import (
	"encoding/json"
	"fmt"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/auth"
	errPackage "github.com/MarlonG1/delivery-backend/internal/domain/error"
	"io"
)

type LoginRequest struct {
	// User email
	Email string `json:"email" validate:"required,email"`
	// User password
	Password string `json:"password" validate:"required"`
	// Device information
	DeviceInfo map[string]interface{} `json:"device_info,omitempty"`
}

// LoginResponse representa la estructura de la respuesta de login
type LoginResponse struct {
	// JWT Token
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
