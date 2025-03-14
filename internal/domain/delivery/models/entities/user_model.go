package entities

import "time"

type User struct {
	ID              string     `gorm:"column:id;type:char(36);primary_key" json:"id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	CompanyID       string     `gorm:"column:company_id;type:char(36);not null" json:"company_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	Email           string     `gorm:"column:email;type:varchar(255);unique;not null" json:"email" example:"example@example.com"`
	PasswordHash    string     `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	FullName        string     `gorm:"column:full_name;type:varchar(255);not null" json:"full_name" example:"John Doe"`
	Phone           string     `gorm:"column:phone;type:varchar(20)" json:"phone" example:"21212828"`
	IsActive        bool       `gorm:"column:is_active;type:boolean;default:true" json:"is_active" example:"true"`
	EmailVerifiedAt *time.Time `gorm:"column:email_verified_at;type:timestamp null" json:"email_verified_at,omitempty" example:"2021-01-01T00:00:00Z"`
	PhoneVerifiedAt *time.Time `gorm:"column:phone_verified_at;type:timestamp null" json:"phone_verified_at,omitempty" example:"2021-01-01T00:00:00Z"`
	CreatedAt       time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at" example:"2021-01-01T00:00:00Z"`
	UpdatedAt       time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at" example:"2021-01-01T00:00:00Z"`
	DeletedAt       *time.Time `gorm:"column:deleted_at;type:timestamp null" json:"deleted_at,omitempty" example:"2021-01-01T00:00:00Z"`

	// Relationships
	Sessions []UserSession `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
	Company  *Company      `gorm:"foreignKey:CompanyID" json:"-"`

	Roles   []UserRole `gorm:"foreignKey:UserID" json:"roles,omitempty"`
	Profile *Profile   `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

func (User) TableName() string {
	return "users"
}
