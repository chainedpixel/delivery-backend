package services

import (
	"context"
	"domain/delivery/interfaces"
	"domain/delivery/models/entities"
	"domain/delivery/ports"
	error2 "domain/error"
	"errors"
	"gorm.io/gorm"
	"shared/logs"
)

type userService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) interfaces.Userer {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetAllUsers(ctx context.Context, companyID string, queryParams *entities.UserQueryParams) ([]entities.User, int64, error) {
	// 1. Obtener los usuarios
	users, total, err := s.userRepo.GetAllUsersFromCompany(ctx, companyID, queryParams)
	if err != nil {
		logs.Error("Failed to list users", map[string]interface{}{
			"error": err.Error(),
		})
		return nil, 0, error2.NewDomainErrorWithCause("UserService", "GetUsers", "failed to list users", err)
	}

	return users, total, nil
}

func (s *userService) GetUserInfo(ctx context.Context, userID string) (*entities.User, error) {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	user, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	user, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) CreateUser(ctx context.Context, user *entities.User) error {
	// 1. Crear el usuario
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		logs.Error("Failed to create user", map[string]interface{}{
			"error": err.Error(),
		})
		return error2.NewDomainErrorWithCause("UserService", "CreateUser", "failed to create user", err)
	}

	return nil
}

func (s *userService) UpdateUser(ctx context.Context, userID string, user *entities.User) error {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	_, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return err
	}

	// 2. Actualizar el usuario
	err = s.userRepo.Update(ctx, userID, user)
	if err != nil {
		logs.Error("Failed to update user", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return error2.NewDomainErrorWithCause("UserService", "UpdateUser", "failed to update user", err)
	}

	return nil
}

func (s *userService) RecoverUser(ctx context.Context, userID string) error {
	// 1. Verificar si el usuario existe
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logs.Error("Failed to get user by ID", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return error2.NewDomainErrorWithCause("UserService", "RecoverUser", "failed to get user by ID", err)
	}

	// 2. Verificar si el usuario está eliminado
	if user.DeletedAt == nil {
		logs.Error("User is not deleted", map[string]interface{}{
			"user_id": userID,
		})
		return error2.NewDomainError("UserService", "RecoverUser", error2.ErrUserNotDeleted.Error())
	}

	// 3. Recuperar el usuario
	err = s.userRepo.Recover(ctx, userID)
	if err != nil {
		logs.Error("Failed to recover user", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return error2.NewDomainErrorWithCause("UserService", "RecoverUser", "failed to recover user", err)
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, userID string) error {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	_, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return err
	}

	// 2. Eliminar el usuario
	err = s.userRepo.Delete(ctx, userID)
	if err != nil {
		logs.Error("Failed to delete user", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return error2.NewDomainErrorWithCause("UserService", "DeleteUser", "failed to delete user", err)
	}

	return nil
}

func (s *userService) validateUserFromRepository(ctx context.Context, userID string) (*entities.User, error) {
	// 1. Verificar si el usuario existe
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logs.Error("Failed to get user by ID", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return nil, error2.NewDomainErrorWithCause("UserService", "GetUserByID", "failed to get user by ID", err)
	}

	// 2. Verificar si el usuario está activo
	isActive, err := s.userRepo.IsUserActive(ctx, userID)
	if err != nil {
		logs.Error("Failed to get user active status", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return nil, error2.NewDomainErrorWithCause("UserService", "GetUserByID", "failed to get user active status", err)
	}
	if !isActive {
		logs.Error("User is not active", map[string]interface{}{
			"user_id": userID,
		})
		return nil, error2.NewDomainError("UserService", "GetUserByID", error2.ErrUserDeactivated.Error())
	}

	// 3. Verificar si el usuario está eliminado
	isDeleted, err := s.userRepo.IsUserDeleted(ctx, userID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logs.Error("Failed to get user deleted status", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return nil, error2.NewDomainErrorWithCause("UserService", "GetUserByID", "failed to get user deleted status", err)
	}
	if isDeleted {
		logs.Error("User is deleted", map[string]interface{}{
			"user_id": userID,
		})
		return nil, error2.NewDomainError("UserService", "GetUserByID", error2.ErrUserDeleted.Error())
	}

	return user, nil
}

func (s *userService) ActivateOrDeactivateUser(ctx context.Context, userID, loggedUser string, active bool) error {
	// 1. Verificar si el usuario existe
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		logs.Error("Failed to get user by ID", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return error2.NewDomainErrorWithCause("UserService", "ActivateOrDeactivateUser", "failed to get user by ID", err)
	}

	// 2. Verificar que el usuario no sea el mismo que el usuario que realiza la acción
	if user.ID == loggedUser {
		logs.Error("User cannot activate or deactivate itself", map[string]interface{}{
			"user_id": userID,
		})
		return error2.NewDomainError("UserService", "ActivateOrDeactivateUser", error2.ErrUserCannotActivateOrDeactivateItself.Error())
	}

	// 3. Verificar que el estado a cambiar sea diferente al actual
	if user.IsActive == active {
		logs.Error("User is already active or inactive", map[string]interface{}{
			"user_id": userID,
		})
		return error2.NewDomainError("UserService", "ActivateOrDeactivateUser", error2.ErrUserAlreadyActiveOrInactive.Error())
	}

	// 4. Activar o desactivar el usuario
	err = s.userRepo.ActivateOrDeactivate(ctx, userID, active)
	if err != nil {
		logs.Error("Failed to activate or deactivate user", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return error2.NewDomainErrorWithCause("UserService", "ActivateOrDeactivateUser", "failed to activate or deactivate user", err)
	}

	return nil
}

func (s *userService) AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy string) error {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	_, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return err
	}

	// 2. Verificar que el usuarios no tenga el rol asignado
	hasRole, err := s.hasRole(ctx, userID, roleID)
	if err != nil {
		return err
	}
	if hasRole {
		logs.Error("User already has role", map[string]interface{}{
			"user_id": userID,
			"role_id": roleID,
		})
		return error2.NewDomainError("UserService", "AssignRoleToUser", error2.ErrUserAlreadyHasRole.Error())
	}

	// 3. Asignar el rol al usuario
	err = s.userRepo.AssignRoleToUser(ctx, userID, roleID, assignedBy)
	if err != nil {
		logs.Error("Failed to assign role to user", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"role_id": roleID,
		})
		return error2.NewDomainErrorWithCause("UserService", "AssignRoleToUser", "failed to assign role to user", err)
	}

	return nil
}

func (s *userService) UpdateRolesToUser(ctx context.Context, userID, loggedUserID string, roles []entities.Role) error {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	_, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return err
	}

	// 2. Actualizar los roles del usuario
	err = s.userRepo.UpdateRolesToUser(ctx, userID, loggedUserID, roles)
	if err != nil {
		logs.Error("Failed to update user roles", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return error2.NewDomainErrorWithCause("UserService", "UpdateRolesToUser", "failed to update user roles", err)
	}

	return nil
}

func (s *userService) UnassignRole(ctx context.Context, userID, roleID string) error {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	_, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return err
	}

	// 2. Verificar que el usuario tenga el rol asignado
	hasRole, err := s.hasRole(ctx, userID, roleID)
	if err != nil {
		return error2.NewDomainErrorWithCause("UserService", "UnassignRole", "failed to check if user has role", err)
	}
	if !hasRole {
		logs.Error("User does not have role", map[string]interface{}{
			"user_id": userID,
			"role_id": roleID,
		})
		return error2.NewDomainError("UserService", "UnassignRole", error2.ErrUserDoesNotHaveRole.Error())
	}

	// 3. Desasignar el rol al usuario
	err = s.userRepo.RemoveRoleFromUser(ctx, userID, roleID)
	if err != nil {
		logs.Error("Failed to remove role from user", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
			"role_id": roleID,
		})
		return error2.NewDomainErrorWithCause("UserService", "UnassignRole", "failed to remove role from user", err)
	}

	return nil
}

func (s *userService) CleanAllSessions(ctx context.Context, userID string) error {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	_, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return err
	}

	// 2. Limpiar las sesiones activas del usuario
	err = s.userRepo.CleanExpiredSessions(ctx, userID)
	if err != nil {
		logs.Error("Failed to clean user sessions", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return error2.NewDomainErrorWithCause("UserService", "CleanAllSessions", "failed to clean user sessions", err)
	}

	return nil
}

func (s *userService) GetUserRoles(ctx context.Context, userID string) ([]entities.Role, error) {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	_, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 2. Obtener los roles del usuario
	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		logs.Error("Failed to get user roles", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return nil, error2.NewDomainErrorWithCause("UserService", "GetUserRoles", "failed to get user roles", err)
	}

	return roles, nil
}

func (s *userService) hasRole(ctx context.Context, userID, roleID string) (bool, error) {
	// 1. Obtener los roles del usuario
	roles, err := s.userRepo.GetUserRoles(ctx, userID)
	if err != nil {
		logs.Error("Failed to get user roles", map[string]interface{}{
			"error":   err.Error(),
			"user_id": userID,
		})
		return false, error2.NewDomainErrorWithCause("UserService", "HasRole", "failed to get user roles", err)
	}

	// 2. Verificar si el usuario tiene el rol
	for _, role := range roles {
		if role.ID == roleID {
			return true, nil
		}
	}

	return false, nil
}
