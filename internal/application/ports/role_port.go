package ports

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type RolerUseCase interface {
	GetRoles(ctx context.Context) ([]entities.Role, error)
	GetRoleByIDOrName(ctx context.Context, id string) (*entities.Role, error)
}
