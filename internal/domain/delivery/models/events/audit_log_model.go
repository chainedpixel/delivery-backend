package events

import (
	user "domain/delivery/models/users"
	"time"
)

type AuditLog struct {
	ID         string    `gorm:"column:id;type:char(36);primaryKey"`
	UserID     string    `gorm:"column:user_id;type:char(36);not null"`
	Action     string    `gorm:"column:action;type:varchar(50);not null"`
	EntityType string    `gorm:"column:entity_type;type:varchar(50);not null"`
	EntityID   string    `gorm:"column:entity_id;type:char(36);not null"`
	OldValues  string    `gorm:"column:old_values;type:json"`
	NewValues  string    `gorm:"column:new_values;type:json"`
	IPAddress  string    `gorm:"column:ip_address;type:varchar(45)"`
	UserAgent  string    `gorm:"column:user_agent;type:text"`
	CreatedAt  time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relación
	User *user.User `gorm:"foreignKey:UserID;references:ID"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}
