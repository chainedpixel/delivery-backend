package interfaces

import (
	"context"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type Roler interface {
	GetRoles(ctx context.Context) ([]entities.Role, error)
	GetRoleByIDOrName(ctx context.Context, param string) (*entities.Role, error)
	IsRoleActive(ctx context.Context, id string) (bool, error)
	IsRoleExist(ctx context.Context, param string) (bool, error)
}
