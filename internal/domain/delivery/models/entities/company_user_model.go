package entities

import (
	"time"
)

type CompanyUser struct {
	UserID          string    `gorm:"column:user_id;type:char(36);primaryKey"`
	CompanyID       string    `gorm:"column:company_id;type:char(36);primaryKey"`
	BranchID        string    `gorm:"column:branch_id;type:char(36);primaryKey"`
	Position        string    `gorm:"column:position;type:varchar(100);not null"`
	Department      string    `gorm:"column:department;type:varchar(100)"`
	Permissions     string    `gorm:"column:permissions;type:json"`
	CanCreateOrders bool      `gorm:"column:can_create_orders;type:boolean;default:false"`
	CreatedAt       time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Relationships
	Branch  *Branch  `gorm:"foreignKey:BranchID;references:ID"`
	User    *User    `gorm:"foreignKey:UserID;references:ID"`
	Company *Company `gorm:"foreignKey:CompanyID;references:ID"`
}

func (CompanyUser) TableName() string {
	return "company_users"
}
