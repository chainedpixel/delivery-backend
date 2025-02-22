package ports

import (
	"context"
	"domain/delivery/models/auth"
	"domain/delivery/models/user"
)

type AuthService interface {
	ValidateCredentials(ctx context.Context, email, password string) (*user.User, error)
	CreateSession(ctx context.Context, user *user.User, deviceInfo map[string]interface{}, ipAddress string) (string, error)
	InvalidateSession(ctx context.Context, token string) error
}

type AuthUseCase interface {
	Authenticate(ctx context.Context, credentials *auth.Credentials) (string, error)
	SignOut(ctx context.Context, token string) error
}
