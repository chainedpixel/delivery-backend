package user

import (
	appPorts "application/ports"
	"context"
	"domain/delivery/interfaces"
	"domain/delivery/models/users"
)

type UserUseCase struct {
	profileService interfaces.UserService
}

func NewUserProfileUseCase(profileService interfaces.UserService) appPorts.UserUseCase {
	return &UserUseCase{
		profileService: profileService,
	}
}

func (uc *UserUseCase) GetProfileInfo(ctx context.Context, userID string) (*users.User, error) {
	return uc.profileService.GetUserInfo(ctx, userID)
}
