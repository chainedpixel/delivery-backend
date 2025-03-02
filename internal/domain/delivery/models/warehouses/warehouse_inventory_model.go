package warehouses

import (
	"domain/delivery/models/orders"
	"time"
)

type Inventory struct {
	ID            string     `gorm:"column:id;type:char(36);primaryKey"`
	WarehouseID   string     `gorm:"column:warehouse_id;type:char(36);not null"`
	OrderID       string     `gorm:"column:order_id;type:char(36);not null"`
	Status        string     `gorm:"column:status;type:varchar(50);not null"`
	ShelfLocation string     `gorm:"column:shelf_location;type:varchar(50)"`
	ReceivedAt    time.Time  `gorm:"column:received_at;type:timestamp;not null"`
	DispatchedAt  *time.Time `gorm:"column:dispatched_at;type:timestamp"`
	CreatedAt     time.Time  `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Warehouse *Warehouse    `gorm:"foreignKey:WarehouseID;references:ID"`
	Order     *orders.Order `gorm:"foreignKey:OrderID;references:ID"`
}

func (Inventory) TableName() string {
	return "warehouse_inventory"
}
