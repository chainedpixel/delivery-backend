package warehouses

import (
	"domain/delivery/models/zones"
	"time"
)

type Warehouse struct {
	ID        string    `gorm:"column:id;type:char(36);primaryKey"`
	ZoneID    string    `gorm:"column:zone_id;type:char(36);not null"`
	Name      string    `gorm:"column:name;type:varchar(100);not null"`
	Address   string    `gorm:"column:address;type:varchar(255);not null"`
	Location  []byte    `gorm:"column:location;type:point;not null"`
	IsActive  bool      `gorm:"column:is_active;type:boolean;default:true"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP"`

	// Inverse Relationships
	Zone *zones.Zone `gorm:"foreignKey:ZoneID;references:ID"`

	// Relationships
	PackageTrackings []PackageTracking `gorm:"foreignKey:WarehouseID"`
	Inventory        []Inventory       `gorm:"foreignKey:WarehouseID"`
}

func (Warehouse) TableName() string {
	return "warehouse"
}
