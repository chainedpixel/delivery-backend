package notifications

import "time"

type NotificationTemplate struct {
	ID              string    `gorm:"column:id;type:char(36);primaryKey"`
	Name            string    `gorm:"column:name;type:varchar(100);not null"`
	Type            string    `gorm:"column:type;type:varchar(50);not null"`
	TitleTemplate   string    `gorm:"column:title_template;type:text;not null"`
	ContentTemplate string    `gorm:"column:content_template;type:text;not null"`
	Variables       string    `gorm:"column:variables;type:json"`
	IsActive        bool      `gorm:"column:is_active;type:boolean;default:true"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	Notifications []Notification `gorm:"foreignKey:TemplateID"`
}

func (NotificationTemplate) TableName() string {
	return "notification_templates"
}
