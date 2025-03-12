package ports

import (
	"context"
	"domain/delivery/models/entities"
)

type UserUseCase interface {
	GetProfileInfo(ctx context.Context) (*entities.User, error)
	GetUserByID(ctx context.Context, userID string) (*entities.User, error)
	ActivateOrDeactivateUser(ctx context.Context, userID string, active bool) error
	AssignRoleToUser(ctx context.Context, userID, roleID string) error
	RecoverUser(ctx context.Context, id string) error
	CreateUser(ctx context.Context, user *entities.User) error
	UpdateUser(ctx context.Context, userID string, user *entities.User) error
	DeleteUser(ctx context.Context, userID string) error
}
