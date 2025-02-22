package auth

import "time"

// Credentials representa los datos necesarios para la autenticación
type Credentials struct {
	Email      string
	Password   string
	DeviceInfo map[string]interface{}
	IPAddress  string
	CreatedAt  time.Time
}

// NewCredentials crea una nueva instancia de credenciales
func NewCredentials(email, password, ipAddress string, deviceInfo map[string]interface{}) *Credentials {
	return &Credentials{
		Email:      email,
		Password:   password,
		DeviceInfo: deviceInfo,
		IPAddress:  ipAddress,
		CreatedAt:  time.Now(),
	}
}
