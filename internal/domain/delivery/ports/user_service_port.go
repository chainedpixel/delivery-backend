package ports

import (
	"context"
	"domain/delivery/models/user"
)

type UserService interface {
	GetUserInfo(ctx context.Context, userID string) (*user.User, error)
}
