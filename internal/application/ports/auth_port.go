package ports

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/auth"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type Authenticator interface {
	ValidateCredentials(ctx context.Context, email, password string) (*entities.User, error)
	CreateSession(ctx context.Context, user *entities.User, deviceInfo map[string]interface{}, ipAddress string) (string, error)
	InvalidateSession(ctx context.Context, token string) error
}

type AuthenticatorUseCase interface {
	Authenticate(ctx context.Context, credentials *auth.Credentials) (string, error)
	SignOut(ctx context.Context, token string) error
}
