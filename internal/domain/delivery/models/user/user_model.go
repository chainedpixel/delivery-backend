package user

import "time"

type User struct {
	ID              string        `gorm:"column:id;type:char(36);primary_key" json:"id"`
	Email           string        `gorm:"column:email;type:varchar(255);uniqueIndex;not null" json:"email"`
	PasswordHash    string        `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	FullName        string        `gorm:"column:full_name;type:varchar(255);not null" json:"full_name"`
	Phone           string        `gorm:"column:phone;type:varchar(20)" json:"phone"`
	IsActive        bool          `gorm:"column:is_active;type:boolean;default:true" json:"is_active"`
	EmailVerifiedAt *time.Time    `gorm:"column:email_verified_at;type:timestamp null" json:"email_verified_at,omitempty"`
	PhoneVerifiedAt *time.Time    `gorm:"column:phone_verified_at;type:timestamp null" json:"phone_verified_at,omitempty"`
	CreatedAt       time.Time     `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt       time.Time     `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt       *time.Time    `gorm:"column:deleted_at;type:timestamp null" json:"deleted_at,omitempty"`
	Profile         *UserProfile  `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	Sessions        []UserSession `gorm:"foreignKey:UserID" json:"sessions,omitempty"`
	Roles           []UserRole    `gorm:"foreignKey:UserID" json:"roles,omitempty"`
}

func (User) TableName() string {
	return "users"
}
