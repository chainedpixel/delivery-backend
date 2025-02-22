package user

import (
	appPorts "application/ports"
	"context"
	"domain/delivery/models/user"
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

func (uc *UserUseCase) GetProfileInfo(ctx context.Context, userID string) (*user.User, error) {
	return uc.profileService.GetUserInfo(ctx, userID)
}
