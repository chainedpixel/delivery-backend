package ports

import (
	"context"
	"domain/delivery/models/user"
)

type UserUseCase interface {
	GetProfileInfo(ctx context.Context, userID string) (*user.User, error)
}
