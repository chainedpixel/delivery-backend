package request_mapper

import (
	"errors"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/constants"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/models/entities"
	"github.com/MarlonG1/delivery-backend/internal/domain/delivery/value_objects"
	"github.com/MarlonG1/delivery-backend/internal/infrastructure/api/dto"
	error2 "github.com/MarlonG1/delivery-backend/internal/infrastructure/error"
	"github.com/MarlonG1/delivery-backend/pkg/shared/logs"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

var formats = []string{
	"02/01/2006",
	"02-01-2006",
}

func UserRequestToModel(req *dto.UserDTO) (*entities.User, error) {
	userID := uuid.NewString()

	if !value_objects.NewEmail(req.Email).IsValid() {
		return nil, error2.NewGeneralServiceError("UserMapper", "UserRequestToModel", errors.New("invalid email format example@example.com"))
	}

	if !value_objects.NewPhoneNumber(req.Phone).IsValid() {
		return nil, error2.NewGeneralServiceError("UserMapper", "UserRequestToModel", errors.New("invalid phone format, must be 00000000"))
	}

	profile, err := userProfileToModel(req.Profile, userID)
	if err != nil {
		return nil, error2.NewGeneralServiceError("UserMapper", "UserRequestToModel", err)
	}

	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, error2.NewGeneralServiceError("UserMapper", "UserRequestToModel", errors.New("error hashing password"))
	}

	user := &entities.User{
		ID:           userID,
		Email:        req.Email,
		FullName:     req.FullName,
		Phone:        req.Phone,
		IsActive:     true,
		PasswordHash: string(password),
		Profile:      profile,
	}

	var roles []entities.UserRole
	for _, role := range req.Roles {
		roles = append(roles, entities.UserRole{
			Role: &entities.Role{
				Name: strings.ToUpper(role),
			},
		})
	}
	user.Roles = roles

	return user, nil
}

func userProfileToModel(req *dto.UserProfileDTO, id string) (*entities.Profile, error) {
	var date time.Time
	var parseErr, err error
	for _, format := range formats {
		date, err = time.Parse(format, req.BirthDate)
		if err == nil {
			break
		}
		parseErr = err
	}

	if parseErr != nil {
		return nil, error2.NewGeneralServiceError("UserMapper", "UserProfileToModel", errors.New("invalid date time for birthdate format, must be 1990-01-01 or 01/01/1990"))
	}

	if !value_objects.NewPhoneNumber(req.EmergencyContactPhone).IsValid() {
		logs.Info(req.EmergencyContactPhone)
		return nil, error2.NewGeneralServiceError("UserMapper", "UserProfileToModel", errors.New("invalid emergency phone format, must be 00000000"))
	}

	return &entities.Profile{
		DocumentType:          req.DocumentType,
		DocumentNumber:        req.DocumentNumber,
		BirthDate:             &date,
		EmergencyContactName:  req.EmergencyContactName,
		EmergencyContactPhone: req.EmergencyContactPhone,
		AdditionalInfo:        req.AdditionalInfo,
	}, nil
}

func UpdateUserRequestToModel(req *dto.UpdateUserDTO) (*entities.User, error) {
	if req == nil {
		return nil, error2.NewGeneralServiceError("UserMapper", "UpdateUserRequestToModel", errors.New("request is nil"))
	}

	if req.Email != "" {
		if !value_objects.NewEmail(req.Email).IsValid() {
			return nil, error2.NewGeneralServiceError("UserMapper", "UpdateUserRequestToModel", errors.New("invalid email format"))
		}
	}

	if req.Phone != "" {
		if !value_objects.NewPhoneNumber(req.Phone).IsValid() {
			return nil, error2.NewGeneralServiceError("UserMapper", "UpdateUserRequestToModel", errors.New("invalid phone format"))
		}
	}

	var date time.Time
	var err error
	if req.Profile != nil {
		if req.Profile.DocumentNumber != "" {
			if !value_objects.NewPhoneNumber(req.Profile.EmergencyContactPhone).IsValid() {
				return nil, error2.NewGeneralServiceError("UserMapper", "UpdateUserRequestToModel", errors.New("invalid phone format"))
			}
		}

		var parseErr error
		for _, format := range formats {
			date, err = time.Parse(format, req.Profile.BirthDate)
			if err == nil {
				break
			}
			parseErr = err
		}

		if parseErr != nil {
			return nil, error2.NewGeneralServiceError("UserMapper", "UserProfileToModel", errors.New("invalid date time for birthdate format, must be 1990-01-01 or 01/01/1990"))
		}
	}

	var password []byte
	if req.Password != "" {
		password, err = bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, error2.NewGeneralServiceError("UserMapper", "UpdateUserRequestToModel", errors.New("error hashing password"))
		}
	}

	user := &entities.User{
		Email:        req.Email,
		Phone:        req.Phone,
		IsActive:     req.Active,
		PasswordHash: string(password),
	}

	if req.Roles != nil {
		if len(req.Roles) == 0 {
			return nil, error2.NewGeneralServiceError("UpdateUserDTO", "Validate", error2.ErrMissingRoles)
		}

		var roles []entities.UserRole
		for _, role := range req.Roles {
			if role == "" {
				return nil, error2.NewGeneralServiceError("UpdateUserDTO", "Validate", error2.ErrRoleMissing)
			}

			if !constants.ValidRoles[strings.ToUpper(role)] {
				return nil, error2.NewGeneralServiceError("UpdateUserDTO", "Validate", error2.ErrInvalidRole)
			}

			roles = append(roles, entities.UserRole{
				Role: &entities.Role{
					Name: strings.ToUpper(role),
				},
			})
		}
		user.Roles = roles
	}

	if req.FullName != "" {
		user.FullName = req.FullName
	}

	if req.Profile != nil {
		user.Profile = &entities.Profile{
			DocumentType:          req.Profile.DocumentType,
			DocumentNumber:        req.Profile.DocumentNumber,
			BirthDate:             &date,
			EmergencyContactName:  req.Profile.EmergencyContactName,
			EmergencyContactPhone: req.Profile.EmergencyContactPhone,
			AdditionalInfo:        req.Profile.AdditionalInfo,
		}
	}

	return user, nil
}
