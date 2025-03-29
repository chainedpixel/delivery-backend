package role

import (
	"context"
	"github.com/MarlonG1/delivery-backend/internal/application/ports"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
)

type RolerUseCase struct {
	roleService interfaces.Roler
}

func NewRolerUseCase(roleRepo interfaces.Roler) ports.RolerUseCase {
	return &RolerUseCase{
		roleService: roleRepo,
	}
}

func (r RolerUseCase) GetRoles(ctx context.Context) ([]entities.Role, error) {
	// 1. Obtener todos los roles
	roles, err := r.roleService.GetRoles(ctx)
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (r RolerUseCase) GetRoleByIDOrName(ctx context.Context, param string) (*entities.Role, error) {
	// 1. Obtener el rol por su ID o nombre
	role, err := r.roleService.GetRoleByIDOrName(ctx, param)
	if err != nil {
		return nil, err
	}

	return role, nil
}
