package events

import "time"

type EventLog struct {
	ID          string    `gorm:"column:id;type:char(36);primaryKey"`
	EventID     string    `gorm:"column:event_id;type:char(36);not null"`
	LogLevel    string    `gorm:"column:log_level;type:varchar(20);not null"`
	Description string    `gorm:"column:description;type:text;not null"`
	Metadata    string    `gorm:"column:metadata;type:json"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Event *SystemEvent `gorm:"foreignKey:EventID;references:ID"`
}

func (EventLog) TableName() string {
	return "event_logs"
}
