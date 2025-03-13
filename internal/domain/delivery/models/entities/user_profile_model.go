package entities

import "time"

type Profile struct {
	UserID                string     `gorm:"column:user_id;type:char(36);primary_key" json:"user_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	DocumentType          string     `gorm:"column:document_type;type:varchar(20)" json:"document_type,omitempty" example:"DNI"`
	DocumentNumber        string     `gorm:"column:document_number;type:varchar(30)" json:"document_number,omitempty" example:"000000000"`
	BirthDate             *time.Time `gorm:"column:birth_date;type:date" json:"birth_date,omitempty" example:"2021-01-01"`
	ProfilePictureURL     string     `gorm:"column:profile_picture_url;type:varchar(255)" json:"profile_picture_url,omitempty"`
	EmergencyContactName  string     `gorm:"column:emergency_contact_name;type:varchar(255)" json:"emergency_contact_name,omitempty" example:"John Doe"`
	EmergencyContactPhone string     `gorm:"column:emergency_contact_phone;type:varchar(20)" json:"emergency_contact_phone,omitempty" example:"21212828"`
	AdditionalInfo        string     `gorm:"column:additional_info;type:text" json:"additional_info,omitempty" example:"Additional information"`
	CreatedAt             time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"-"`
	UpdatedAt             time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"-"`

	// Inverse Relationships
	User *User `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (Profile) TableName() string {
	return "user_profiles"
}
