package drivers

import (
	"domain/delivery/models/zones"
	"time"
)

type Availability struct {
	DriverID        string    `gorm:"column:driver_id;type:char(36);primaryKey"`
	CurrentZoneID   string    `gorm:"column:current_zone_id;type:char(36);not null"`
	CurrentLocation []byte    `gorm:"column:current_location;type:point;not null"`
	Status          string    `gorm:"column:status;type:varchar(20);not null"`
	LastUpdate      time.Time `gorm:"column:last_update;type:timestamp;default:CURRENT_TIMESTAMP"`
	ActiveOrders    int       `gorm:"column:active_orders;type:int;default:0"`
	CanTakeOrders   bool      `gorm:"column:can_take_orders;type:boolean;default:true"`
	ShiftStart      time.Time `gorm:"column:shift_start;type:timestamp;not null"`
	ShiftEnd        time.Time `gorm:"column:shift_end;type:timestamp;not null"`

	// Inverse Relationships
	Driver *Driver     `gorm:"foreignKey:DriverID;references:UserID"`
	Zone   *zones.Zone `gorm:"foreignKey:CurrentZoneID;references:ID"`
}

func (Availability) TableName() string {
	return "driver_availability"
}
