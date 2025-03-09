package entities

import (
	"time"
)

type NotificationDevice struct {
	ID          string     `gorm:"column:id;type:char(36);primaryKey"`
	UserID      string     `gorm:"column:user_id;type:char(36);not null"`
	DeviceToken string     `gorm:"column:device_token;type:text;not null"`
	DeviceType  string     `gorm:"column:device_type;type:varchar(50);not null"`
	IsActive    bool       `gorm:"column:is_active;type:boolean;default:true"`
	LastUsedAt  *time.Time `gorm:"column:last_used_at;type:timestamp"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;references:ID"`
}

func (NotificationDevice) TableName() string {
	return "notification_devices"
}
