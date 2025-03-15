package dto

import (
	"domain/delivery/constants"
	error2 "infrastructure/error"
	"strings"
)

type UserDTO struct {
	// Email of the user to be created
	// @required
	Email string `json:"email" example:"example@example.com"`

	// Fullname of the user to be created
	// @required
	FullName string `json:"full_name" example:"John Doe"`

	// Phone number of the user to be created
	// @required
	Phone string `json:"phone" example:"21212828"`

	// Password of the user to be created
	// @required
	Password string `json:"password" example:"password"`

	// Roles of the user to be created
	// @required
	Roles []string `json:"roles" example:"[\"admin\"]"`

	// Profile of the user to be created
	// @required
	Profile *UserProfileDTO `json:"profile"`
}

type UserProfileDTO struct {
	// DocumentType of the user to be created
	// @required
	DocumentType string `json:"document_type" example:"DNI"`

	// DocumentNumber of the user to be created
	// @required
	DocumentNumber string `json:"document_number" example:"12345678"`

	// BirthDate of the user to be created
	// @required
	BirthDate string `json:"birth_date" example:"01-01-1990 or 01/01/1990"`

	// EmergencyContactName of the user to be created
	// @required
	EmergencyContactName string `json:"emergency_contact_name" example:"John Doe"`

	// EmergencyContactPhone of the user to be created
	// @required
	EmergencyContactPhone string `json:"emergency_contact_phone" example:"21212828"`

	// AdditionalInfo of the user to be created
	AdditionalInfo string `json:"additional_info" example:"Additional information"`
}

func (u *UserDTO) Validate() error {
	if u.Email == "" || u.FullName == "" || u.Phone == "" || u.Password == "" {
		return error2.NewGeneralServiceError("UserDTO", "Validate", error2.ErrInvalidUser)
	}

	if len(u.Roles) == 0 {
		return error2.NewGeneralServiceError("UserDTO", "Validate", error2.ErrMissingRoles)
	}

	for _, role := range u.Roles {
		if role == "" {
			return error2.NewGeneralServiceError("UserDTO", "Validate", error2.ErrRoleMissing)
		}

		if !constants.ValidRoles[strings.ToUpper(role)] {
			return error2.NewGeneralServiceError("UserDTO", "Validate", error2.ErrInvalidRole)
		}
	}

	if u.Profile == nil {
		return error2.NewGeneralServiceError("UserDTO", "Validate", error2.ErrMissingProfileSection)
	}

	if u.Profile.DocumentType == "" || u.Profile.DocumentNumber == "" || u.Profile.BirthDate == "" || u.Profile.EmergencyContactName == "" || u.Profile.EmergencyContactPhone == "" {
		return error2.NewGeneralServiceError("UserDTO", "Validate", error2.ErrInvalidProfileUser)
	}

	return nil
}

type AssignRoleDTO struct {
	// RoleID or RoleName of the role to be assigned to the user
	// @required
	Role string `json:"role" example:"admin or 3fa85f64-5717-4562-b3fc-2c963f66afa6"`
}

func (r *AssignRoleDTO) Validate() error {
	if r.Role == "" {
		return error2.NewGeneralServiceError("AssignRoleDTO", "Validate", error2.ErrRoleMissing)
	}

	return nil
}

type ActivateUserDTO struct {
	// Active status of the user to be updated
	// @required
	Active bool `json:"active" example:"true"`
}

type UpdateUserDTO struct {
	// Email of the user to be updated
	Email string `json:"email" example:"example@example.com"`

	// Fullname of the user to be updated
	FullName string `json:"full_name" example:"John Doe"`

	// Phone of the user to be updated
	Phone string `json:"phone" example:"21212828"`

	// Active status of the user to be updated
	Active bool `json:"active" example:"true"`

	// Password of the user to be updated
	Password string `json:"password" example:"password"`

	Roles []string `json:"roles" example:"[\"admin\"]"`

	// Profile of the user to be updated
	Profile *UpdateUserProfileDTO `json:"profile"`
}

type UpdateUserProfileDTO struct {
	// DocumentType of the user to be updated
	DocumentType string `json:"document_type" example:"DNI"`

	// DocumentNumber of the user to be updated
	DocumentNumber string `json:"document_number" example:"12345678"`

	// BirthDate of the user to be updated
	BirthDate string `json:"birth_date" example:"01-01-1990 or 01/01/1990"`

	// EmergencyContactName of the user to be updated
	EmergencyContactName string `json:"emergency_contact_name" example:"John Doe"`

	// EmergencyContactPhone of the user to be updated
	EmergencyContactPhone string `json:"emergency_contact_phone" example:"21212828"`

	// AdditionalInfo of the user to be updated
	AdditionalInfo string `json:"additional_info" example:"Additional information"`
}
