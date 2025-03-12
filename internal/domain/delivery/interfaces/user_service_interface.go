package interfaces

import (
	"context"
	"domain/delivery/models/entities"
)

type Userer interface {
	GetUserInfo(ctx context.Context, userID string) (*entities.User, error)
	GetUserByID(ctx context.Context, userID string) (*entities.User, error)
	ActivateOrDeactivateUser(ctx context.Context, userID, loggedUser string, active bool) error
	AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy string) error
	RecoverUser(ctx context.Context, id string) error
	CreateUser(ctx context.Context, user *entities.User) error
	UpdateUser(ctx context.Context, userID string, user *entities.User) error
	DeleteUser(ctx context.Context, userID string) error
}
