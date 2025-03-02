package ports

import (
	"context"
	"domain/delivery/models/users"
)

type UserUseCase interface {
	GetProfileInfo(ctx context.Context, userID string) (*users.User, error)
}
