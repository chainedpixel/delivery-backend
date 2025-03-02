package orders

import "time"

type PackageDetail struct {
	OrderID             string    `gorm:"column:order_id;type:char(36);primaryKey"`
	IsFragile           bool      `gorm:"column:is_fragile;type:boolean;default:false"`
	IsUrgent            bool      `gorm:"column:is_urgent;type:boolean;default:false"`
	Weight              float64   `gorm:"column:weight;type:decimal(10,2)"`
	Dimensions          string    `gorm:"column:dimensions;type:json"`
	SpecialInstructions string    `gorm:"column:special_instructions;type:text"`
	CreatedAt           time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Order *Order `gorm:"foreignKey:OrderID;references:ID"`
}

func (PackageDetail) TableName() string {
	return "package_details"
}
