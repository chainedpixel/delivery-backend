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

type userProfileService struct {
	userRepo ports.UserRepository
}

func NewUserService(userRepo ports.UserRepository) interfaces.Userer {
	return &userProfileService{
		userRepo: userRepo,
	}
}

func (s *userProfileService) GetUserInfo(ctx context.Context, userID string) (*entities.User, error) {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	user, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userProfileService) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	user, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userProfileService) CreateUser(ctx context.Context, user *entities.User) error {
	// Crear el usuario
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		logs.Error("Failed to create user", map[string]interface{}{
			"error": err.Error(),
		})
		return error2.NewDomainErrorWithCause("UserService", "CreateUser", "failed to create user", err)
	}

	return nil
}

func (s *userProfileService) UpdateUser(ctx context.Context, userID string, user *entities.User) error {
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

func (s *userProfileService) RecoverUser(ctx context.Context, userID string) error {
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

func (s *userProfileService) DeleteUser(ctx context.Context, userID string) error {
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

func (s *userProfileService) validateUserFromRepository(ctx context.Context, userID string) (*entities.User, error) {
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

func (s *userProfileService) ActivateOrDeactivateUser(ctx context.Context, userID, loggedUser string, active bool) error {
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

func (s *userProfileService) AssignRoleToUser(ctx context.Context, userID, roleID, assignedBy string) error {
	// 1. Verificar si el usuario existe, está activo y no está eliminado
	_, err := s.validateUserFromRepository(ctx, userID)
	if err != nil {
		return err
	}

	// 2. Asignar el rol al usuario
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
