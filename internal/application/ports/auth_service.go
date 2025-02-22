package ports

import (
	"context"
	"domain/delivery/models/auth"
	"domain/delivery/models/user"
)

type AuthService interface {
	Login(ctx context.Context, credentials *auth.Credentials) (*user.User, string, error)
	Logout(ctx context.Context, token string) error
}
