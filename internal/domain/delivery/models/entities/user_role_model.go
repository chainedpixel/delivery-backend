package entities

import (
	"time"
)

type UserRole struct {
	UserID     string    `gorm:"column:user_id;type:char(36);primaryKey;not null" json:"user_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	RoleID     string    `gorm:"column:role_id;type:char(36);primaryKey;not null" json:"role_id" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	AssignedAt time.Time `gorm:"column:assigned_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"assigned_at" example:"2021-01-01T00:00:00Z"`
	AssignedBy string    `gorm:"column:assigned_by;type:char(36);not null" json:"assigned_by" example:"3fa85f64-5717-4562-b3fc-2c963f66afa6"`
	IsActive   bool      `gorm:"column:is_active;type:boolean;default:true" json:"-" example:"true"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"-" example:"2021-01-01T00:00:00Z"`

	// Inverse Relationships
	User *User `gorm:"foreignKey:UserID;references:ID" json:"-"`

	// Relationships
	Role *Role `gorm:"foreignKey:RoleID;references:ID" json:"auth,omitempty"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
