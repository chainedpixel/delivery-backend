package user

import (
	appPorts "application/ports"
	"context"
	"domain/delivery/interfaces"
	"domain/delivery/models/entities"
)

type UserUseCase struct {
	profileService interfaces.Userer
}

func NewUserProfileUseCase(profileService interfaces.Userer) appPorts.UserUseCase {
	return &UserUseCase{
		profileService: profileService,
	}
}

func (uc *UserUseCase) GetProfileInfo(ctx context.Context, userID string) (*entities.User, error) {
	return uc.profileService.GetUserInfo(ctx, userID)
}
