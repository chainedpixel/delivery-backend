package companies

import "time"

type Address struct {
	CompanyID    string    `gorm:"column:company_id;type:char(36);primaryKey"`
	AddressLine1 string    `gorm:"column:address_line1;type:varchar(255);not null"`
	AddressLine2 string    `gorm:"column:address_line2;type:varchar(255)"`
	City         string    `gorm:"column:city;type:varchar(100);not null"`
	State        string    `gorm:"column:state;type:varchar(100);not null"`
	PostalCode   string    `gorm:"column:postal_code;type:varchar(20)"`
	Location     []byte    `gorm:"column:location;type:point;not null"`
	IsMain       bool      `gorm:"column:is_main;type:boolean;default:false"`
	CreatedAt    time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse relationships
	Company *Company `gorm:"foreignKey:CompanyID;references:ID"`
}

func (Address) TableName() string {
	return "company_addresses"
}
