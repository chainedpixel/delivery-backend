package entities

import "time"

type RolePermission struct {
	RoleID       string    `gorm:"column:role_id;type:char(36);primaryKey;not null" json:"role_id"`
	PermissionID string    `gorm:"column:permission_id;type:char(36);primaryKey;not null" json:"permission_id"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"-"`

	// Inverse Relationships
	Role       *Role       `gorm:"foreignKey:RoleID;references:ID" json:"-"`
	Permission *Permission `gorm:"foreignKey:PermissionID;references:ID" json:"-"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}
