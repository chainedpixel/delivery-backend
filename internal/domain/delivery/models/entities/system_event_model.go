package entities

import (
	"time"
)

type SystemEvent struct {
	ID         string    `gorm:"column:id;type:char(36);primaryKey"`
	EventType  string    `gorm:"column:event_type;type:varchar(50);not null"`
	Source     string    `gorm:"column:source;type:varchar(50);not null"`
	SourceID   string    `gorm:"column:source_id;type:char(36);not null"`
	EventData  string    `gorm:"column:event_data;type:json;not null"`
	Severity   string    `gorm:"column:severity;type:varchar(20);not null;default:INFO"`
	OccurredAt time.Time `gorm:"column:occurred_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	Logs []EventLog `gorm:"foreignKey:EventID"`
}

func (SystemEvent) TableName() string {
	return "system_events"
}
