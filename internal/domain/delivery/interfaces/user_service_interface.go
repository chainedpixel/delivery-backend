package interfaces

import (
	"context"
	"domain/delivery/models/entities"
)

type Userer interface {
	GetUserInfo(ctx context.Context, userID string) (*entities.User, error)
}
