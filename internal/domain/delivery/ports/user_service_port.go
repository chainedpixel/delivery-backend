package ports

import (
	"context"
	"domain/delivery/models/users"
)

type UserService interface {
	GetUserInfo(ctx context.Context, userID string) (*users.User, error)
}
