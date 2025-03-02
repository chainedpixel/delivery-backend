package ports

import (
	"context"
	"domain/delivery/models/auth"
	"domain/delivery/models/users"
)

type AuthService interface {
	ValidateCredentials(ctx context.Context, email, password string) (*users.User, error)
	CreateSession(ctx context.Context, user *users.User, deviceInfo map[string]interface{}, ipAddress string) (string, error)
	InvalidateSession(ctx context.Context, token string) error
}

type AuthUseCase interface {
	Authenticate(ctx context.Context, credentials *auth.Credentials) (string, error)
	SignOut(ctx context.Context, token string) error
}
