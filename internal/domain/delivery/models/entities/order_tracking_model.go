package entities

import (
	"time"
)

type Tracking struct {
	OrderID         string    `gorm:"column:order_id;type:char(36);primaryKey"`
	CurrentLocation []string  `gorm:"column:current_location;type:point"`
	CurrentStatus   string    `gorm:"column:current_status;type:varchar(20);not null"`
	LastUpdated     time.Time `gorm:"column:last_updated;type:timestamp;default:CURRENT_TIMESTAMP"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Order *Order `gorm:"foreignKey:OrderID;references:ID"`
}

func (Tracking) TableName() string {
	return "order_tracking"
}
