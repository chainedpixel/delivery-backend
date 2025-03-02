package orders

import "time"

type StatusHistory struct {
	ID          string    `gorm:"column:id;type:char(36);primaryKey"`
	OrderID     string    `gorm:"column:order_id;type:char(36);not null"`
	Status      string    `gorm:"column:status;type:varchar(20);not null"`
	Description string    `gorm:"column:description;type:text"`
	CreatedAt   time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Order *Order `gorm:"foreignKey:OrderID;references:ID"`
}

func (StatusHistory) TableName() string {
	return "order_status_history"
}
