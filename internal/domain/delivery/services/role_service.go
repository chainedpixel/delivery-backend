package services

import (
	"context"

	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/interfaces"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/ports"
	error2 "github.com/MarlonG1/delivery-backend/internal/domain/error"
)

type RolerService struct {
	roleRepo ports.RolerRepository
}

func NewRoleService(roleRepo ports.RolerRepository) interfaces.Roler {
	return &RolerService{
		roleRepo: roleRepo,
	}
}

func (r RolerService) GetRoles(ctx context.Context) ([]entities.Role, error) {
	roles, err := r.roleRepo.ListRoles(ctx)
	if err != nil {
		return nil, error2.NewDomainErrorWithCause("RoleService", "GetRoles", "failed to get roles", err)
	}

	return roles, nil
}

func (r RolerService) GetRoleByIDOrName(ctx context.Context, param string) (*entities.Role, error) {
	// 1. Verificar si el rol existe
	role, err := r.roleRepo.GetRoleByIDOrName(ctx, param)
	if err != nil {
		return nil, error2.NewDomainErrorWithCause("RoleService", "GetRoleByID", "failed to get role", err)
	}

	// 2. Verificar si el rol está activo
	isActive, err := r.roleRepo.IsRoleActive(ctx, param)
	if err != nil {
		return nil, error2.NewDomainErrorWithCause("RoleService", "GetRoleByID", "failed to check if role is active", err)
	}
	if !isActive {
		return nil, error2.NewDomainError("RoleService", "GetRoleByID", error2.ErrRoleIsNotActive.Error())
	}

	return role, nil
}

func (r RolerService) IsRoleExist(ctx context.Context, param string) (bool, error) {
	// 1. Verificar si el rol existe
	isExist, err := r.roleRepo.IsRoleExist(ctx, param)
	if err != nil {
		return false, error2.NewDomainErrorWithCause("RoleService", "IsRoleExist", "failed to check if role exists", err)
	}

	return isExist, nil
}

func (r RolerService) IsRoleActive(ctx context.Context, id string) (bool, error) {
	// 1. Verificar si el rol está activo
	isActive, err := r.roleRepo.IsRoleActive(ctx, id)
	if err != nil {
		return false, error2.NewDomainErrorWithCause("RoleService", "IsRoleActive", "failed to check if role is active", err)
	}

	return isActive, nil
}
