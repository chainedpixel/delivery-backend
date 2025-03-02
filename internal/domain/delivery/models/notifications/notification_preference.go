package notifications

import (
	user "domain/delivery/models/users"
	"time"
)

type NotificationPreference struct {
	UserID           string    `gorm:"column:user_id;type:char(36);primaryKey"`
	NotificationType string    `gorm:"column:notification_type;type:varchar(50);primaryKey"`
	EmailEnabled     bool      `gorm:"column:email_enabled;type:boolean;default:true"`
	PushEnabled      bool      `gorm:"column:push_enabled;type:boolean;default:true"`
	SMSEnabled       bool      `gorm:"column:sms_enabled;type:boolean;default:false"`
	CreatedAt        time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt        time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	User *user.User `gorm:"foreignKey:UserID;references:ID"`
}

func (NotificationPreference) TableName() string {
	return "notification_preferences"
}
