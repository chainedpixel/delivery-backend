package user

import (
	appPorts "application/ports"
	"context"
	"domain/delivery/models/users"
	"domain/delivery/ports"
)

type UserUseCase struct {
	profileService ports.UserService
}

func NewUserProfileUseCase(profileService ports.UserService) appPorts.UserUseCase {
	return &UserUseCase{
		profileService: profileService,
	}
}

func (uc *UserUseCase) GetProfileInfo(ctx context.Context, userID string) (*users.User, error) {
	return uc.profileService.GetUserInfo(ctx, userID)
}
