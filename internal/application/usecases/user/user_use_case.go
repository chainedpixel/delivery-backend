package user

import (
	appPorts "application/ports"
	"context"
	"domain/delivery/interfaces"
	"domain/delivery/models/auth"
	"domain/delivery/models/entities"
)

type UsererUseCase struct {
	userService interfaces.Userer
}

func NewUserProfileUseCase(profileService interfaces.Userer) appPorts.UserUseCase {
	return &UsererUseCase{
		userService: profileService,
	}
}

func (uc *UsererUseCase) GetProfileInfo(ctx context.Context) (*entities.User, error) {
	// 1. Extraer el ID de los claims del contexto
	claims := ctx.Value("claims").(*auth.AuthClaims)

	// 1. Obtener la información del usuario
	user, err := uc.userService.GetUserInfo(ctx, claims.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UsererUseCase) CreateUser(ctx context.Context, user *entities.User) error {
	// 1. Crear el usuario
	err := uc.userService.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) UpdateUser(ctx context.Context, userID string, user *entities.User) error {
	// 1. Actualizar el usuario
	err := uc.userService.UpdateUser(ctx, userID, user)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) DeleteUser(ctx context.Context, userID string) error {
	// 1. Eliminar el usuario
	err := uc.userService.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) GetUserByID(ctx context.Context, userID string) (*entities.User, error) {
	// 1. Obtener la información del usuario
	user, err := uc.userService.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UsererUseCase) RecoverUser(ctx context.Context, id string) error {
	// 1. Recuperar el usuario
	err := uc.userService.RecoverUser(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) ActivateOrDeactivateUser(ctx context.Context, userID string, active bool) error {
	// 1. Extraer el ID de los claims del contexto
	claims := ctx.Value("claims").(*auth.AuthClaims)

	// 2. Activar o desactivar el usuario
	err := uc.userService.ActivateOrDeactivateUser(ctx, userID, claims.UserID, active)
	if err != nil {
		return err
	}

	return nil
}

func (uc *UsererUseCase) AssignRoleToUser(ctx context.Context, userID, roleID string) error {
	// 1. Extraer el ID de los claims del contexto
	claims := ctx.Value("claims").(*auth.AuthClaims)

	// 2. Asignar rol al usuario
	err := uc.userService.AssignRoleToUser(ctx, userID, roleID, claims.UserID)
	if err != nil {
		return err
	}

	return nil
}
