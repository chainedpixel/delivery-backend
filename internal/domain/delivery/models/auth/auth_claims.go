package auth

import "time"

// AuthClaims representa la información que se almacenará en el token
type AuthClaims struct {
	UserID    string    `json:"user_id"`
	Role      string    `json:"auth"`
	IssuedAt  time.Time `json:"issued_at"`
	ExpiresAt time.Time `json:"expires_at"`
}
