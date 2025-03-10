package entities

import (
	"time"
)

type PackageTracking struct {
	ID          string     `gorm:"column:id;type:char(36);primaryKey"`
	OrderID     string     `gorm:"column:order_id;type:char(36);not null"`
	WarehouseID string     `gorm:"column:warehouse_id;type:char(36);not null"`
	Status      string     `gorm:"column:status;type:varchar(50);not null"`
	CollectorID string     `gorm:"column:collector_id;type:char(36);not null"`
	CollectedAt *time.Time `gorm:"column:collected_at;type:timestamp"`
	Notes       string     `gorm:"column:notes;type:text"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Order     *Order     `gorm:"foreignKey:OrderID;references:ID"`
	Warehouse *Warehouse `gorm:"foreignKey:WarehouseID;references:ID"`
	Collector *User      `gorm:"foreignKey:CollectorID;references:ID"`
}

func (PackageTracking) TableName() string {
	return "package_warehouse_tracking"
}
