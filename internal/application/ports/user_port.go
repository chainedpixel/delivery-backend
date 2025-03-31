package ports

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"net/http"
)

type UserUseCase interface {
	GetProfileInfo(ctx context.Context) (*entities.User, error)
	GetAllUsers(ctx context.Context, request *http.Request) ([]entities.User, *entities.UserQueryParams, int64, error)
	GetUserByID(ctx context.Context, userID string) (*entities.User, error)
	GetUserRoles(ctx context.Context, userID string) ([]entities.Role, error)
	ActivateOrDeactivateUser(ctx context.Context, userID string, active bool) error
	AssignRoleToUser(ctx context.Context, userID, param string) error
	UnassignRole(ctx context.Context, userID, param string) error
	CleanAllSessions(ctx context.Context, userID string) error
	RecoverUser(ctx context.Context, id string) error
	CreateUser(ctx context.Context, user *entities.User) error
	UpdateUser(ctx context.Context, userID string, user *entities.User) error
	DeleteUser(ctx context.Context, userID string) error
}
