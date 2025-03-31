package interfaces

import (
	"context"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type Userer interface {
	GetUserInfo(ctx context.Context, userID string) (*entities.User, error)
	GetUserByID(ctx context.Context, userID string) (*entities.User, error)
	GetUserRoles(ctx context.Context, userID string) ([]entities.Role, error)
	GetAllUsers(ctx context.Context, companyID string, queryParams *entities.UserQueryParams) ([]entities.User, int64, error)
	ActivateOrDeactivateUser(ctx context.Context, userID, loggedUser string, active bool) error
	AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy string) error
	UnassignRole(ctx context.Context, userID, roleID string) error
	CleanAllSessions(ctx context.Context, userID string) error
	UpdateRolesToUser(ctx context.Context, userID string, loggedUserID string, roles []entities.Role) error
	RecoverUser(ctx context.Context, id string) error
	CreateUser(ctx context.Context, user *entities.User) error
	UpdateUser(ctx context.Context, userID string, user *entities.User) error
	DeleteUser(ctx context.Context, userID string) error
}
