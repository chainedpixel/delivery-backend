package companies

import (
	"domain/delivery/models/orders"
	"domain/delivery/models/zones"
	"time"
)

type Branch struct {
	ID             string    `gorm:"column:id;type:char(36);primaryKey"`
	CompanyID      string    `gorm:"column:company_id;type:char(36);not null"`
	Name           string    `gorm:"column:name;type:varchar(255);not null"`
	Code           string    `gorm:"column:code;type:varchar(50);not null"`
	ContactName    string    `gorm:"column:contact_name;type:varchar(255);not null"`
	ContactPhone   string    `gorm:"column:contact_phone;type:varchar(20);not null"`
	ContactEmail   string    `gorm:"column:contact_email;type:varchar(255);not null"`
	IsActive       bool      `gorm:"column:is_active;type:boolean;default:true"`
	ZoneID         string    `gorm:"column:zone_id;type:char(36);not null"`
	OperatingHours string    `gorm:"column:operating_hours;type:json"`
	CreatedAt      time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt      time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Company *Company    `gorm:"foreignKey:CompanyID;references:ID"`
	Zone    *zones.Zone `gorm:"foreignKey:ZoneID;references:ID"`

	// Relationships
	Orders []orders.Order `gorm:"foreignKey:BranchID"`
}

func (Branch) TableName() string {
	return "company_branches"
}
