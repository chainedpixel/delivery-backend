package ports

import (
	"context"
	"domain/delivery/models/entities"
)

type UserUseCase interface {
	GetProfileInfo(ctx context.Context, userID string) (*entities.User, error)
}
