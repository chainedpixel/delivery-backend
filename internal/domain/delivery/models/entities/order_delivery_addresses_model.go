package entities

import (
	"time"
)

type DeliveryAddress struct {
	OrderID        string    `gorm:"column:order_id;type:char(36);primaryKey"`
	RecipientName  string    `gorm:"column:recipient_name;type:varchar(255);not null"`
	RecipientPhone string    `gorm:"column:recipient_phone;type:varchar(20);not null"`
	AddressLine1   string    `gorm:"column:address_line1;type:varchar(255);not null"`
	AddressLine2   string    `gorm:"column:address_line2;type:varchar(255)"`
	City           string    `gorm:"column:city;type:varchar(100);not null"`
	State          string    `gorm:"column:state;type:varchar(100);not null"`
	PostalCode     string    `gorm:"column:postal_code;type:varchar(20)"`
	Location       []byte    `gorm:"column:location;type:point;null"`
	AddressNotes   string    `gorm:"column:address_notes;type:varchar(200)"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	Latitude  float64 `gorm:"-"`
	Longitude float64 `gorm:"-"`

	// Inverse Relationships
	Order *Order `gorm:"foreignKey:OrderID;references:ID"`
}

func (DeliveryAddress) TableName() string {
	return "delivery_addresses"
}
