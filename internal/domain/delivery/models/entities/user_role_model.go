package entities

import (
	"time"
)

type UserRole struct {
	UserID     string    `gorm:"column:user_id;type:char(36);primaryKey;not null" json:"user_id"`
	RoleID     string    `gorm:"column:role_id;type:char(36);primaryKey;not null" json:"role_id"`
	AssignedAt time.Time `gorm:"column:assigned_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"assigned_at"`
	AssignedBy string    `gorm:"column:assigned_by;type:char(36);not null" json:"assigned_by"`
	IsActive   bool      `gorm:"column:is_active;type:boolean;default:true" json:"-"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"-"`

	// Inverse Relationships
	User *User `gorm:"foreignKey:UserID;references:ID" json:"-"`

	// Relationships
	Role *Role `gorm:"foreignKey:RoleID;references:ID" json:"auth,omitempty"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
