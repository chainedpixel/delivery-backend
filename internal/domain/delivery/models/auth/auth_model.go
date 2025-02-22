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
