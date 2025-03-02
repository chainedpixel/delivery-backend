package companies

import (
	"domain/delivery/models/orders"
	"time"
)

type Company struct {
	ID                string     `gorm:"column:id;type:char(36);primaryKey"`
	Name              string     `gorm:"column:name;type:varchar(255);not null"`
	LegalName         string     `gorm:"column:legal_name;type:varchar(255);not null"`
	TaxID             string     `gorm:"column:tax_id;type:varchar(50);not null"`
	ContactEmail      string     `gorm:"column:contact_email;type:varchar(255);not null"`
	ContactPhone      string     `gorm:"column:contact_phone;type:varchar(20);not null"`
	Website           string     `gorm:"column:website;type:varchar(255)"`
	IsActive          bool       `gorm:"column:is_active;type:boolean;default:true"`
	ContractDetails   string     `gorm:"column:contract_details;type:json"`
	DeliveryRate      float64    `gorm:"column:delivery_rate;type:decimal(10,2);not null"`
	LogoURL           string     `gorm:"column:logo_url;type:varchar(255)"`
	ContractStartDate time.Time  `gorm:"column:contract_start_date;type:timestamp;not null"`
	ContractEndDate   *time.Time `gorm:"column:contract_end_date;type:timestamp"`
	CreatedAt         time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt         time.Time  `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	Address  *Address       `gorm:"foreignKey:CompanyID"`
	Branches []Branch       `gorm:"foreignKey:CompanyID"`
	Orders   []orders.Order `gorm:"foreignKey:CompanyID"`
}

func (Company) TableName() string {
	return "companies"
}
