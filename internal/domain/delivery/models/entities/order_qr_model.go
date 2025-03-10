package entities

import (
	"time"
)

type QRCode struct {
	OrderID   string    `gorm:"column:order_id;type:char(36);primaryKey"`
	QRData    string    `gorm:"column:qr_data;type:text;not null"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Order *Order `gorm:"foreignKey:OrderID;references:ID"`
}

func (QRCode) TableName() string {
	return "qr_codes"
}
