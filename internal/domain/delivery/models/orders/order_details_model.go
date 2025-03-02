package orders

import "time"

type Details struct {
	OrderID           string     `gorm:"column:order_id;type:char(36);primaryKey"`
	Price             float64    `gorm:"column:price;type:decimal(10,2);not null"`
	Distance          float64    `gorm:"column:distance;type:decimal(10,2);not null"`
	PickupTime        time.Time  `gorm:"column:pickup_time;type:timestamp;not null"`
	DeliveryDeadline  time.Time  `gorm:"column:delivery_deadline;type:timestamp;not null"`
	DeliveredAt       *time.Time `gorm:"column:delivered_at;type:timestamp"`
	RequiresSignature bool       `gorm:"column:requires_signature;type:boolean;default:false"`
	DeliveryNotes     string     `gorm:"column:delivery_notes;type:json"`
	CreatedAt         time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Order *Order `gorm:"foreignKey:OrderID;references:ID"`
}

func (Details) TableName() string {
	return "order_details"
}
