package entities

import (
	"time"
)

type Notification struct {
	ID         string     `gorm:"column:id;type:char(36);primaryKey"`
	UserID     string     `gorm:"column:user_id;type:char(36);not null"`
	TemplateID string     `gorm:"column:template_id;type:char(36);not null"`
	Title      string     `gorm:"column:title;type:varchar(255);not null"`
	Content    string     `gorm:"column:content;type:text;not null"`
	Type       string     `gorm:"column:type;type:varchar(50);not null"`
	Metadata   string     `gorm:"column:metadata;type:json"`
	IsRead     bool       `gorm:"column:is_read;type:boolean;default:false"`
	ReadAt     *time.Time `gorm:"column:read_at;type:timestamp"`
	SentAt     *time.Time `gorm:"column:sent_at;type:timestamp"`
	CreatedAt  time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	User     *User                 `gorm:"foreignKey:UserID;references:ID"`
	Template *NotificationTemplate `gorm:"foreignKey:TemplateID;references:ID"`
}

func (Notification) TableName() string {
	return "notifications"
}
