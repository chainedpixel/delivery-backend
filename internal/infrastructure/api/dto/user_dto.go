package dto

import error2 "infrastructure/error"

type UserDTO struct {
	// Email of the user to be created
	// @example example@example.com
	// @required
	Email string `json:"email"`

	// FirstName of the user to be created
	// @example John
	// @required
	FirstName string `json:"first_name"`

	// LastName of the user to be created
	// @example Doe
	// @required
	LastName string `json:"last_name"`

	// Phone number of the user to be created
	// @example 21212828
	// @required
	Phone string `json:"phone"`

	// Password of the user to be created
	// @example password
	// @required
	Password string `json:"password"`

	// Profile of the user to be created
	// @required
	Profile *UserProfileDTO `json:"profile"`
}

type UserProfileDTO struct {
	// DocumentType of the user to be created
	// @example DUI
	// @required
	DocumentType string `json:"document_type"`

	// DocumentNumber of the user to be created
	// @example 123456789
	// @required
	DocumentNumber string `json:"document_number"`

	// BirthDate of the user to be created
	// @example 1990-01-01
	// @required
	BirthDate string `json:"birth_date"`

	// EmergencyContactName of the user to be created
	// @example Jane Doe
	// @required
	EmergencyContactName string `json:"emergency_contact_name"`

	// EmergencyContactPhone of the user to be created
	// @example 21212828
	// @required
	EmergencyContactPhone string `json:"emergency_contact_phone"`

	// AdditionalInfo of the user to be created
	AdditionalInfo string `json:"additional_info"`
}

func (u *UserDTO) Validate() error {
	if u.Email == "" || u.FirstName == "" || u.LastName == "" || u.Phone == "" || u.Password == "" {
		return error2.NewGeneralServiceError("UserDTO", "Validate", error2.ErrInvalidUser)
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
	// RoleID UUID of the user to be assigned
	// @example admin
	// @required
	Role string `json:"role_id"`
}

func (r *AssignRoleDTO) Validate() error {
	if r.Role == "" {
		return error2.NewGeneralServiceError("AssignRoleDTO", "Validate", error2.ErrRoleMissing)
	}

	return nil
}

type ActivateUserDTO struct {
	// Active status of the user to be updated
	// @example true
	// @required
	Active bool `json:"active"`
}

type UpdateUserDTO struct {
	// Email of the user to be updated
	// @example example@example.com
	Email string `json:"email"`

	// FirstName of the user to be updated
	// @example John
	FirstName string `json:"first_name"`

	// LastName of the user to be updated
	// @example Doe
	LastName string `json:"last_name"`

	// Phone of the user to be updated
	// @example 21212828
	Phone string `json:"phone"`

	// Active status of the user to be updated
	// @example true
	Active bool `json:"active"`

	// Password of the user to be updated
	// @example password
	Password string `json:"password"`

	// Profile of the user to be updated
	Profile *UpdateUserProfileDTO `json:"profile"`
}

type UpdateUserProfileDTO struct {
	// DocumentType of the user to be updated
	// @example DUI
	DocumentType string `json:"document_type"`

	// DocumentNumber of the user to be updated
	// @example 123456789
	DocumentNumber string `json:"document_number"`

	// BirthDate of the user to be updated
	// @example 1990-01-01
	BirthDate string `json:"birth_date"`

	// EmergencyContactName of the user to be updated
	// @example Jane Doe
	EmergencyContactName string `json:"emergency_contact_name"`

	// EmergencyContactPhone of the user to be updated
	// @example 21212828
	EmergencyContactPhone string `json:"emergency_contact_phone"`

	// AdditionalInfo of the user to be updated
	AdditionalInfo string `json:"additional_info"`
}
