package interfaces

import (
	"context"
	"domain/delivery/models/entities"
)

type UserService interface {
	GetUserInfo(ctx context.Context, userID string) (*entities.User, error)
}
