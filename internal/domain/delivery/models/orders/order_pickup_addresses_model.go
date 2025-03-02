package orders

import "time"

type PickupAddress struct {
	OrderID      string    `gorm:"column:order_id;type:char(36);primaryKey"`
	ContactName  string    `gorm:"column:contact_name;type:varchar(255);not null"`
	ContactPhone string    `gorm:"column:contact_phone;type:varchar(20);not null"`
	AddressLine1 string    `gorm:"column:address_line1;type:varchar(255);not null"`
	AddressLine2 string    `gorm:"column:address_line2;type:varchar(255)"`
	City         string    `gorm:"column:city;type:varchar(100);not null"`
	State        string    `gorm:"column:state;type:varchar(100);not null"`
	PostalCode   string    `gorm:"column:postal_code;type:varchar(20)"`
	Location     []byte    `gorm:"column:location;type:point;not null"`
	AddressNotes string    `gorm:"column:address_notes;type:json"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Order *Order `gorm:"foreignKey:OrderID;references:ID"`
}

func (PickupAddress) TableName() string {
	return "pickup_addresses"
}
